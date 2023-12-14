package mysql

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
)

type MysqlMigrator struct{}

func (m MysqlMigrator) ObjectsTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
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
		"CREATE INDEX " + tableName + "_name_index ON " + tableName + " (name);" +
			" CREATE INDEX " + tableName + "_object_guid_index ON " + tableName + " (object_guid);" +
			" CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_type_name_index ON " + tableName + " (type_name);",
	)
}

func (m MysqlMigrator) ObjectTypesTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"level INT NOT NULL DEFAULT 0," +
			"short_name VARCHAR(255) NOT NULL DEFAULT ''," +
			"name VARCHAR(255) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
}

func (m MysqlMigrator) HierarchyTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`parent_object_id` INT NOT NULL DEFAULT 0) ENGINE=InnoDB;",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_parent_object_id_index ON " + tableName + " (parent_object_id);",
	)
}

func (m MysqlMigrator) KladrTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`kladr_id` VARCHAR(50) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_kladr_id_index ON " + tableName + " (kladr_id);",
	)
}

func (m MysqlMigrator) MigrateFromTempTables() error {
	err := m.dropOldTables()
	if err != nil {
		return err
	}
	return m.renameTables()
}

func (m MysqlMigrator) dropOldTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	originalObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE")
	originalObjectTypesTableName := config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE")
	originalHierarchyObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE")
	originalFiasKladrTableName := config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE")

	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalObjectTypesTableName + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalHierarchyObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalFiasKladrTableName + ";")
	return err
}

func (m MysqlMigrator) renameTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	originalObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE")
	originalObjectTypesTableName := config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE")
	originalHierarchyObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE")
	originalFiasKladrTableName := config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE")

	tempObjectsTable := config.GetConfig("DB_OBJECTS_TABLE")
	tempObjectTypesTableName := config.GetConfig("DB_OBJECT_TYPES_TABLE")
	tempHierarchyObjectsTable := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	tempFiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	err = dbInstance.Exec("RENAME TABLE " + tempObjectsTable + " TO " + originalObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("RENAME TABLE " + tempObjectTypesTableName + " TO " + originalObjectTypesTableName + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("RENAME TABLE " + tempHierarchyObjectsTable + " TO " + originalHierarchyObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("RENAME TABLE " + tempFiasKladrTableName + " TO " + originalFiasKladrTableName + ";")
	return err
}
