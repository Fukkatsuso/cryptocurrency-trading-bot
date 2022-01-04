package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

type DataFrameHandler interface {
	Get(productCode string, tradeSize float64) http.HandlerFunc
}

type dataFrameHandler struct {
	dataFrameUsecase usecase.DataFrameUsecase
}

func NewDataFrameHandler(du usecase.DataFrameUsecase) DataFrameHandler {
	return &dataFrameHandler{
		dataFrameUsecase: du,
	}
}

func (dh *dataFrameHandler) Get(productCode string, tradeSize float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := reqUrlToTradeParams(r, productCode, tradeSize)

		// [0, 1000]の範囲に限定
		candleLimit := getQueryUintDefault(r, "limit", 1000)
		if candleLimit > 1000 {
			candleLimit = 1000
		}

		backtestEnable := r.URL.Query().Get("backtest") == "true"

		df, err := dh.dataFrameUsecase.Get(params, int64(candleLimit), backtestEnable)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// jsonに変換できるように入れ替える
		dto := dto.ConvertDataFrame(df)

		js, err := json.Marshal(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func reqUrlToTradeParams(r *http.Request, productCode string, tradeSize float64) *model.TradeParams {
	sma := r.URL.Query().Get("sma")
	smaEnable := sma == "true"
	var smaPeriod1, smaPeriod2, smaPeriod3 int
	if smaEnable {
		smaPeriod1 = getQueryUintDefault(r, "smaPeriod1", 7)
		smaPeriod2 = getQueryUintDefault(r, "smaPeriod2", 14)
		smaPeriod3 = getQueryUintDefault(r, "smaPeriod3", 50)
	}

	ema := r.URL.Query().Get("ema")
	emaEnable := ema == "true"
	var emaPeriod1, emaPeriod2, emaPeriod3 int
	if emaEnable {
		emaPeriod1 = getQueryUintDefault(r, "emaPeriod1", 7)
		emaPeriod2 = getQueryUintDefault(r, "emaPeriod2", 14)
		emaPeriod3 = getQueryUintDefault(r, "emaPeriod3", 50)
	}

	bbands := r.URL.Query().Get("bbands")
	bbandsEnable := bbands == "true"
	var bbandsN int
	var bbandsK float64
	if bbandsEnable {
		bbandsN = getQueryUintDefault(r, "bbandsN", 20)
		bbandsK = getQueryFloatDefault(r, "bbandsK", 2)
	}

	ichimoku := r.URL.Query().Get("ichimoku")
	ichimokuEnable := ichimoku == "true"

	rsi := r.URL.Query().Get("rsi")
	rsiEnable := rsi == "true"
	var rsiPeriod int
	var rsiBuyThread, rsiSellThread float64
	if rsiEnable {
		rsiPeriod = getQueryUintDefault(r, "rsiPeriod", 14)
		rsiBuyThread = 30.0
		rsiSellThread = 70.0
	}

	macd := r.URL.Query().Get("macd")
	macdEnable := macd == "true"
	var macdFastPeriod, macdSlowPeriod, macdSignalPeriod int
	if macdEnable {
		macdFastPeriod = getQueryUintDefault(r, "macdPeriod1", 12)
		macdSlowPeriod = getQueryUintDefault(r, "macdPeriod2", 26)
		macdSignalPeriod = getQueryUintDefault(r, "macdPeriod3", 9)
	}

	stopLimitPercent := 0.75

	params := model.NewTradeParams(
		false,
		productCode,
		tradeSize,
		smaEnable,
		smaPeriod1,
		smaPeriod2,
		smaPeriod3,
		emaEnable,
		emaPeriod1,
		emaPeriod2,
		emaPeriod3,
		bbandsEnable,
		bbandsN,
		bbandsK,
		ichimokuEnable,
		rsiEnable,
		rsiPeriod,
		rsiBuyThread,
		rsiSellThread,
		macdEnable,
		macdFastPeriod,
		macdSlowPeriod,
		macdSignalPeriod,
		stopLimitPercent,
	)

	return params
}

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

func getQueryFloatDefault(r *http.Request, query string, def float64) float64 {
	strVal := r.URL.Query().Get(query)
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if strVal == "" || err != nil || floatVal < 0 {
		floatVal = def
	}
	return floatVal
}
