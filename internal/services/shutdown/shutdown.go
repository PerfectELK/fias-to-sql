package shutdown

import (
	"encoding/json"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/fias/types"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var archivePath string
var fileNames []string

var archivePathToDump string
var fileNamesToDump []string

var IsReboot bool

var c chan os.Signal

func PutFileNameToDump(fileName string) {
	fileNamesToDump = append(fileNamesToDump, fileName)
}

func GetFileNames() []string {
	return fileNames
}

func GetArchivePath() string {
	return archivePath
}

func SetArchivePathToDump(newArchivePath string) {
	archivePathToDump = newArchivePath
}

func OnShutdown(fn func()) {
	c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fn()
	}()
}

func MakeDump() error {
	dumpFile := filepath.Join(os.Getenv("APP_ROOT"), "storage", types.IMPORT_DUMP_FILENAME)

	var dump types.Dump
	dump.ArchivePath = archivePathToDump
	dump.TablesType = config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT")
	dump.Files = fileNamesToDump

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
	IsReboot = true
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

	archivePath = dump.ArchivePath
	fileNames = dump.Files
	config.SetConfig("DB_TABLE_TYPES_FOR_IMPORT", dump.TablesType)

	err = os.Remove(dumpFile)
	if err != nil {
		return err
	}

	return nil
}
