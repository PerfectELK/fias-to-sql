package pgsql

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
	"fmt"
)

type Migrator struct {
	objectsTable     string
	objectTypesTable string
	hierarchyTable   string
	kladrTable       string
}

func (p *Migrator) SetObjectsTable(table string) {
	p.objectsTable = table
}

func (p *Migrator) SetObjectTypesTable(table string) {
	p.objectTypesTable = table
}

func (p *Migrator) SetHierarchyTable(table string) {
	p.hierarchyTable = table
}

func (p *Migrator) SetKladrTable(table string) {
	p.kladrTable = table
}

func (p *Migrator) ObjectsTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + p.objectsTable + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INTEGER NOT NULL DEFAULT 0," +
			"object_guid VARCHAR(100) NOT NULL DEFAULT ''," +
			"type_name VARCHAR(100) NOT NULL DEFAULT ''," +
			"level INT NOT NULL DEFAULT 0," +
			"name VARCHAR(255) NOT NULL DEFAULT '');",
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Migrator) ObjectTypesTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	return dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + p.objectsTable + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"level INT NOT NULL DEFAULT 0," +
			"short_name VARCHAR(255) NOT NULL DEFAULT ''," +
			"name VARCHAR(255) NOT NULL DEFAULT '');",
	)
}

func (p *Migrator) HierarchyTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + p.hierarchyTable + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"parent_object_id INT NOT NULL DEFAULT 0);",
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Migrator) KladrTableCreate() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + p.kladrTable + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"kladr_id VARCHAR(50) NOT NULL DEFAULT '');",
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Migrator) CreateIndexes() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")

	err = dbInstance.Exec(
		"CREATE INDEX " + p.kladrTable + "_object_id_index ON " + dbSchema + "." + p.kladrTable + " (object_id);" +
			" CREATE INDEX " + p.kladrTable + "_kladr_id_index ON " + dbSchema + "." + p.kladrTable + " (kladr_id);",
	)
	if err != nil {
		return err
	}

	err = dbInstance.Exec(
		"CREATE INDEX " + p.hierarchyTable + "_object_id_index ON " + dbSchema + "." + p.hierarchyTable + " (object_id);" +
			" CREATE INDEX " + p.hierarchyTable + "_parent_object_id_index ON " + dbSchema + "." + p.hierarchyTable + " (parent_object_id);",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + p.objectsTable + "_name_index ON " + dbSchema + "." + p.objectsTable + " (name);" +
			" CREATE INDEX " + p.objectsTable + "_object_guid_index ON " + dbSchema + "." + p.objectsTable + " (object_guid);" +
			" CREATE INDEX " + p.objectsTable + "_object_id_index ON " + dbSchema + "." + p.objectsTable + " (object_id);" +
			" CREATE INDEX " + p.objectsTable + "_type_name_index ON " + dbSchema + "." + p.objectsTable + " (type_name);",
	)
}

func (p *Migrator) MigrateFromTempTables() error {
	err := p.dropOldTables()
	if err != nil {
		return err
	}
	err = p.renameTables()
	if err != nil {
		return err
	}
	return p.renameIndexes()
}

func (p *Migrator) dropOldTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbSchema := config.GetConfig("DB_SCHEMA")

	originalObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	originalObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE"))
	originalHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE"))
	originalFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))

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

func (p *Migrator) renameTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbSchema := config.GetConfig("DB_SCHEMA")

	originalObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	originalObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE"))
	originalHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE"))
	originalFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))

	tempObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_TABLE"))
	tempObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECT_TYPES_TABLE"))
	tempHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
	tempFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_KLADR_TABLE"))

	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempObjectsTable, originalObjectsTable))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempObjectTypesTableName, originalObjectTypesTableName))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempHierarchyObjectsTable, originalHierarchyObjectsTable))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempFiasKladrTableName, originalFiasKladrTableName))
	return err
}

func (p *Migrator) renameIndexes() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	originalObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE")
	originalHierarchyObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE")
	originalFiasKladrTableName := config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE")

	tempObjectsTable := config.GetConfig("DB_OBJECTS_TABLE")
	tempHierarchyObjectsTable := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	tempFiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_name_index RENAME TO %s_name_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_guid_index RENAME TO %s_object_guid_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_type_name_index RENAME TO %s_type_name_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempHierarchyObjectsTable,
		originalHierarchyObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_parent_object_id_index RENAME TO %s_parent_object_id_index",
		tempHierarchyObjectsTable,
		originalHierarchyObjectsTable,
	))
	if err != nil {
		return err
	}

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempFiasKladrTableName,
		originalFiasKladrTableName,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_kladr_id_index RENAME TO %s_kladr_id_index",
		tempFiasKladrTableName,
		originalFiasKladrTableName,
	))
	return err
}

