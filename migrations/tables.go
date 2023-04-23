package migrations

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
)

func CreateTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbName := config.GetConfig("DB_NAME")
	err = dbInstance.Use(dbName)
	if err != nil {
		return err
	}

	return nil
}
