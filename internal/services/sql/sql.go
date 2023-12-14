package sql

import (
	"io"
	"os"
	"path/filepath"
)

func LoadSql(name string) (string, error) {
	file, err := os.OpenFile(filepath.Join(os.Getenv("APP_ROOT"), "sql", name), os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
