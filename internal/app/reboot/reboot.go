package reboot

import (
	"fias_to_sql/internal/app/dump"
	"os"
	"path/filepath"
)

func CheckSoftTerminate() bool {
	if _, err := os.Stat(filepath.Join(os.Getenv("APP_ROOT"), "storage", dump.IMPORT_DUMP_FILENAME)); err != nil {
		return false
	}
	return true
}

func RebootAfterSoftTerminate() error {
	return nil
}
