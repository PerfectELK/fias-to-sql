package migrations

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/migrations/mysql"
	"fias_to_sql/migrations/pgsql"
	"fias_to_sql/pkg/db"
)

func CreateTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbName := config.GetConfig("DB_NAME")
	dbDriver := config.GetConfig("DB_DRIVER")
	err = dbInstance.Use(dbName)
	if err != nil {
		return err
	}

	fiasTableName := config.GetConfig("DB_OBJECTS_TABLE")
	fiasHierarchyTableName := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	fiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	_, tableCheck := dbInstance.Query("select * from " + fiasTableName + " LIMIT 1;")
	if tableCheck == nil {
		config.SetConfig("DB_ORIGINAL_OBJECTS_TABLE", config.GetConfig("DB_OBJECTS_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE", config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE", config.GetConfig("DB_OBJECTS_KLADR_TABLE"))
		fiasTableName = config.GetConfig("DB_OBJECTS_TABLE") + "_temp"
		fiasHierarchyTableName = config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE") + "_temp"
		fiasKladrTableName = config.GetConfig("DB_OBJECTS_KLADR_TABLE") + "_temp"
		_, tempTableCheck := dbInstance.Query("select * from " + fiasTableName + " LIMIT 1;")
		if tempTableCheck == nil {
			return errors.New("fias tables and temp tables is exists")
		}
		config.SetConfig("DB_OBJECTS_TABLE", fiasTableName)
		config.SetConfig("DB_OBJECTS_HIERARCHY_TABLE", fiasHierarchyTableName)
		config.SetConfig("DB_OBJECTS_KLADR_TABLE", fiasKladrTableName)
		config.SetConfig("DB_USE_TEMP_TABLE", "true")
		return createFiasTables(dbDriver, fiasTableName, fiasHierarchyTableName, fiasKladrTableName)
	}
	return createFiasTables(dbDriver, fiasTableName, fiasHierarchyTableName, fiasKladrTableName)
}

func MigrateDataFromTempTables() error {
	if config.GetConfig("DB_USE_TEMP_TABLE") != "true" {
		return nil
	}

	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	originalObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE")
	originalHierarchyObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE")

	tempObjectsTable := config.GetConfig("DB_OBJECTS_TABLE")
	tempHierarchyObjectsTable := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")

	dbInstance.Exec("DROP TABLE IF EXISTS " + originalObjectsTable + ";")
	dbInstance.Exec("DROP TABLE IF EXISTS " + originalHierarchyObjectsTable + ";")

	err = dbInstance.Exec("RENAME TABLE " + tempObjectsTable + " TO " + originalObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("RENAME TABLE " + tempHierarchyObjectsTable + " TO " + originalHierarchyObjectsTable + ";")
	return err
}

func createFiasTables(
	dbDriver string,
	fiasTableName string,
	fiasHierarchyTableName string,
	fiasKladrTableName string,
) error {
	switch dbDriver {
	case "MYSQL":
		err := mysql.ObjectsTableCreate(fiasTableName)
		if err != nil {
			return err
		}
		err = mysql.HierarchyTableCreate(fiasHierarchyTableName)
		if err != nil {
			return err
		}
		return mysql.KladrTableCreate(fiasKladrTableName)
	case "PGSQL":
		err := pgsql.ObjectsTableCreate(fiasTableName)
		if err != nil {
			return err
		}
		err = pgsql.HierarchyTableCreate(fiasHierarchyTableName)
		if err != nil {
			return err
		}
		return pgsql.KladrTableCreate(fiasKladrTableName)
	default:
		return nil
	}
}
