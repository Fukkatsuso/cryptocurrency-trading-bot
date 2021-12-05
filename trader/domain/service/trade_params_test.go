package service_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
)

func TestTradeParamsService(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	candles, err := candleRepository.FindAll(config.ProductCode, config.CandleDuration, -1)
	if err != nil {
		t.Fatal(err.Error())
	}
	df := model.NewDataFrame(config.ProductCode, candles, nil)

	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)
	tradeParamsService := service.NewTradeParamsService(tradeParamsRepository, dataFrameService)

	params := model.NewBasicTradeParams(config.ProductCode, 0.01)

	t.Run("save trade_params", func(t *testing.T) {
		err := tradeParamsService.Save(*params)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("find trade_params", func(t *testing.T) {
		findParams, err := tradeParamsService.Find(config.ProductCode)
		if err != nil {
			t.Fatal(err.Error())
		}
		if findParams == nil ||
			*findParams != *params {
			t.Fatalf("%+v != %+v", *findParams, *params)
		}
	})

	t.Run("optimize EMA", func(t *testing.T) {
		performance, fastPeriod, slowPeriod, changed := tradeParamsService.OptimizeEMA(df, params.EMAPeriod1(), params.EMAPeriod2(), params.Size())
		t.Logf("performance=%f, fastPeriod=%d, slowPeriod=%d", performance, fastPeriod, slowPeriod)
		if changed &&
			(fastPeriod == params.EMAPeriod1() && slowPeriod == params.EMAPeriod2()) {
			t.Fatal("params is not changed")
		} else if !changed &&
			(fastPeriod != params.EMAPeriod1() || slowPeriod != params.EMAPeriod2()) {
			t.Fatal("params is changed")
		}
	})

	t.Run("optimize bbands", func(t *testing.T) {
		performance, n, k, changed := tradeParamsService.OptimizeBBands(df, params.BBandsN(), params.BBandsK(), params.Size())
		t.Logf("performance=%f, n=%d, k=%f", performance, n, k)
		if changed &&
			(n == params.BBandsN() && k == params.BBandsK()) {
			t.Fatal("params is not changed")
		} else if !changed &&
			(n != params.BBandsN() || k != params.BBandsK()) {
			t.Fatal("params is changed")
		}
	})

	t.Run("optimize ichimoku cloud", func(t *testing.T) {
		performance, changed := tradeParamsService.OptimizeIchimoku(df, params.Size())
		t.Logf("performance=%f", performance)
		if changed {
			t.Fatal("params is changed(?)")
		}
	})

	t.Run("optimize rsi", func(t *testing.T) {
		performance, period, buyThread, sellThread, changed := tradeParamsService.OptimizeRSI(df, params.RSIPeriod(), params.RSIBuyThread(), params.RSISellThread(), params.Size())
		t.Logf("performance=%f, period=%d, buyThread=%f, sellThread=%f", performance, period, buyThread, sellThread)
		if changed &&
			(period == params.RSIPeriod() && buyThread == params.RSIBuyThread() && sellThread == params.RSISellThread()) {
			t.Fatal("params is not changed")
		} else if !changed &&
			(period != params.RSIPeriod() || buyThread != params.RSIBuyThread() || sellThread != params.RSISellThread()) {
			t.Fatal("params is changed")
		}
	})

	t.Run("optimize macd", func(t *testing.T) {
		performance, fastPeriod, slowPeriod, signalPeriod, changed := tradeParamsService.OptimizeMACD(df, params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSignalPeriod(), params.Size())
		t.Logf("performance=%f, fastPeriod=%d, slowPeriod=%d, signalPeriod=%d", performance, fastPeriod, slowPeriod, signalPeriod)
		if changed &&
			(fastPeriod == params.MACDFastPeriod() && slowPeriod == params.MACDSlowPeriod() && signalPeriod == params.MACDSignalPeriod()) {
			t.Fatal("params is not changed")
		} else if !changed &&
			(fastPeriod != params.MACDFastPeriod() || slowPeriod != params.MACDSlowPeriod() || signalPeriod != params.MACDSignalPeriod()) {
			t.Fatal("params is changed")
		}
	})

	t.Run("optimize all", func(t *testing.T) {
		optimizedParams, changed := tradeParamsService.OptimizeAll(df, params)
		t.Logf("params=%+v", optimizedParams)
		if changed && *optimizedParams == *params {
			t.Fatal("params is not changed")
		} else if !changed && *optimizedParams != *params {
			t.Fatal("params is changed")
		}
	})
}
