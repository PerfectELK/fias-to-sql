package fias

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/models"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/fias/types"
	"fias_to_sql/internal/services/terminal"
	"fmt"
	"github.com/go-rod/rod"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func GetLinkOnNewestArchive() (string, error) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(config.GetConfig("ARCHIVE_PAGE_LINK"))
	page.MustWaitLoad()

	links := page.MustElements(config.GetConfig("ARCHIVE_LINK_SELECTOR"))
	r, err := regexp.Compile(`(\d\d\d\d.\d\d.\d\d)`)
	if err != nil {
		return "", err
	}
	var highestTime int64
	var lastLink string
	for _, link := range links {
		href := link.MustAttribute("href")
		if strings.Contains(*href, "delta") {
			continue
		}
		timeStr := r.FindString(*href)
		linkTime, err := time.Parse("2006.01.02", timeStr)
		if err != nil {
			return "", err
		}
		if linkTime.Unix() > highestTime {
			highestTime = linkTime.Unix()
			lastLink = *href
		}

	}

	return lastLink, nil
}

func GetArchivePath() (string, error) {
	//Todo debug
	return "E:\\Sources\\my_project\\fias_to_sql\\archive.zip", nil
	isHaveArchive := terminal.YesNoPrompt("do you have fias archive?")
	if isHaveArchive {
		archivePath := terminal.InputPrompt("enter full path to archive file")
		if _, err := os.Stat(archivePath); err == nil {
			return archivePath, nil
		} else {
			return "", errors.New("fias archive does not exists")
		}
	}

	fmt.Println("start downloading fias archive")
	link, err := GetLinkOnNewestArchive()
	if err != nil {
		return "", err
	}

	pwd, _ := os.Getwd()
	pwd = path.Join(pwd, "archive.zip")
	err = download.File(link, pwd)
	if err != nil {
		return "", err
	}
	fmt.Println("download complete")

	return pwd, nil
}

func getSortedXmlFiles(zf *zip.ReadCloser) []*zip.File {
	filesMap := make(map[string][]*zip.File)
	filesMap["object"] = make([]*zip.File, 0)
	filesMap["house"] = make([]*zip.File, 0)
	filesMap["hierarchy"] = make([]*zip.File, 0)

	for _, file := range zf.File {
		var objectType string
		if strings.Contains(file.Name, config.GetConfig("HOUSES_FILE_PART")) {
			objectType = "house"
		}
		if strings.Contains(file.Name, config.GetConfig("OBJECT_FILE_PART")) {
			objectType = "object"
		}
		if strings.Contains(file.Name, config.GetConfig("HIERARCHY_FILE_PART")) {
			objectType = "hierarchy"
		}
		if objectType == "" {
			continue
		}
		if strings.Contains(file.Name, "_PARAMS_") {
			continue
		}
		if strings.Contains(file.Name, "_DIVISION_") {
			continue
		}
		if strings.Contains(file.Name, "_OBJ_TYPES_") {
			continue
		}
		filesMap[objectType] = append(filesMap[objectType], file)
	}

	var sortedFiles []*zip.File
	for i := 0; ; i++ {
		endFilesCounter := 0
		if len(filesMap["object"]) > i {
			sortedFiles = append(sortedFiles, filesMap["object"][i])
		} else {
			endFilesCounter++
		}

		if len(filesMap["house"]) > i {
			sortedFiles = append(sortedFiles, filesMap["house"][i])
		} else {
			endFilesCounter++
		}

		if len(filesMap["hierarchy"]) > i {
			sortedFiles = append(sortedFiles, filesMap["hierarchy"][i])
		} else {
			endFilesCounter++
		}

		if endFilesCounter == 3 {
			break
		}
	}
	return sortedFiles
}

func GetImportDestination() (string, error) {
	importDestination := config.GetConfig("IMPORT_DESTINATION")
	if importDestination == "" {
		importDestination = strings.ToLower(terminal.InputPrompt("input import destination (json/db): "))
	}
	if importDestination != "json" &&
		importDestination != "db" {
		return "", errors.New("incorrect import destination choose")
	}
	return importDestination, nil
}