type ViewCreator struct{}

func (v *ViewCreator) CreateSettlementsParentsView() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")

	ObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	FiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))
	HierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE"))

	query := fmt.Sprintf("CREATE MATERIALIZED VIEW settlements_parents AS "+
		"SELECT fias.id, "+
		"fias.settlement_id, "+
		"fias.parent_id "+
		"FROM (WITH cities AS (SELECT %s.object_id "+
		"FROM %s "+
		"JOIN %s ON %s.object_id = %s.object_id "+
		"WHERE (level < 6 OR type_name IN "+
		"('г', 'г.', 'пгт', 'пгт.', 'Респ', 'обл', 'обл.', 'Аобл', 'а.обл.', 'а.окр.', "+
		"'АО', 'г.ф.з.'))) "+
		"SELECT %s.id, %s.object_id AS settlement_id, parent_object_id AS parent_id "+
		"FROM %s "+
		"JOIN cities AS c1 ON c1.object_id = %s.object_id "+
		"JOIN cities AS c2 ON c2.object_id = %s.parent_object_id) AS fias;",
		ObjectsTable,
		ObjectsTable,
		FiasKladrTableName,
		ObjectsTable,
		FiasKladrTableName,
		HierarchyObjectsTable,
		HierarchyObjectsTable,
		HierarchyObjectsTable,
		HierarchyObjectsTable,
		HierarchyObjectsTable,
	)

	dbInstance.Exec("DROP MATERIALIZED VIEW IF EXISTS settlements_parents")
	err = dbInstance.Exec(query)
	if err != nil {
		return err
	}

	query = "create index settlements_parents_id " +
		"on settlements_parents (id); " +
		"create index settlements_parents_settlement_id " +
		"on settlements_parents (settlement_id); " +
		"create index settlements_parents_parent_id " +
		"on settlements_parents (parent_id);"

	return dbInstance.Exec(query)
}

func (v *ViewCreator) CreateSettlementsView() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")

	ObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	ObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE"))
	FiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))

	query := fmt.Sprintf("CREATE MATERIALIZED VIEW settlements AS "+
		"SELECT fias.id, "+
		"fias.fias_id, "+
		"fias.kladr_id, "+
		"fias.type, "+
		"fias.type_short, "+
		"fias.name, "+
		"fias.created_at "+
		"FROM (SELECT %s.object_id                                            as id, "+
		"object_guid                                                                as fias_id, "+
		"%s.kladr_id,"+
		"replace(LOWER(%s.name), '.', '')                            as type, "+
		"replace(LOWER(type_name), '.', '')                                             as type_short, "+
		"%s.name, "+
		"to_char(now(), 'YYYY-MM-DD HH12:MI:SS'::text)::timestamp without time zone AS created_at "+
		"FROM %s "+
		"JOIN %s ON %s.object_id = %s.object_id "+
		"LEFT JOIN %s ON "+
		"%s.type_name = %s.short_name AND %s.level = %s.level "+
		"WHERE %s.level < 6 "+
		"OR type_name IN "+
		"('г', 'г.', 'пгт', 'пгт.', 'Респ', 'обл', 'обл.', 'Аобл', 'а.обл.', 'а.окр.', 'АО', 'г.ф.з.')) AS fias;",
		ObjectsTable,
		FiasKladrTableName,
		ObjectTypesTableName,
		ObjectsTable,
		ObjectsTable,
		FiasKladrTableName,
		ObjectsTable,
		FiasKladrTableName,
		ObjectTypesTableName,
		ObjectsTable,
		ObjectTypesTableName,
		ObjectsTable,
		ObjectTypesTableName,
		ObjectsTable,
	)

	dbInstance.Exec("DROP MATERIALIZED VIEW IF EXISTS settlements")
	err = dbInstance.Exec(query)
	if err != nil {
		return err
	}

	query = "create index settlements_id " +
		"on settlements (id);" +
		"create index settlements_fias_id " +
		"on settlements (fias_id); " +
		"create index settlements_kladr_id " +
		"on settlements (kladr_id); " +
		"create index settlements_type_short " +
		"on settlements (type_short);"

	return dbInstance.Exec(query)
}
