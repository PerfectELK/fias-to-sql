package db

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db/abstract"
	"fias_to_sql/pkg/db/mysql"
)

var dbInstance abstract.DbProcessor

func GetDbInstance() (abstract.DbProcessor, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	dbDriver := config.GetConfig("DB_DRIVER")
	switch dbDriver {
	case "MYSQL":
		dbInstance = &mysql.Processor{}
		if !dbInstance.IsConnected() {
			err := dbInstance.Connect()
			if err != nil {
				return nil, err
			}
		}
		return dbInstance, nil
	case "PGSQL":
		return nil, errors.New("PGSQL is not exists")
	}

	return nil, errors.New("doesn't selected db driver (MYSQL or PGSQL)")
}
