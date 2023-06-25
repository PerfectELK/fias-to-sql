package dump

const (
	IMPORT_DUMP_FILENAME = "dump.json"
)

type Dump struct {
	ArchivePath string
	Files       []string
}

func MakeDump() error {
	return nil
}
