package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var (
	MYSQL_USER            = os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD        = os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST            = os.Getenv("MYSQL_HOST")
	MYSQL_PORT            = os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE        = os.Getenv("MYSQL_DATABASE")
	MYSQL_CONNECTION_NAME = os.Getenv("MYSQL_CONNECTION_NAME")
)

var timeFormat = "2006-01-02 15:04:05"

func init() {
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dsn string
	if MYSQL_CONNECTION_NAME == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE)
	} else {
		dsn = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s", MYSQL_USER, MYSQL_PASSWORD, socketDir, MYSQL_CONNECTION_NAME, MYSQL_DATABASE)
	}
	fmt.Println("dsn:", dsn)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// DBに接続してeth_candlesテーブルの情報(レコードのcount?)を取得，何か値を表示させる
	cmd := `SELECT COUNT(time) FROM eth_candles`
	row := db.QueryRow(cmd)
	var count int
	if err := row.Scan(&count); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s! (trader)\ncount: %d\n", name, count)
}

func saveCandle(tm time.Time) {
	cmd := `INSERT INTO eth_candles (time, open, close, high, low, volume) VALUES (?, 1.0, 2.0, 3.0, 4.0, 5.0)`
	_, err := db.Exec(cmd, tm.Format(timeFormat))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// 3秒おきにDBのデータを更新(適当な値を挿入する)
	// あくまでもDB接続のテストであり，定期実行はGoでは扱わないため，動けばOK
	go func() {
		t := time.NewTicker(3 * time.Second)
		defer t.Stop()
		for {
			select {
			case now := <-t.C:
				fmt.Println(now.Format(timeFormat))
				saveCandle(now)
			}
		}
	}()

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
