package migrations

import (
	"errors"
	"fias_to_sql/internal/config"
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

	_, tableCheck := dbInstance.Query("select * from " + fiasTableName + " LIMIT 1;")
	if tableCheck == nil {
		config.SetConfig("DB_ORIGINAL_OBJECTS_TABLE", config.GetConfig("DB_OBJECTS_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE", config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
		fiasTableName = config.GetConfig("DB_OBJECTS_TABLE") + "_temp"
		fiasHierarchyTableName = config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE") + "_temp"
		_, tempTableCheck := dbInstance.Query("select * from " + fiasTableName + " LIMIT 1;")
		if tempTableCheck == nil {
			return errors.New("fias tables and temp tables is exists")
		}
		config.SetConfig("DB_OBJECTS_TABLE", fiasTableName)
		config.SetConfig("DB_OBJECTS_HIERARCHY_TABLE", fiasHierarchyTableName)
		config.SetConfig("DB_USE_TEMP_TABLE", "true")
		return createFiasTables(dbDriver, fiasTableName, fiasHierarchyTableName)
	}
	return createFiasTables(dbDriver, fiasTableName, fiasHierarchyTableName)
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
) error {
	switch dbDriver {
	case "MYSQL":
		err := mysqlObjectsTableCreate(fiasTableName)
		if err != nil {
			return err
		}
		return mysqlHierarchyTableCreate(fiasHierarchyTableName)
	case "PGSQL":
		//Todo PGSQL db driver
		return nil
	default:
		return nil
	}
}

func mysqlObjectsTableCreate(fiasTableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + fiasTableName + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`object_guid` VARCHAR(100) NOT NULL DEFAULT ''," +
			"`type_name` VARCHAR(100) NOT NULL DEFAULT ''," +
			"`level` INT NOT NULL DEFAULT 0," +
			"`name` VARCHAR(255) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE INDEX " + fiasTableName + "_name_index ON " + fiasTableName + " (name);" +
			" CREATE INDEX " + fiasTableName + "_object_guid_index ON " + fiasTableName + " (object_guid);" +
			" CREATE INDEX " + fiasTableName + "_object_id_index ON " + fiasTableName + " (object_id);" +
			" CREATE INDEX " + fiasTableName + "_type_name_index ON " + fiasTableName + " (type_name);",
	)
}

func mysqlHierarchyTableCreate(fiasHierarchyTableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + fiasHierarchyTableName + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`parent_object_id` INT NOT NULL DEFAULT 0) ENGINE=InnoDB;",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + fiasHierarchyTableName + "_object_id_index ON " + fiasHierarchyTableName + " (object_id);" +
			" CREATE INDEX " + fiasHierarchyTableName + "_parent_object_id_index ON " + fiasHierarchyTableName + " (parent_object_id);",
	)
}
