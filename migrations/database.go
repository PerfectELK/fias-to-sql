package migrations

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
)

func CreateDatabase() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbName := config.GetConfig("DB_NAME")
	return dbInstance.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
}
