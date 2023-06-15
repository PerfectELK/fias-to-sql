package fias

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/logger"
	"fmt"
	"github.com/go-rod/rod"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetLinkOnNewestArchive() (string, *time.Time) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage(config.GetConfig("ARCHIVE_PAGE_LINK"))
	page.MustWaitLoad()

	err := page.WaitElementsMoreThan(config.GetConfig("ARCHIVE_LINK_SELECTOR"), 1)
	links := page.MustElements(config.GetConfig("ARCHIVE_LINK_SELECTOR"))

	timeRegex, err := regexp.Compile(`(\d\d\d\d.\d\d.\d\d)`)
	weightRegex, err := regexp.Compile(`(\([\d]+ (б|мб)\))`)
	if err != nil {
		return "", nil
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
			return "", nil
		}
		if linkTime.Unix() > highestTime.Unix() {
			highestTime = linkTime
			lastLink = *href
		}

	}
	return lastLink, &highestTime
}

func GetLastLocalArchivePath(highestTime time.Time) (string, bool) {
	entries, err := os.ReadDir(filepath.Join(os.Getenv("APP_ROOT"), "storage"))
	if err != nil {
		return "", false
	}

	lastLocalArchivePath := ""
	isLastArchive := false
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
			lastLocalArchivePath = fmt.Sprintf("./storage/%s", entry.Name())
			isLastArchive = true
			continue
		}
	}

	return lastLocalArchivePath, isLastArchive
}

func GetArchivePath() (string, error) {
	link, highestTime := GetLinkOnNewestArchive()
	if link == "" {
		return "", errors.New("cannot get link")
	}

	if config.GetConfig("IS_NEED_DOWNLOAD_ARCHIVE") != "true" {
		if archivePath, isNewest := GetLastLocalArchivePath(*highestTime); isNewest {
			return archivePath, nil
		}
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
