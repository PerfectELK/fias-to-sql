package fias

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/logger"
	"fias_to_sql/internal/services/shutdown"
	"fmt"
	"github.com/go-rod/rod"
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
