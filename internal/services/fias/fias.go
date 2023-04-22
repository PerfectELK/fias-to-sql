package fias

import (
	"archive/zip"
	"errors"
	"fias_to_sql/internal/config/fias"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/terminal"
	"fmt"
	"github.com/go-rod/rod"
	"os"
	"path"
	"regexp"
	"strings"
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

	for _, file := range zf.File {
		if !strings.Contains(file.Name, fias.ADDRRESS_FILE_PART) {
			continue
		}
		if strings.Contains(file.Name, "_PARAMS_") {
			continue
		}
		if strings.Contains(file.Name, "_DIVISION_") {
			continue
		}

		c, err := file.Open()
		if err != nil {
			return err
		}
		ProcessingXml(c)
		// Todo debug
		return nil
	}

	return nil
}
