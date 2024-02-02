package fias

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/internal/services/download"
	"github.com/PerfectELK/go-import-fias/internal/services/logger"
	"github.com/PerfectELK/go-import-fias/internal/services/shutdown"
	"github.com/go-rod/rod"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetLinkOnNewestArchive() string {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage(config.GetConfig("ARCHIVE_PAGE_LINK"))
	page.MustWaitLoad()

	err := page.WaitElementsMoreThan(config.GetConfig("ARCHIVE_LINK_SELECTOR"), 1)
	links := page.MustElements(config.GetConfig("ARCHIVE_LINK_SELECTOR"))

	timeRegex, err := regexp.Compile(`(\d\d\d\d.\d\d.\d\d)`)
	weightRegex, err := regexp.Compile(`(\([\d]+ (б|мб)\))`)
	if err != nil {
		return ""
	}
	var highestTime time.Time
	var lastLink string
	for _, link := range links {
		href := link.MustAttribute("href")
		if strings.Contains(*href, "delta") {
			continue
		}
		lowWeight := weightRegex.FindString(link.MustHTML())
		if lowWeight != "" {
			continue
		}
		timeStr := timeRegex.FindString(*href)
		linkTime, err := time.Parse("2006.01.02", timeStr)
		if err != nil {
			return ""
		}
		if linkTime.Unix() > highestTime.Unix() {
			highestTime = linkTime
			lastLink = *href
		}

	}
	return lastLink
}

func GetLastLocalArchivePath() string {
	entries, err := os.ReadDir(filepath.Join(os.Getenv("APP_ROOT"), "storage"))
	if err != nil {
		return ""
	}

	lastLocalArchivePath := ""
	var highestTime time.Time
	for _, entry := range entries {
		if !strings.Contains(entry.Name(), ".zip") {
			continue
		}
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		timeMod := info.ModTime()
		if timeMod.Unix() >= highestTime.Unix() {
			highestTime = timeMod
			lastLocalArchivePath = fmt.Sprintf("./storage/%s", entry.Name())
		}
	}

	return lastLocalArchivePath
}

func GetArchivePath() (string, error) {
	if shutdown.IsReboot {
		return shutdown.GetArchivePath(), nil
	}

	if localPath := config.GetConfig("ARCHIVE_LOCAL_PATH"); localPath != "" {
		return localPath, nil
	}

	if config.GetConfig("ARCHIVE_SOURCE") == "local" {
		if archivePath := GetLastLocalArchivePath(); archivePath != "" {
			return archivePath, nil
		}
	}

	var link string
	if config.GetConfig("ARCHIVE_SOURCE") == "link" {
		if link = config.GetConfig("ARCHIVE_LINK"); link == "" {
			return "", errors.New("archive-link var is empty")
		}
	}

	if link == "" {
		link = GetLinkOnNewestArchive()
	}

	if link == "" {
		return "", errors.New("cannot get link")
	}

	logger.Println("start downloading fias archive")
	now := time.Now()
	pwd, _ := os.Getwd()
	pwd = path.Join(pwd, fmt.Sprintf("storage/archive-%s.zip", now.Format("2006-01-02")))
	err := download.File(link, pwd)
	if err != nil {
		return "", err
	}
	logger.Println("download complete")

	return pwd, nil
}

func ExtractZipFiles(zf []*zip.File, destPath string) ([]string, error) {
	retArr := make([]string, 0, 0)

	_, err := os.ReadDir(destPath)
	if err != nil {
		return nil, err
	}

	// for tests only
	//counter := 0
	// for tests only
	buff := make([]byte, 1024*1024)
	for _, file := range zf {
		// for tests only
		//if counter > 50 {
		//	break
		//}
		// for tests only
		name := strings.Replace(file.Name, "\\", "_", -1)
		name = strings.Replace(name, "/", "_", -1)
		newFile, err := os.OpenFile(filepath.Join(destPath, name), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		for {
			n, err := reader.Read(buff)
			if n > 0 {
				_, err = newFile.Write(buff[:n])
				if err != nil {
					return nil, err
				}
			}
			if err != nil {
				if err != io.EOF {
					return nil, err
				}
				break
			}
		}
		retArr = append(retArr, newFile.Name())
		err = newFile.Close()
		if err != nil {
			return nil, err
		}
		err = reader.Close()
		if err != nil {
			return nil, err
		}
		// for tests only
		//counter++
		// for tests only
	}

	return retArr, nil
}
