package fias

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fias_to_sql/internal/config"
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

func ImportXml(
	archivePath string,
	importDestinationStr ...string,
) error {
	importDestination := "json"
	if len(importDestinationStr) > 0 {
		importDestination = importDestinationStr[0]
	}

	zf, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zf.Close()

	mutexChan := make(chan struct{}, 4)
	g, ctx := errgroup.WithContext(context.Background())
	for _, file := range zf.File {
		var objectType string

		if strings.Contains(file.Name, config.GetConfig("OBJECT_FILE_PART")) {
			objectType = "object"
		}
		if strings.Contains(file.Name, config.GetConfig("HOUSES_FILE_PART")) {
			objectType = "house"
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
		_file := file
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				mutexChan <- struct{}{}
				c, err := _file.Open()
				if err != nil {
					return err
				}

				list, err := ProcessingXml(c, objectType)
				if err != nil {
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
					return err
				}
				<-mutexChan
				fmt.Println(_file.Name, ": records amount (", listLen, ") [OK]")
				return nil
			}
		})
		//Todo debug
		break
	}

	err = g.Wait()
	if err != nil {
		return err
	}

	if importDestination == "json" {
		err = fixJsons()
	}

	return err
}

func importToDb(list *types.FiasObjectList) error {
	list.Clear()
	return nil
}

func importToJson(list *types.FiasObjectList) error {
	pwd, _ := os.Getwd()
	//housesFile, _ := os.OpenFile(path.Join(pwd, "/storage/houses.json"), os.O_CREATE|os.O_WRONLY, 0644)
	//hierarchyFile, _ := os.OpenFile(path.Join(pwd, "/storage/hierarchy.json"), os.O_CREATE|os.O_WRONLY, 0644)

	for _, item := range list.Addresses {
		switch fiasObj := item.(type) {
		case *types.Address:
			addressesFile, _ := os.OpenFile(path.Join(pwd, "/storage/addresses.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			defer addressesFile.Close()
			j, err := json.Marshal(fiasObj)
			if err != nil {
				return err
			}
			addressesFile.Write(j)
			//case *types.House
			//case *types.Hierarchy
		}
	}
	return nil
}

func fixJsons() error {
	pwd, _ := os.Getwd()
	addressesFile, _ := os.Open(path.Join(pwd, "/storage/addresses.json"))
	addressesTmpFile, _ := os.OpenFile(path.Join(pwd, "/storage/addresses.tmp.json"), os.O_CREATE|os.O_WRONLY, 644)

	b := make([]byte, 32*1024)
	addressesTmpFile.WriteString("[")
	for {
		_, err := addressesFile.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		str := string(b[:])
		startPos := 0
		for pos, ch := range str {
			if ch == '}' && str[pos+1] == '{' {
				addressesTmpFile.WriteString(str[startPos:pos+1] + ",")
				startPos = pos + 1
			}
			if ch == '}' && str[pos+1] != '{' {
				addressesTmpFile.WriteString(str[startPos : pos+1])
				startPos = 0
				break
			}
		}
		if startPos != 0 {
			addressesTmpFile.WriteString(str[startPos:] + ",")
		}
	}
	addressesTmpFile.WriteString("]")

	addressesFile.Close()
	os.Remove(path.Join(pwd, "/storage/addresses.json"))
	addressesTmpFile.Close()
	os.Rename(path.Join(pwd, "/storage/addresses.tmp.json"), path.Join(pwd, "/storage/addresses.json"))

	return nil
}
