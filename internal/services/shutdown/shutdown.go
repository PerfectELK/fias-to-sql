package shutdown

import (
	"encoding/json"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/slice"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	IMPORT_DUMP_FILENAME = "dump.json"
)

var archivePath string
var files []DumpFile

var archivePathToDump string
var filesToDump []DumpFile

var IsReboot bool

var c chan os.Signal

type Dump struct {
	ArchivePath string     `json:"archive_path"`
	TablesType  string     `json:"tables_type"`
	Files       []DumpFile `json:"files"`
}

type DumpFile struct {
	FileName      string `json:"file_name"`
	RecordsAmount int64  `json:"records_amount"`
}

func PutFileToDump(fileName DumpFile) {
	filesToDump = append(filesToDump, fileName)
}

func GetFiles() []DumpFile {
	return files
}

func GetFilesNames() []string {
	return slice.Map(files, func(file DumpFile) string {
		return file.FileName
	})
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
	dumpFile := filepath.Join(os.Getenv("APP_ROOT"), "storage", IMPORT_DUMP_FILENAME)

	var dump Dump
	dump.ArchivePath = archivePathToDump
	dump.TablesType = config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT")
	dump.Files = filesToDump

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
	if _, err := os.Stat(filepath.Join(os.Getenv("APP_ROOT"), "storage", IMPORT_DUMP_FILENAME)); err != nil {
		return false
	}
	IsReboot = true
	return true
}

func RebootAfterGracefulShutdown() error {
	dumpFile := filepath.Join(os.Getenv("APP_ROOT"), "storage", IMPORT_DUMP_FILENAME)
	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return err
	}

	var dump Dump
	err = json.Unmarshal(data, &dump)
	if err != nil {
		return err
	}

	archivePath = dump.ArchivePath
	files = dump.Files
	config.SetConfig("DB_TABLE_TYPES_FOR_IMPORT", dump.TablesType)

	err = os.Remove(dumpFile)
	if err != nil {
		return err
	}

	return nil
}