func ImportXml(
	archivePath string,
	importDestinationStr ...string,
) error {
	importDestination := "db"
	if len(importDestinationStr) > 0 {
		importDestination = importDestinationStr[0]
	}

	zf, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zf.Close()

	files := getSortedXmlFiles(zf)
	mutexChan := make(chan struct{}, 6)
	g, ctx := errgroup.WithContext(context.Background())
	for _, file := range files {
		var objectType string
		if strings.Contains(file.Name, config.GetConfig("HOUSES_FILE_PART")) {
			objectType = "house"
		}
		if strings.Contains(file.Name, config.GetConfig("OBJECT_FILE_PART")) {
			objectType = "object"
		}
		if strings.Contains(file.Name, config.GetConfig("HIERARCHY_FILE_PART")) {
			objectType = "hierarchy"
		}

		_file := file
		mutexChan <- struct{}{}
		g.Go(func() error {
			select {
			case <-ctx.Done():
				<-mutexChan
				return nil
			default:
				c, err := _file.Open()
				if err != nil {
					<-mutexChan
					return err
				}

				list, err := ProcessingXml(c, objectType)
				if err != nil {
					<-mutexChan
					return err
				}
				listLen := len(list.Addresses)

				switch importDestination {
				case "db":
					err = importToDb(list)
				case "json":
					err = importToJson(list)
				default:
					err = importToDb(list)
				}

				if err != nil {
					<-mutexChan
					return err
				}
				<-mutexChan
				fmt.Println(_file.Name, ": records amount (", listLen, ") [OK]")
				return nil
			}
		})
	}

	err = g.Wait()
	if err != nil {
		return err
	}

	if importDestination == "json" {
		pwd, _ := os.Getwd()
		err = fixJson(pwd + "/storage/addresses.json")
		if err != nil {
			return err
		}
		err = fixJson(pwd + "/storage/houses.json")
		if err != nil {
			return err
		}
		err = fixJson(pwd + "/storage/hierarchy.json")
		if err != nil {
			return err
		}
	}

	return err
}

func importToDb(list *types.FiasObjectList) error {
	if len(list.Addresses) == 0 {
		return nil
	}
	var modelList models.ModelListStruct

	for _, item := range list.Addresses {
		var err error
		switch fiasObj := item.(type) {
		case *types.Address:
			model := models.NewObject()
			model.SetName(fiasObj.Name)
			model.SetObject_id(fiasObj.ObjectId)
			model.SetObject_guid(fiasObj.ObjectGuid)
			model.SetLevel(fiasObj.Level)
			model.SetType_name(fiasObj.TypeName)
			modelList.AppendModel(model)
		case *types.House:
			model := models.NewObject()
			model.SetName(fiasObj.HouseNum)
			model.SetObject_id(fiasObj.ObjectId)
			model.SetObject_guid(fiasObj.ObjectGuid)
			model.SetLevel(12)
			model.SetType_name("дом")
			modelList.AppendModel(model)
		case *types.Hierarchy:
			model := models.NewHierarchy()
			model.SetObject_id(fiasObj.ObjectId)
			model.SetParent_object_id(fiasObj.ParentObjId)
			modelList.AppendModel(model)
		}
		if err != nil {
			return err
		}
	}
	err := modelList.SaveModelList()
	if err != nil {
		return err
	}
	list.Clear()
	return nil
}

func importToJson(list *types.FiasObjectList) error {
	pwd, _ := os.Getwd()

	for _, item := range list.Addresses {
		switch fiasObj := item.(type) {
		case *types.Address:
			addressesFile, _ := os.OpenFile(path.Join(pwd, "/storage/addresses.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			j, err := json.Marshal(fiasObj)
			if err != nil {
				return err
			}
			addressesFile.Write(j)
			addressesFile.Close()
		case *types.House:
			housesFile, _ := os.OpenFile(path.Join(pwd, "/storage/houses.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			j, err := json.Marshal(fiasObj)
			if err != nil {
				return err
			}
			housesFile.Write(j)
			housesFile.Close()
		case *types.Hierarchy:
			hierarchyFile, _ := os.OpenFile(path.Join(pwd, "/storage/hierarchy.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			j, err := json.Marshal(fiasObj)
			if err != nil {
				return err
			}
			hierarchyFile.Write(j)
			hierarchyFile.Close()
		}
	}
	return nil
}

func fixJson(filePath string) error {
	addressesFile, _ := os.Open(filePath)
	addressesTmpFile, _ := os.OpenFile(filePath+".tmp", os.O_CREATE|os.O_WRONLY, 644)

	byteLength := 1024
	b := make([]byte, byteLength)
	br := bufio.NewReader(addressesFile)
	addressesTmpFile.WriteString("[")
	for {
		n, err := br.Read(b)
		if err != nil && err != io.EOF {
			return err
		}

		if err != nil {
			break
		}
		str := string(b[0:n])

		begin := 0
		for p, ch := range str {
			if ch == '}' {
				addressesTmpFile.WriteString(str[begin:p+1] + ",")
				begin = p + 1
			}
		}

		addressesTmpFile.WriteString(str[begin:n])
	}

	addressesTmpFile.WriteString("]")
	stat, _ := addressesTmpFile.Stat()
	size := stat.Size()
	addressesTmpFile.WriteAt([]byte(" "), size-2)

	addressesFile.Close()
	os.Remove(filePath)
	addressesTmpFile.Close()
	os.Rename(filePath+".tmp", filePath)

	return nil
}
