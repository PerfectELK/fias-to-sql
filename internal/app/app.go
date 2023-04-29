package app

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/error/handler"
	"fias_to_sql/internal/services/fias"
	"fias_to_sql/internal/services/terminal"
	"fias_to_sql/migrations"
	"fias_to_sql/pkg/db"
	"log"
	"time"
)

func App() error {
	err := config.InitConfig()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	usageGB, err := disk.Usage()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	if usageGB.FreeGB < 70 {
		return errors.New("no space left on device")
	}

	path, err := fias.GetArchivePath()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	_, err = db.GetDbInstance()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	err = migrations.CreateDatabase()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	err = migrations.CreateTables()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	importDestination := terminal.InputPrompt("input import destination (json/db): ")
	if importDestination != "json" &&
		importDestination != "db" {
		return errors.New("incorrect import destination choose")
	}
	beginTime := time.Now()
	err = fias.ImportXml(path, importDestination)
	if err != nil {
		return handler.ErrorHandler(err)
	}
	endTime := time.Now()
	log.Println("import time ", float64((endTime.Unix()-beginTime.Unix())/60), " minutes")

	return nil
}
