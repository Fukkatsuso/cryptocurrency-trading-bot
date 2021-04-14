package controller

import (
	"fmt"
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
)

func FetchTickerHandler(w http.ResponseWriter, r *http.Request) {
	client := bitflyer.NewClient(config.APIKey, config.APISecret)
	ticker, err := client.GetTicker(config.ProductCode)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to fetch ticker")
		return
	}

	// 時刻をeth_candlesに保存
	fmt.Println(ticker)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Fetch ticker")
}
