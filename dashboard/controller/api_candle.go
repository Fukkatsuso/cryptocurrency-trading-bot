package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	candles, _ := model.GetAllCandle(config.DB, config.CandleTableName, config.TimeFormat,
		config.ProductCode, config.CandleDuration, limit)

	df := model.DataFrame{
		ProductCode: config.ProductCode,
		Candles:     candles,
	}

	tradeParams := model.TradeParams{
		ProductCode: config.ProductCode,
		Size:        0.01,
	}

	// 移動平均線
	sma := r.URL.Query().Get("sma")
	if sma == "true" {
		period1 := getQueryUintDefault(r, "smaPeriod1", 7)
		period2 := getQueryUintDefault(r, "smaPeriod2", 14)
		period3 := getQueryUintDefault(r, "smaPeriod3", 50)
		enable := df.AddSMA(period1) && df.AddSMA(period2) && df.AddSMA(period3)
		tradeParams.SMAEnable = enable
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
		enable := df.AddEMA(period1) && df.AddEMA(period2) && df.AddEMA(period3)
		tradeParams.EMAEnable = enable
		tradeParams.EMAPeriod1 = period1
		tradeParams.EMAPeriod2 = period2
		tradeParams.EMAPeriod3 = period3
	}

	// ボリンジャーバンド
	bbands := r.URL.Query().Get("bbands")
	if bbands == "true" {
		n := getQueryUintDefault(r, "bbandsN", 20)
		k := getQueryUintDefault(r, "bbandsK", 2)
		enable := df.AddBBands(n, float64(k))
		tradeParams.BBandsEnable = enable
		tradeParams.BBandsN = n
		tradeParams.BBandsK = float64(k)
	}

	// 一目均衡表
	ichimoku := r.URL.Query().Get("ichimoku")
	if ichimoku == "true" {
		enable := df.AddIchimoku()
		tradeParams.IchimokuEnable = enable
	}

	// RSI(Relative Strength Index, 相対力指数)
	rsi := r.URL.Query().Get("rsi")
	if rsi == "true" {
		period := getQueryUintDefault(r, "rsiPeriod", 14)
		enable := df.AddRSI(period)
		tradeParams.RSIEnable = enable
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
		enable := df.AddMACD(period1, period2, period3)
		tradeParams.MACDEnable = enable
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
