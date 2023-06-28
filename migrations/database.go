package migrations

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
	"fias_to_sql/pkg/db/helpers"
	"fmt"
)

func CreateDatabase() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbName := config.GetConfig("DB_NAME")
	switch dbInstance.GetDriverName() {
	case "PGSQL":
		dbSchema := config.GetConfig("DB_SCHEMA")
		rows, err := dbInstance.Query(fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s'", dbName))
		if err != nil {
			return err
		}
		rowsArr := helpers.Scan(rows)
		if len(rowsArr) == 0 {
			return dbInstance.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		}
		return dbInstance.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS \"%s\"", dbSchema))
	case "MYSQL":
		return dbInstance.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	default:
		return errors.New("doesn't selected db driver")
	}

}
