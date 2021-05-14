package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/model"
)

// candleデータとテクニカル指標をjsonで返す
func APICandleHandler(w http.ResponseWriter, r *http.Request) {
	// limit: candle最大数．[0, 1000]の範囲に限定
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	candles, _ := model.GetAllCandle(config.ProductCode, 24*time.Hour, limit)

	df := model.DataFrame{
		Candles: candles,
	}

	df.AddSMA(3)
	df.AddSMA(7)
	df.AddSMA(14)
	df.AddEMA(3)
	df.AddEMA(7)
	df.AddEMA(14)
	df.AddBBands(3, float64(2))
	df.AddIchimoku()
	df.AddRSI(7)
	df.AddMACD(3, 7, 5)

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
