package fias

import (
	"fias_to_sql/internal/config/fias"
	"github.com/go-rod/rod"
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
