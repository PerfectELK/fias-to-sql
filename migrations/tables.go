package migrations

import (
	"errors"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/internal/services/shutdown"
	"github.com/PerfectELK/go-import-fias/migrations/mysql"
	"github.com/PerfectELK/go-import-fias/migrations/pgsql"
	"github.com/PerfectELK/go-import-fias/pkg/db"
)

type migrator interface {
	ObjectsTableCreate() error
	ObjectTypesTableCreate() error
	HierarchyTableCreate() error
	CreateIndexes() error
	KladrTableCreate() error
	MigrateFromTempTables() error

	SetObjectsTable(table string)
	SetObjectTypesTable(table string)
	SetHierarchyTable(table string)
	SetKladrTable(table string)
}

type viewCreator interface {
	CreateSettlementsView() error
	CreateSettlementsParentsView() error
}

var _m migrator

func getMigrator() migrator {
	if _m != nil {
		return _m
	}
	dbDriver := config.GetConfig("DB_DRIVER")
	switch dbDriver {
	case "MYSQL":
		_m = &mysql.Migrator{}
	case "PGSQL":
		_m = &pgsql.Migrator{}
	}
	return _m
}

func getViewCreator() viewCreator {
	dbDriver := config.GetConfig("DB_DRIVER")
	var c viewCreator
	switch dbDriver {
	case "MYSQL":
		return nil
	case "PGSQL":
		c = &pgsql.ViewCreator{}
	}
	return c
}

func CreateTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	fiasTableName := config.GetConfig("DB_OBJECTS_TABLE")
	fiasObjectTypesTableName := config.GetConfig("DB_OBJECT_TYPES_TABLE")
	fiasHierarchyTableName := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	fiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	_, tableCheck := dbInstance.Table(fiasTableName).Limit(1).Get()
	if tableCheck == nil &&
		shutdown.IsReboot &&
		config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT") == "original" {
		return nil
	}

	if tableCheck == nil {
		config.SetConfig("DB_ORIGINAL_OBJECTS_TABLE", config.GetConfig("DB_OBJECTS_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE", config.GetConfig("DB_OBJECT_TYPES_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE", config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE", config.GetConfig("DB_OBJECTS_KLADR_TABLE"))

		fiasTableName = config.GetConfig("DB_OBJECTS_TABLE") + "_temp"
		fiasObjectTypesTableName = config.GetConfig("DB_OBJECT_TYPES_TABLE") + "_temp"
		fiasHierarchyTableName = config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE") + "_temp"
		fiasKladrTableName = config.GetConfig("DB_OBJECTS_KLADR_TABLE") + "_temp"
		_, tempTableCheck := dbInstance.Table(fiasTableName).Limit(1).Get()

		config.SetConfig("DB_OBJECTS_TABLE", fiasTableName)
		config.SetConfig("DB_OBJECT_TYPES_TABLE", fiasObjectTypesTableName)
		config.SetConfig("DB_OBJECTS_HIERARCHY_TABLE", fiasHierarchyTableName)
		config.SetConfig("DB_OBJECTS_KLADR_TABLE", fiasKladrTableName)
		config.SetConfig("DB_TABLE_TYPES_FOR_IMPORT", "temp")
		if tempTableCheck == nil && !shutdown.IsReboot {
			return errors.New("fias tables and temp tables is exists")
		} else if tempTableCheck == nil {
			return nil
		}
	}

	return createFiasTables(
		fiasTableName,
		fiasObjectTypesTableName,
		fiasHierarchyTableName,
		fiasKladrTableName,
	)
}

func CreateIndexes() error {
	m := getMigrator()
	return m.CreateIndexes()
}

func MigrateDataFromTempTables() error {
	if config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT") != "temp" {
		return nil
	}

	m := getMigrator()
	return m.MigrateFromTempTables()
}

func CreateAdditionalViews() error {
	c := getViewCreator()
	if c == nil {
		return nil
	}

	err := c.CreateSettlementsView()
	if err != nil {
		return err
	}

	return c.CreateSettlementsParentsView()
}

func createFiasTables(
	fiasTableName string,
	fiasObjectTypesTableName string,
	fiasHierarchyTableName string,
	fiasKladrTableName string,
) error {
	m := getMigrator()

	if m == nil {
		return errors.New("error when get migrator instance")
	}
	m.SetObjectsTable(fiasTableName)
	m.SetObjectTypesTable(fiasObjectTypesTableName)
	m.SetHierarchyTable(fiasHierarchyTableName)
	m.SetKladrTable(fiasKladrTableName)

	err := m.ObjectsTableCreate()
	if err != nil {
		return err
	}
	err = m.ObjectTypesTableCreate()
	if err != nil {
		return err
	}
	err = m.HierarchyTableCreate()
	if err != nil {
		return err
	}
	return m.KladrTableCreate()
}
