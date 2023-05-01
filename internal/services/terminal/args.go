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

	return nil
}
