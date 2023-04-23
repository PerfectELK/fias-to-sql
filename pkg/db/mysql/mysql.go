package mysql

import (
	"database/sql"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db/abstract"
	_ "github.com/go-sql-driver/mysql"
)

type Processor struct {
	abstract.DbProcessor
	db          *sql.DB
	isConnected bool
	table       string
	sel         []string
	where       [][]string
}

func (m *Processor) Connect() error {
	connectStr := config.GetConfig("DB_USER") + ":" + config.GetConfig("DB_PASSWORD") + "@tcp(" + config.GetConfig("DB_HOST") + ":" + config.GetConfig("DB_PORT") + ")/"
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return err
	}
	m.db = db
	m.isConnected = true
	return nil
}

func (m *Processor) Disconnect() error {
	m.isConnected = false
	return m.db.Close()
}

func (m *Processor) Exec(q string) error {
	_, err := m.db.Query(q)
	return err
}

func (m *Processor) Insert(mm map[string]string) error {
	return nil
}

func (m *Processor) IsConnected() bool {
	return m.isConnected
}

func (m *Processor) Where(q [][]string) abstract.DbProcessor {
	m.where = q
	return m
}

func (m *Processor) Table(t string) abstract.DbProcessor {
	m.table = t
	return m
}

func (m *Processor) Select(s []string) abstract.DbProcessor {
	m.sel = s
	return m
}

func (m *Processor) Get() map[string]string {
	return nil
}

func (m *Processor) Use(q string) error {
	return m.Exec("USE " + q)
}
