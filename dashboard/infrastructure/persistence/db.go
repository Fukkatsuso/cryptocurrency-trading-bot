package persistence

import (
	"database/sql"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func NewMySQLTransaction(dsn string) *sql.Tx {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	tx, _ := db.Begin()

	return tx
}
