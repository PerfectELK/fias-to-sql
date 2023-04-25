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

	_, tableCheck := dbInstance.Query("select * from fias_objects;")
	if tableCheck == nil {
		dbInstance.Exec("drop table fias_objects")
	}

	err = dbInstance.Exec("CREATE TABLE `fias_objects` (" +
		"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"`object_id` INT DEFAULT NULL," +
		"`object_guid` VARCHAR(100) DEFAULT NULL," +
		"`type_name` VARCHAR(100) DEFAULT NULL," +
		"`level` INT DEFAULT NULL," +
		"`name` VARCHAR(255) DEFAULT NULL," +
		"`add_name` VARCHAR(255) DEFAULT NULL," +
		"`add_name2` VARCHAR(255) DEFAULT NULL) ENGINE=InnoDB;",
	)

	if err != nil {
		return err
	}

	_, tableCheck = dbInstance.Query("select * from fias_objects_hierarchy;")
	if tableCheck == nil {
		dbInstance.Exec("drop table fias_objects_hierarchy")
	}
	err = dbInstance.Exec("CREATE TABLE `fias_objects_hierarchy` (" +
		"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"`object_id` INT DEFAULT NULL," +
		"`parent_object_id` INT DEFAULT NULL) ENGINE=InnoDB;",
	)
	return err
}
