package db

import (
	"errors"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/pkg/db/interfaces"
	"github.com/PerfectELK/go-import-fias/pkg/db/mysql"
	"github.com/PerfectELK/go-import-fias/pkg/db/pgsql"
)

var dbInstance interfaces.DbProcessor

func GetDbInstance() (interfaces.DbProcessor, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	dbDriver := config.GetConfig("DB_DRIVER")
	switch dbDriver {
	case "MYSQL":
		dbInstance = &mysql.Processor{}
	case "PGSQL":
		dbInstance = &pgsql.Processor{}
	}

	if dbInstance == nil {
		return nil, errors.New("doesn't selected db driver (MYSQL or PGSQL)")
	}

	if !dbInstance.IsConnected() {
		err := dbInstance.Connect()
		if err != nil {
			return nil, err
		}
	}

	return dbInstance, nil
}
