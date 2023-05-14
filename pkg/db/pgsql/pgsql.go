package pgsql

import (
	"database/sql"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db/abstract"
	"fias_to_sql/pkg/db/helpers"
	"fias_to_sql/pkg/db/types"
	"fmt"
	_ "github.com/lib/pq"
)

type Processor struct {
	abstract.DbProcessor
	db          *sql.DB
	isConnected bool
	table       string
	sel         []string
	where       [][]string
}

func (m *Processor) Connect(dbName ...string) error {
	connectStr := "postgres://" + config.GetConfig("DB_USER") + ":" + config.GetConfig("DB_PASSWORD") + "@" + config.GetConfig("DB_HOST") + ":" + config.GetConfig("DB_PORT") + "/"
	if len(dbName) > 0 {
		connectStr += dbName[0]
	}
	connectStr += "?sslmode=disable"
	db, err := sql.Open("postgres", connectStr)
	err = db.Ping()
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
	rows, err := m.db.Query(q)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (m *Processor) Insert(table string, mm map[string]string) error {
	queryStr := "INSERT INTO " + table
	var keys []string
	var values []string
	for key, elem := range mm {
		if elem == "" {
			continue
		}
		keys = append(keys, key)
		values = append(values, elem)
	}

	keysStr := ""
	valuesStr := ""
	for key, _ := range keys {
		afterStr := ""
		if key != len(keys)-1 {
			afterStr += ", "
		}
		keysStr += keys[key] + afterStr
		valuesStr += "\"" + values[key] + "\"" + afterStr
	}

	queryStr += " (" + keysStr + ") VALUES (" + valuesStr + ");"
	return m.Exec(queryStr)
}

func (m *Processor) InsertList(table string, keys []types.Key, values [][]string) error {
	queryStr := "INSERT INTO " + table

	keysStr := ""
	valuesStr := ""
	for i, val := range keys {
		afterStr := ""
		if i != len(keys)-1 {
			afterStr += ", "
		}
		keysStr += val.Name + afterStr
	}
	queryStr += " (" + keysStr + ") "

	queryCount := 0
	for i, vals := range values {
		queryCount++
		valuesStr += "( "
		for key, val := range vals {
			afterStr := ""
			if key != len(vals)-1 {
				afterStr += ", "
			}
			valuesStr += "\"" + helpers.SqlRealEscapeString(val) + "\"" + afterStr
		}
		closeStr := ") "
		if i != len(values)-1 && queryCount < 4000 {
			closeStr += ", "
		}
		valuesStr += closeStr
		if queryCount >= 4000 {
			q := queryStr + "VALUES " + valuesStr + ";"
			err := m.Exec(q)
			valuesStr = ""
			if err != nil {
				fmt.Println(err)
				return err
			}
			queryCount = 0
		}
	}

	if valuesStr != "" {
		return m.Exec(queryStr + "VALUES " + valuesStr + ";")
	}
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

func (m *Processor) Get() (map[string]string, error) {
	return nil, nil
}

func (m *Processor) Use(q string) error {
	m.db.Close()
	m.isConnected = false
	return m.Connect(q)
}

func (m *Processor) Query(q string) (*sql.Rows, error) {
	return m.db.Query(q)
}

func (m *Processor) GetDriverName() string {
	return "PGSQL"
}
