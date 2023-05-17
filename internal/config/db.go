package config

func dbConfig(m map[string]string) {
	m["DB_DRIVER"] = "MYSQL"
	m["DB_HOST"] = "127.0.0.1"
	m["DB_PORT"] = "3306"
	m["DB_NAME"] = "fias"
	m["DB_USER"] = "root"
	m["DB_PASSWORD"] = "123"
	m["DB_OBJECTS_TABLE"] = "fias_objects"
	m["DB_OBJECTS_HIERARCHY_TABLE"] = "fias_objects_hierarchy"
	m["DB_OBJECTS_KLADR_TABLE"] = "fias_object_kladr"
	m["DB_USE_TEMP_TABLE"] = "false"
}
