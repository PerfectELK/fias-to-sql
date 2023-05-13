package mysql

import (
	"fias_to_sql/internal/config"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestConnect(t *testing.T) {
	p := Processor{}

	err := p.Connect()
	fmt.Println(p.Query("SELECT * FROM table"))
	if err != nil {
		t.Error(err)
	}
}
