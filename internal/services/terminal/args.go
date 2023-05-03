package terminal

import (
	"fias_to_sql/internal/config"
	"flag"
)

func ParseArgs() error {
	var importDestination string
	var dbDriver string
	var dbHost string
	var dbPort string
	var dbName string
	var dbUser string
	var dbPassword string
	var objectsTableName string
	var objectsHierarchyTableName string

	flag.StringVar(&importDestination, "import-destination", "", "")
	flag.StringVar(&dbDriver, "db-driver", "", "")
	flag.StringVar(&dbHost, "db-host", "", "")
	flag.StringVar(&dbPort, "db-port", "", "")
	flag.StringVar(&dbName, "db-name", "", "")
	flag.StringVar(&dbUser, "db-user", "", "")
	flag.StringVar(&dbPassword, "db-password", "", "")
	flag.StringVar(&objectsTableName, "objects-table", "", "")
	flag.StringVar(&objectsHierarchyTableName, "objects-hierarchy-table", "", "")

	if importDestination != "" {
		config.SetConfig("IMPORT_DESTINATION", importDestination)
	}
	if dbDriver != "" {
		config.SetConfig("DB_DRIVER", dbDriver)
	}
	if dbHost != "" {
		config.SetConfig("DB_HOST", dbHost)
	}
	if dbPort != "" {
		config.SetConfig("DB_PORT", dbPort)
	}
	if dbName != "" {
		config.SetConfig("DB_NAME", dbName)
	}
	if dbUser != "" {
		config.SetConfig("DB_USER", dbUser)
	}
	if dbPassword != "" {
		config.SetConfig("DB_PASSWORD", dbPassword)
	}
	if objectsTableName != "" {
		config.SetConfig("DB_OBJECTS_TABLE", objectsTableName)
	}
	if objectsHierarchyTableName != "" {
		config.SetConfig("DB_OBJECTS_HIERARCHY_TABLE", objectsHierarchyTableName)
	}

	return nil
}
