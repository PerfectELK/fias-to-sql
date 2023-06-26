package fias

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/models"
	"fias_to_sql/internal/services/fias/types"
	"fias_to_sql/internal/services/logger"
	"fias_to_sql/internal/services/terminal"
	"fias_to_sql/pkg/slice"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

var FileNames []string
var FileNamesToDump []string

func getSortedXmlFiles(zf *zip.ReadCloser) []*zip.File {
	filesMap := make(map[string][]*zip.File)
	filesMap["object"] = make([]*zip.File, 0)
	filesMap["house"] = make([]*zip.File, 0)
	filesMap["hierarchy"] = make([]*zip.File, 0)

	for _, file := range zf.File {
		if len(FileNames) > 0 && !slice.Contains(FileNames, file.Name) {
			continue
		}

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
		if strings.Contains(file.Name, "_DIVISION_") {
			continue
		}
		if strings.Contains(file.Name, "_OBJ_TYPES_") {
			continue
		}
		if strings.Contains(file.Name, "_PARAMS_") {
			objectType = "param"
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

		if len(filesMap["param"]) > i {
			sortedFiles = append(sortedFiles, filesMap["param"][i])
		} else {
			endFilesCounter++
		}

		if endFilesCounter == 4 {
			break
		}
	}
	return sortedFiles
}

func GetImportDestination() (string, error) {
	importDestination := config.GetConfig("IMPORT_DESTINATION")
	if importDestination == "" {
		importDestination = strings.ToLower(terminal.InputPrompt("input import destination (json/db): "))
		config.SetConfig("IMPORT_DESTINATION", importDestination)
	}
	if importDestination != "json" &&
		importDestination != "db" {
		return "", errors.New("incorrect import destination choose")
	}
	return importDestination, nil
}

func ImportXml(
	ctx context.Context,
	archivePath string,
	importDestinationStr ...string,
) error {
	ArchivePathToDump = archivePath
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

	threadNumber := 3
	if tn := config.GetConfig("APP_THREAD_NUMBER"); tn != "" {
		threadNumber, _ = strconv.Atoi(tn)
	}
	mutexChan := make(chan struct{}, threadNumber)
	g, onErrCtx := errgroup.WithContext(context.Background())
	for _, file := range files {
		if ctx.Err() != nil {
			FileNamesToDump = append(FileNamesToDump, file.Name)
			continue
		}

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
		if strings.Contains(file.Name, "_PARAMS_") {
			objectType = "param"
		}

		_file := file
		mutexChan <- struct{}{}
		g.Go(func() error {
			select {
			case <-onErrCtx.Done():
				<-mutexChan
				return nil
			default:
				c, err := _file.Open()
				if err != nil {
					<-mutexChan
					return err
				}

				amount, err := ProcessingXml(
					c,
					objectType,
					func(ol *types.FiasObjectList) {
						switch importDestination {
						case "db":
							err = importToDb(ol)
						case "json":
							err = importToJson(ol)
						default:
							err = importToDb(ol)
						}
					},
				)
				if err != nil {
					<-mutexChan
					return err
				}
				<-mutexChan
				logger.Println(_file.Name, ": records amount (", amount, ") [OK]")
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
	if len(list.List) == 0 {
		return nil
	}
	var modelList models.ModelListStruct

	for _, item := range list.List {
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
		case *types.Param:
			model := models.NewParam()
			model.SetObject_id(fiasObj.ObjectId)
			model.SetKladr_id(fiasObj.Value)
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
	return nil
}

func importToJson(list *types.FiasObjectList) error {
	pwd, _ := os.Getwd()

	for _, item := range list.List {
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
