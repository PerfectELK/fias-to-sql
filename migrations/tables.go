package migrations

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
)

func createTempTables() error {
	fiasTableName := config.GetConfig("DB_OBJECTS_TABLE") + "_temp"
	fiasHierarchyTableName := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE") + "_temp"
	dbDriver := config.GetConfig("DB_DRIVER")

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

	fiasTableName := config.GetConfig("DB_OBJECTS_TABLE")
	fiasHierarchyTableName := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")

	_, tableCheck := dbInstance.Query("select * from " + fiasTableName + " ;")
	if tableCheck == nil {
		return createTempTables()
	}

	dbDriver := config.GetConfig("DB_DRIVER")

	switch dbDriver {
	case "MYSQL":
		err = mysqlObjectsTableCreate(fiasTableName)
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
	return dbInstance.Exec("CREATE TABLE " + fiasTableName + " (" +
		"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"`object_id` INT DEFAULT NULL," +
		"`object_guid` VARCHAR(100) DEFAULT NULL," +
		"`type_name` VARCHAR(100) DEFAULT NULL," +
		"`level` INT DEFAULT NULL," +
		"`name` VARCHAR(255) DEFAULT NULL) ENGINE=InnoDB;" +
		"create index " + fiasTableName + "_name_index on " + fiasTableName + " (name);" +
		"create index " + fiasTableName + "_object_guid_index on " + fiasTableName + " (object_guid);" +
		"create index " + fiasTableName + "_object_id_index on " + fiasTableName + " (object_id);" +
		"create index " + fiasTableName + "_type_name_index on " + fiasTableName + " (type_name);",
	)
}

func mysqlHierarchyTableCreate(fiasHierarchyTableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec("CREATE TABLE " + fiasHierarchyTableName + " (" +
		"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"`object_id` INT DEFAULT NULL," +
		"`parent_object_id` INT DEFAULT NULL) ENGINE=InnoDB;" +
		"create index " + fiasHierarchyTableName + "_object_id_index on " + fiasHierarchyTableName + " (object_id);" +
		"create index " + fiasHierarchyTableName + "_parent_object_id_index on " + fiasHierarchyTableName + " (parent_object_id);",
	)
}
