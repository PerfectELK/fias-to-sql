package app

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/fias"
	"fias_to_sql/pkg/db"
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

	err = config.InitConfig()
	if err != nil {
		return err
	}

	_, err = db.GetDbInstance()
	if err != nil {
		return err
	}

	err = fias.ImportXmlToDb(path)
	if err != nil {
		return err
	}

	return nil
}
