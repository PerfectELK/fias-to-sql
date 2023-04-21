package app

import (
	"errors"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/fias"
	"fmt"
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

	fmt.Println(link)

	return nil
}
