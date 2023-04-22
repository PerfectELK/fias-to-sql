package app

import (
	"errors"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/fias"
)

func App() error {
	usageGB, err := disk.Usage()
	if err != nil {
		return err
	}
	if usageGB.FreeGB < 70 {
		return errors.New("no space left on device")
	}

	path, err := fias.GetArchivePath()
	if err != nil {
		return err
	}

	fias.ParseArchive(path)

	return nil
}
