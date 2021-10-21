package persistence_test

import (
	"database/sql"
)

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
