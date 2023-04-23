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
}

func (m *Processor) Connect() error {
	connectStr := config.GetConfig("DB_USER") + ":" + config.GetConfig("DB_PASSWORD") + "tcp(" + config.GetConfig("DB_HOST") + ":" + config.GetConfig("DB_PORT") + ")/" + config.GetConfig("DB_NAME")
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
	return nil
}

func (m *Processor) Insert(mm map[string]string) error {
	return nil
}

func (m *Processor) Query(q [][]string) struct{} {
	return struct{}{}
}

func (m *Processor) IsConnected() bool {
	return m.isConnected
}
