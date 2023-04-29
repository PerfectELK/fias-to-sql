package abstract

import (
	"database/sql"
	"fias_to_sql/pkg/db/types"
)

type DbProcessor interface {
	Connect(...string) error
	Disconnect() error
	Use(q string) error
	Exec(q string) error
	Insert(table string, m map[string]string) error
	InsertList(table string, keys []types.Key, values [][]string) error
	Table(t string) DbProcessor
	Select(s []string) DbProcessor
	Where(q [][]string) DbProcessor
	Get() (map[string]string, error)
	IsConnected() bool
	Query(string) (*sql.Rows, error)
}
