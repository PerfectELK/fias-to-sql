package fias

import (
	"encoding/json"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/fias/types"
	"os"
	"path/filepath"
)

func MakeDump() error {
	dumpFile := filepath.Join(os.Getenv("APP_ROOT"), "storage", types.IMPORT_DUMP_FILENAME)

	var dump types.Dump
	dump.ArchivePath = ArchivePathToDump
	dump.TablesType = config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT")
	dump.Files = FileNamesToDump

	b, err := json.Marshal(dump)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(dumpFile, os.O_CREATE|os.O_WRONLY, 0644)
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func CheckGracefulShutdown() bool {
	if _, err := os.Stat(filepath.Join(os.Getenv("APP_ROOT"), "storage", types.IMPORT_DUMP_FILENAME)); err != nil {
		return false
	}
	return true
}

func RebootAfterGracefulShutdown() error {
	dumpFile := filepath.Join(os.Getenv("APP_ROOT"), "storage", types.IMPORT_DUMP_FILENAME)
	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return err
	}

	var dump types.Dump
	err = json.Unmarshal(data, &dump)
	if err != nil {
		return err
	}

	ArchivePath = dump.ArchivePath
	FileNames = dump.Files
	config.SetConfig("DB_TABLE_TYPES_FOR_IMPORT", dump.TablesType)

	err = os.Remove(dumpFile)
	if err != nil {
		return err
	}

	return nil
}
