package download

import (
	"github.com/PerfectELK/go-import-fias/internal/services/disk"
	"github.com/PerfectELK/go-import-fias/internal/services/logger"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func File(
	link string,
	saveTo string,
) error {
	req, _ := http.NewRequest("GET", link, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(saveTo, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	buf := make([]byte, disk.MB)
	var downloaded int64
	var downloadedForSpeed int64
	tBegin := time.Now()
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n > 0 {
			tNow := time.Now()
			if tNow.Unix()-tBegin.Unix() > 10 || downloaded == 0 {
				var speed int64
				if downloadedForSpeed != 0 {
					speed = (downloadedForSpeed / (tNow.Unix() - tBegin.Unix())) / 1024 / 1024
					downloadedForSpeed = 0
				}

				tBegin = tNow
				go func() {
					message := "Downloading... " + strconv.FormatFloat(float64(downloaded)/float64(resp.ContentLength)*100, 'f', 6, 64) + "%"
					logger.Println(message)
					logger.Println("Speed: " + strconv.FormatInt(speed, 10) + " mb/sec")
				}()
			}

			f.Write(buf[:n])
			downloaded += int64(n)
			downloadedForSpeed += int64(n)
		}
	}
	return nil
}
