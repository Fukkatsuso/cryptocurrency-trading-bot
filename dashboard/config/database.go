package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var (
	MYSQL_USER            = os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD        = os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST            = os.Getenv("MYSQL_HOST")
	MYSQL_PORT            = os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE        = os.Getenv("MYSQL_DATABASE")
	MYSQL_CONNECTION_NAME = os.Getenv("MYSQL_CONNECTION_NAME")
	MYSQL_OPTION          = "?parseTime=true"
)

const (
	CandleTableName     = "eth_candles"
	TradeParamTableName = "trade_params"
	TimeFormat          = "2006-01-02 15:04:05"
)

func DSN() string {
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dsn string
	if MYSQL_CONNECTION_NAME == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE, MYSQL_OPTION)
	} else {
		dsn = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s%s", MYSQL_USER, MYSQL_PASSWORD, socketDir, MYSQL_CONNECTION_NAME, MYSQL_DATABASE, MYSQL_OPTION)
	}

	return dsn
}

func init() {
	dsn := DSN()
	// テスト時に出力しないほうがよさげ
	// fmt.Println("dsn:", dsn)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
}
