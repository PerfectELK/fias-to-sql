package pgsql

import "fias_to_sql/pkg/db"

func ObjectsTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
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
	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_name_index ON " + tableName + " (name);" +
			" CREATE INDEX " + tableName + "_object_guid_index ON " + tableName + " (object_guid);" +
			" CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_type_name_index ON " + tableName + " (type_name);",
	)
}

func HierarchyTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"parent_object_id INT NOT NULL DEFAULT 0);",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_parent_object_id_index ON " + tableName + " (parent_object_id);",
	)
}

func KladrTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	err = dbInstance.Exec(
		"CREATE TABLE " + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"kladr_id VARCHAR(50) NOT NULL DEFAULT '');",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_kladr_id_index ON " + tableName + " (kladr_id);",
	)
}
