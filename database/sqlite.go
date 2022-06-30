package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type sqliteDb struct {
	DB *sql.DB
}

func NewSqliteDb(dbFilename string) *sqliteDb {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	sqliteDb := new(sqliteDb)
	sqliteDb.DB = db
	return sqliteDb
}

func (s *sqliteDb) Execute(statement string) (err error) {
	_, err = s.DB.Exec(statement)
	return err
}

func (s *sqliteDb) Query(statement string) Row {
	rows, err := s.DB.Query(statement)
	if err != nil {
		log.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (r SqliteRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}
