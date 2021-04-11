package controller

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
)

// eth_candlesテーブルの情報(レコードのcount)を取得して表示する
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	cmd := `SELECT COUNT(time) FROM eth_candles`
	row := config.DB.QueryRow(cmd)
	var count int
	if err := row.Scan(&count); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

    fmt.Fprintf(w, "Hello World! (trader)\neth_candles count: %d\n", count)
}
