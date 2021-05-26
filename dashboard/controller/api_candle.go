package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/model"
)

// URLクエリパラメータから非負整数を取り出す
// エラーが生じたらデフォルトの値defを返す
func getQueryUintDefault(r *http.Request, query string, def int) int {
	strVal := r.URL.Query().Get(query)
	intVal, err := strconv.Atoi(strVal)
	if strVal == "" || err != nil || intVal < 0 {
		intVal = def
	}
	return intVal
}

// candleデータとテクニカル指標をjsonで返す
func APICandleHandler(w http.ResponseWriter, r *http.Request) {
	// limit: candle最大数．[0, 1000]の範囲に限定
	limit := getQueryUintDefault(r, "limit", 1000)
	if limit > 1000 {
		limit = 1000
	}

	candles, _ := model.GetAllCandle(config.ProductCode, 24*time.Hour, limit)

	df := model.DataFrame{
		ProductCode: config.ProductCode,
		Candles:     candles,
	}

	tradeParams := model.TradeParams{}

	// 移動平均線
	sma := r.URL.Query().Get("sma")
	if sma == "true" {
		period1 := getQueryUintDefault(r, "smaPeriod1", 7)
		period2 := getQueryUintDefault(r, "smaPeriod2", 14)
		period3 := getQueryUintDefault(r, "smaPeriod3", 50)
		df.AddSMA(period1)
		df.AddSMA(period2)
		df.AddSMA(period3)
		tradeParams.SMAEnable = true
		tradeParams.SMAPeriod1 = period1
		tradeParams.SMAPeriod2 = period2
		tradeParams.SMAPeriod3 = period3
	}

	// 指数平滑移動平均線
	ema := r.URL.Query().Get("ema")
	if ema == "true" {
		period1 := getQueryUintDefault(r, "emaPeriod1", 7)
		period2 := getQueryUintDefault(r, "emaPeriod2", 14)
		period3 := getQueryUintDefault(r, "emaPeriod3", 50)
		df.AddEMA(period1)
		df.AddEMA(period2)
		df.AddEMA(period3)
		tradeParams.EMAEnable = true
		tradeParams.EMAPeriod1 = period1
		tradeParams.EMAPeriod2 = period2
		tradeParams.EMAPeriod3 = period3
	}

	// ボリンジャーバンド
	bbands := r.URL.Query().Get("bbands")
	if bbands == "true" {
		n := getQueryUintDefault(r, "bbandsN", 20)
		k := getQueryUintDefault(r, "bbandsK", 2)
		df.AddBBands(n, float64(k))
		tradeParams.BBandsEnable = true
		tradeParams.BBandsN = n
		tradeParams.BBandsK = float64(k)
	}

	// 一目均衡表
	ichimoku := r.URL.Query().Get("ichimoku")
	if ichimoku == "true" {
		df.AddIchimoku()
		tradeParams.IchimokuEnable = true
	}

	// RSI(Relative Strength Index, 相対力指数)
	rsi := r.URL.Query().Get("rsi")
	if rsi == "true" {
		period := getQueryUintDefault(r, "rsiPeriod", 14)
		df.AddRSI(period)
		tradeParams.RSIEnable = true
		tradeParams.RSIPeriod = period
		tradeParams.RSIBuyThread = 30.0
		tradeParams.RSISellThread = 70.0
	}

	// MACD(Moving Average Convergence/Divergence, 移動平均・収束拡散)
	macd := r.URL.Query().Get("macd")
	if macd == "true" {
		period1 := getQueryUintDefault(r, "macdPeriod1", 12)
		period2 := getQueryUintDefault(r, "macdPeriod2", 26)
		period3 := getQueryUintDefault(r, "macdPeriod3", 9)
		df.AddMACD(period1, period2, period3)
		tradeParams.MACDEnable = true
		tradeParams.MACDFastPeriod = period1
		tradeParams.MACDSlowPeriod = period2
		tradeParams.MACDSignalPeriod = period3
	}

	backtest := r.URL.Query().Get("backtest")
	if backtest == "true" {
		df.BackTest(&tradeParams)
	}

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
