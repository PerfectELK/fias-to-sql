package app

import (
	"errors"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/download"
	"fias_to_sql/internal/services/fias"
	"fmt"
	"os"
	"path"
)

func App() error {
	usageGB, err := disk.Usage()
	if err != nil {
		return err
	}
	if usageGB.FreeGB < 70 {
		return errors.New("no space left on device")
	}

	link, err := fias.GetLinkOnNewestArchive()
	if err != nil {
		return err
	}

	pwd, _ := os.Getwd()
	pwd = path.Join(pwd, "archive.zip")
	download.File(link, pwd)
	fmt.Println("download complete")

	return nil
}
