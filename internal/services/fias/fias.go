package fias

import (
	"archive/zip"
	"errors"
	"fias_to_sql/internal/config/fias"
	"fias_to_sql/internal/services/db"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/terminal"
	"fmt"
	"github.com/go-rod/rod"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

func GetLinkOnNewestArchive() (string, error) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(fias.ARCHIVE_PAGE_LINK)
	page.MustWaitLoad()

	links := page.MustElements(fias.ARCHIVE_LINK_SELECTOR)
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
	// Todo debug
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

func ParseArchive(archivePath string) error {
	zf, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zf.Close()

	var wg sync.WaitGroup
	gorutinesCount := 0
	for _, file := range zf.File {
		var objectType string

		if strings.Contains(file.Name, fias.OBJECT_FILE_PART) {
			objectType = "object"
		}
		if strings.Contains(file.Name, fias.HOUSES_FILE_PART) {
			objectType = "house"
		}
		if strings.Contains(file.Name, fias.HIERARCHY_FILE_PART) {
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

		for {
			if gorutinesCount > 5 {
				time.Sleep(time.Second * 2)
			} else {
				break
			}
		}

		wg.Add(1)
		gorutinesCount += 1
		go func(file *zip.File) (err error) {
			defer wg.Done()
			c, err := file.Open()
			if err != nil {
				fmt.Println(file.Name+" [FAIL]", err)
				return err
			}

			list, err := ProcessingXml(c, objectType)
			if err != nil {
				fmt.Println(file.Name+" [FAIL]", err)
				return err
			}
			listLen := len(list.Addresses)
			err = db.ImportToDb(list)
			if err != nil {
				fmt.Println(file.Name+" [FAIL]", err)
				return err
			}
			fmt.Println(file.Name, ": records amount (", listLen, ") [OK]")
			gorutinesCount -= 1
			return err
		}(file)
	}
	wg.Wait()

	return nil
}
