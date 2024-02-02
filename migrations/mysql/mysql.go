package mysql

import (
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/pkg/db"
)

type Migrator struct {
	objectsTable     string
	objectTypesTable string
	hierarchyTable   string
	kladrTable       string
}

func (m *Migrator) SetObjectsTable(table string) {
	m.objectsTable = table
}

func (m *Migrator) SetObjectTypesTable(table string) {
	m.objectTypesTable = table
}

func (m *Migrator) SetHierarchyTable(table string) {
	m.hierarchyTable = table
}

func (m *Migrator) SetKladrTable(table string) {
	m.kladrTable = table
}

func (m *Migrator) ObjectsTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE TABLE " + m.objectsTable + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`object_guid` VARCHAR(100) NOT NULL DEFAULT ''," +
			"`type_name` VARCHAR(100) NOT NULL DEFAULT ''," +
			"`level` INT NOT NULL DEFAULT 0," +
			"`name` VARCHAR(255) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
}

func (m *Migrator) ObjectTypesTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE TABLE " + m.objectTypesTable + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"level INT NOT NULL DEFAULT 0," +
			"short_name VARCHAR(255) NOT NULL DEFAULT ''," +
			"name VARCHAR(255) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
}

func (m *Migrator) HierarchyTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE TABLE " + m.hierarchyTable + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`parent_object_id` INT NOT NULL DEFAULT 0) ENGINE=InnoDB;",
	)
}

func (m *Migrator) KladrTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	return dbInstance.Exec(
		"CREATE TABLE " + m.kladrTable + " (" +
			"`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
			"`object_id` INT NOT NULL DEFAULT 0," +
			"`kladr_id` VARCHAR(50) NOT NULL DEFAULT '') ENGINE=InnoDB;",
	)
}

func (p *Migrator) CreateIndexes() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	err = dbInstance.Exec(
		"CREATE INDEX " + p.kladrTable + "_object_id_index ON " + p.kladrTable + " (object_id);" +
			" CREATE INDEX " + p.kladrTable + "_kladr_id_index ON " + p.kladrTable + " (kladr_id);",
	)
	if err != nil {
		return err
	}

	err = dbInstance.Exec(
		"CREATE INDEX " + p.hierarchyTable + "_object_id_index ON " + p.hierarchyTable + " (object_id);" +
			" CREATE INDEX " + p.hierarchyTable + "_parent_object_id_index ON " + p.hierarchyTable + " (parent_object_id);",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + p.objectsTable + "_name_index ON " + p.objectsTable + " (name);" +
			" CREATE INDEX " + p.objectsTable + "_object_guid_index ON " + p.objectsTable + " (object_guid);" +
			" CREATE INDEX " + p.objectsTable + "_object_id_index ON " + p.objectsTable + " (object_id);" +
			" CREATE INDEX " + p.objectsTable + "_type_name_index ON " + p.objectsTable + " (type_name);",
	)
}

func (m *Migrator) MigrateFromTempTables() error {
	err := m.dropOldTables()
	if err != nil {
		return err
	}
	return m.renameTables()
}

func (m *Migrator) dropOldTables() error {
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

func (m *Migrator) renameTables() error {
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
