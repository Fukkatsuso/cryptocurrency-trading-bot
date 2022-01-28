package service_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestDataFrameService(t *testing.T) {
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	candles, err := candleRepository.FindAll(config.ProductCode, config.CandleDuration, -1)
	if err != nil {
		t.Fatal(err.Error())
	}
	df := model.NewDataFrame(config.ProductCode, candles, nil)

	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)

	t.Run("EMA", func(t *testing.T) {
		events := dataFrameService.BacktestEMA(df, 7, 14, 0.01)
		t.Logf("BacktestEMA: %v", events)
	})

	t.Run("BBands", func(t *testing.T) {
		events := dataFrameService.BacktestBBands(df, 20, 2, 0.01)
		t.Logf("BacktestBBands: %v", events)
	})

	t.Run("Ichimoku Cloud", func(t *testing.T) {
		events := dataFrameService.BacktestIchimoku(df, 0.01)
		t.Logf("BacktestIchimoku: %v", events)
	})

	t.Run("RSI", func(t *testing.T) {
		events := dataFrameService.BacktestRSI(df, 14, 30, 70, 0.01)
		t.Logf("BacktestRSI: %v", events)
	})

	t.Run("MACD", func(t *testing.T) {
		events := dataFrameService.BacktestMACD(df, 12, 26, 9, 0.01)
		t.Logf("BacktestMACD: %v", events)
	})

	params := model.NewBasicTradeParams(config.ProductCode, 0.01)
	// addXXX()するタイミングは再考の余地あり
	df.AddEMA(params.EMAPeriod1())
	df.AddEMA(params.EMAPeriod2())
	df.AddBBands(params.BBandsN(), params.BBandsK())
	df.AddIchimoku()
	df.AddRSI(params.RSIPeriod())
	df.AddMACD(params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSlowPeriod())

	t.Run("Analyze", func(t *testing.T) {
		buy, sell := dataFrameService.Analyze(df, len(candles)-1, params)
		t.Logf("Analyze: buy=%t, sell=%t", buy, sell)
	})

	t.Run("Backtest", func(t *testing.T) {
		dataFrameService.Backtest(df, params)
		events := df.BacktestEvents()
		if events == nil {
			t.Fatal("Backtest() does not set SignalEvents")
		}
		t.Logf("Signals: %v", events.Signals())
		t.Logf("Profit: %f", events.Profit())
	})
}

func TestMRBaseDataFrameService(t *testing.T) {
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	candles, err := candleRepository.FindAll(config.ProductCode, config.CandleDuration, -1)
	if err != nil {
		t.Fatal(err.Error())
	}
	df := model.NewDataFrame(config.ProductCode, candles, nil)

	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewMRBaseDataFrameService(indicatorService)

	params := model.NewBasicTradeParams(config.ProductCode, 0.01)
	df.AddRSI(params.RSIPeriod())
	df.AddMACD(params.MACDFastPeriod(), params.MACDSlowPeriod(), params.MACDSlowPeriod())

	t.Run("Analyze", func(t *testing.T) {
		buy, sell := dataFrameService.Analyze(df, len(candles)-1, params)
		t.Logf("Analyze: buy=%t, sell=%t", buy, sell)
	})

	t.Run("Backtest", func(t *testing.T) {
		dataFrameService.Backtest(df, params)
		events := df.BacktestEvents()
		if events == nil {
			t.Fatal("Backtest() does not set SignalEvents")
		}
		t.Logf("Signals: %v", events.Signals())
		t.Logf("Profit: %f", events.Profit())
	})
}
