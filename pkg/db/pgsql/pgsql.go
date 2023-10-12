package pgsql

import (
	"database/sql"
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db/abstract"
	"fias_to_sql/pkg/db/helpers"
	"fias_to_sql/pkg/db/types"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"unicode/utf8"
)

type Processor struct {
	abstract.DbProcessor
	db          *sql.DB
	isConnected bool
	table       string
	sel         []string
	where       [][]string
	limit       int
	schema      string
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
	m.schema = config.GetConfig("DB_SCHEMA")
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
	queryStr := fmt.Sprintf("INSERT INTO %s.%s", m.schema, table)
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
		valuesStr += "'" + values[key] + "'" + afterStr
	}

	queryStr += " (" + keysStr + ") VALUES (" + valuesStr + ");"
	return m.Exec(queryStr)
}

const LEN = 20

func (m *Processor) InsertList(table string, keys []types.Key, values [][]string) error {
	querySB := strings.Builder{}
	querySB.Grow(len(keys) * LEN)

	fmt.Fprintf(&querySB, "INSERT INTO %s.%s", m.schema, table)
	keysStr := ""
	for i, val := range keys {
		afterStr := ""
		if i != len(keys)-1 {
			afterStr += ", "
		}
		keysStr += val.Name + afterStr
	}
	fmt.Fprintf(&querySB, " ( %s ) ", keysStr)

	valuesSB := strings.Builder{}
	valuesSB.Grow(len(keys) * len(values) * LEN)
	queryCount := 0
	for i, vals := range values {
		queryCount++
		valuesSB.WriteString("( ")
		for key, val := range vals {
			afterStr := ""
			if key != len(vals)-1 {
				afterStr += ", "
			}
			fmt.Fprintf(&valuesSB, "'%s'%s", helpers.PgsqlRealEscapeString(val), afterStr)
		}
		closeStr := ") "
		if i != len(values)-1 && queryCount < 4000 {
			closeStr += ", "
		}
		valuesSB.WriteString(closeStr)
		if queryCount >= 4000 {
			q := fmt.Sprintf("%sVALUES %s;", querySB.String(), valuesSB.String())
			err := m.Exec(q)
			valuesSB = strings.Builder{}
			valuesSB.Grow(len(keys) * len(values) * LEN)
			if err != nil {
				return err
			}
			queryCount = 0
		}
	}
	if utf8.RuneCountInString(valuesSB.String()) != 0 {
		q := fmt.Sprintf("%sVALUES %s;", querySB.String(), valuesSB.String())
		return m.Exec(q)
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

func (m *Processor) Limit(l int) abstract.DbProcessor {
	m.limit = l
	return m
}

func (m *Processor) Get() (*sql.Rows, error) {
	queryString := "SELECT "
	if len(m.sel) == 0 {
		queryString += " * "
	}
	for i, field := range m.sel {
		queryString += fmt.Sprintf("%s ", field)
		if i != len(m.sel)-1 {
			queryString += ", "
		}
	}
	queryString += fmt.Sprintf("FROM %s.%s ", m.schema, m.table)
	if len(m.where) > 0 {
		queryString += "WHERE "
		for i, whereCond := range m.where {
			queryString += fmt.Sprintf("%s %s %s", whereCond[0], whereCond[1], whereCond[2])
			if i != len(m.where)-1 {
				queryString += " AND "
			}
		}
	}

	if m.limit > 0 {
		queryString += fmt.Sprintf("LIMIT %d", m.limit)
	}

	m.where = [][]string{}
	m.sel = []string{}
	m.limit = 0

	rows, err := m.db.Query(queryString)
	return rows, err
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
