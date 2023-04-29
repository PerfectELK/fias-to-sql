package config

func dbConfig(m map[string]string) {
	m["DB_DRIVER"] = "MYSQL"
	m["DB_HOST"] = ""
	m["DB_PORT"] = ""
	m["DB_NAME"] = ""
	m["DB_USER"] = ""
	m["DB_PASSWORD"] = ""
	m["DB_OBJECTS_TABLE"] = "fias_objects"
	m["DB_OBJECTS_HIERARCHY_TABLE"] = "fias_objects_hierarchy"
}
