package database

import (
	"fmt"
	"os"
	"testing"
)

func Test_SqliteHandler(t *testing.T) {
	tester := NewSqliteDb("test.db")
	if err := tester.Execute("DROP TABLE IF EXISTS test"); err != nil {
		t.Error(err)
	}
	if err := tester.Execute("CREATE TABLE test (id integer, currency varchar(42), rate integer)"); err != nil {
		t.Error(err)
	}
	if err := tester.Execute("INSERT INTO test (id, currency, rate) VALUES (3, 'naira', 600)"); err != nil {
		t.Error(err)
	}
	row := tester.Query("SELECT id, currency, rate FROM test LIMIT 1")
	var (
		id       int
		currency string
		rate     int
	)

	row.Next()
	if err := row.Scan(&id, &currency, &rate); err != nil {
		t.Error(err)
	}
	if id != 3 {
		fmt.Println(id)
		t.Error()
	}
	if err := os.Remove("test.db"); err != nil {
		t.Error(err)
	}
}
