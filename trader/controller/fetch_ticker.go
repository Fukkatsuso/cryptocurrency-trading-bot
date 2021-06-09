package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/go-lib/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/go-lib/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
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

	fmt.Println("[fetchTicker]", *ticker)

	// 時刻をeth_candlesに保存
	err = model.CreateCandleWithDuration(config.DB, config.CandleTableName, config.TimeFormat, config.LocalTime, config.TradeHour,
		ticker, config.ProductCode, 24*time.Hour)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create candle")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Fetch ticker")
}
