package model

import (
	"fmt"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
)

func TestTradeParams(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	err := deleteTradeParamsAll(tx)
	if err != nil {
		t.Fatal("failed to exec deleteTradeParamsAll")
	}

	t.Run("get nil trade_params", func(t *testing.T) {
		getTradeParams := GetTradeParams(tx, config.ProductCode)

		if getTradeParams != nil {
			t.Fatalf("Failed to Get TradeParams: should nil, but %v", getTradeParams)
		}
	})

	var tradeParams *TradeParams

	t.Run("create trade_params", func(t *testing.T) {
		tradeParams = &TradeParams{
			TradeEnable:      true,
			ProductCode:      config.ProductCode,
			Size:             0.01,
			SMAEnable:        true,
			SMAPeriod1:       7,
			SMAPeriod2:       14,
			SMAPeriod3:       50,
			EMAEnable:        true,
			EMAPeriod1:       7,
			EMAPeriod2:       14,
			EMAPeriod3:       50,
			BBandsEnable:     true,
			BBandsN:          20,
			BBandsK:          2.2,
			IchimokuEnable:   true,
			RSIEnable:        true,
			RSIPeriod:        14,
			RSIBuyThread:     30.5,
			RSISellThread:    70.5,
			MACDEnable:       true,
			MACDFastPeriod:   12,
			MACDSlowPeriod:   26,
			MACDSignalPeriod: 9,
			StopLimitPercent: 0.75,
		}

		err := tradeParams.Create(tx)
		if err != nil {
			t.Fatal("Failed to Create TradeParams:", tradeParams, err.Error())
		}
	})

	t.Run("get trade_params", func(t *testing.T) {
		getTradeParams := GetTradeParams(tx, config.ProductCode)

		if getTradeParams == nil {
			t.Fatal("Failed to Get TradeParams")
		}
		if *getTradeParams != *tradeParams {
			t.Fatalf("%v != %v", *getTradeParams, *tradeParams)
		}
	})
}

func TestBackTest(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	candles, err := CandleMockData()
	if err != nil {
		t.Fatal(err.Error())
	}
	df := DataFrame{
		ProductCode: config.ProductCode,
		Candles:     candles,
	}

	t.Run("backtest by default params", func(t *testing.T) {
		params := &TradeParams{
			TradeEnable:      true,
			ProductCode:      config.ProductCode,
			Size:             0.01,
			SMAEnable:        true,
			SMAPeriod1:       7,
			SMAPeriod2:       14,
			SMAPeriod3:       50,
			EMAEnable:        true,
			EMAPeriod1:       7,
			EMAPeriod2:       14,
			EMAPeriod3:       50,
			BBandsEnable:     true,
			BBandsN:          20,
			BBandsK:          2,
			IchimokuEnable:   true,
			RSIEnable:        true,
			RSIPeriod:        14,
			RSIBuyThread:     30,
			RSISellThread:    70,
			MACDEnable:       true,
			MACDFastPeriod:   12,
			MACDSlowPeriod:   26,
			MACDSignalPeriod: 9,
			StopLimitPercent: 0.75,
		}

		df.AddSMA(params.SMAPeriod1)
		df.AddSMA(params.SMAPeriod2)
		df.AddSMA(params.SMAPeriod3)
		df.AddEMA(params.EMAPeriod1)
		df.AddEMA(params.EMAPeriod2)
		df.AddEMA(params.EMAPeriod3)
		df.AddBBands(params.BBandsN, params.BBandsK)
		df.AddIchimoku()
		df.AddRSI(params.RSIPeriod)
		df.AddMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod)

		df.BackTest(params)

		t.Log("events:", df.BacktestEvents.Signals)
		t.Log("profit:", df.BacktestEvents.Profit)
	})
}

func TestOptimizeTradeParams(t *testing.T) {
	candles, err := CandleMockData()
	if err != nil {
		t.Fatal(err.Error())
	}
	df := DataFrame{
		ProductCode: config.ProductCode,
		Candles:     candles,
	}

	params := &TradeParams{
		TradeEnable:      true,
		ProductCode:      config.ProductCode,
		Size:             0.01,
		SMAEnable:        true,
		SMAPeriod1:       7,
		SMAPeriod2:       14,
		SMAPeriod3:       50,
		EMAEnable:        true,
		EMAPeriod1:       7,
		EMAPeriod2:       14,
		EMAPeriod3:       50,
		BBandsEnable:     true,
		BBandsN:          20,
		BBandsK:          2,
		IchimokuEnable:   true,
		RSIEnable:        true,
		RSIPeriod:        14,
		RSIBuyThread:     30,
		RSISellThread:    70,
		MACDEnable:       true,
		MACDFastPeriod:   12,
		MACDSlowPeriod:   26,
		MACDSignalPeriod: 9,
		StopLimitPercent: 0.75,
	}

	t.Run("optimize ema", func(t *testing.T) {
		emaPerformance, emaPeriod1, emaPeriod2 := df.OptimizeEMA(params.EMAPeriod1, params.EMAPeriod2, params.Size)
		t.Logf("ema profit: %f (period1=%d, period2=%d)", emaPerformance, emaPeriod1, emaPeriod2)
	})

	t.Run("optimize bbands", func(t *testing.T) {
		bbandsPerformance, bbandsN, bbandsK := df.OptimizeBBands(params.BBandsN, params.BBandsK, params.Size)
		t.Logf("bbands profit: %f (n=%d, k=%f)", bbandsPerformance, bbandsN, bbandsK)
	})

	t.Run("optimize ichimoku (only backtest)", func(t *testing.T) {
		ichimokuPerformance := df.OptimizeIchimoku(params.Size)
		t.Logf("ichimoku profit: %f", ichimokuPerformance)
	})

	t.Run("optimize rsi", func(t *testing.T) {
		rsiPerformance, rsiPeriod, rsiBuyThread, rsiSellThread := df.OptimizeRSI(params.RSIPeriod, params.RSIBuyThread, params.RSISellThread, params.Size)
		t.Logf("rsi profit: %f (period=%d, buyThread=%f sellThread=%f)", rsiPerformance, rsiPeriod, rsiBuyThread, rsiSellThread)
	})

	t.Run("optimize macd", func(t *testing.T) {
		macdPerformance, macdFastPeriod, macdSlowPeriod, macdSignalPeriod := df.OptimizeMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod, params.Size)
		t.Logf("macd profit: %f (fastPeriod=%d, slowPeriod=%d, signalPeriod=%d)", macdPerformance, macdFastPeriod, macdSlowPeriod, macdSignalPeriod)
	})

	t.Run("optimize all", func(t *testing.T) {
		optimizedParams := df.OptimizeTradeParams(params)

		if optimizedParams.SMAEnable {
			df.AddSMA(optimizedParams.SMAPeriod1)
			df.AddSMA(optimizedParams.SMAPeriod2)
			df.AddSMA(optimizedParams.SMAPeriod3)
		}
		if optimizedParams.EMAEnable {
			df.AddEMA(optimizedParams.EMAPeriod1)
			df.AddEMA(optimizedParams.EMAPeriod2)
			df.AddEMA(optimizedParams.EMAPeriod3)
		}
		if optimizedParams.BBandsEnable {
			df.AddBBands(optimizedParams.BBandsN, optimizedParams.BBandsK)
		}
		if optimizedParams.IchimokuEnable {
			df.AddIchimoku()
		}
		if optimizedParams.RSIEnable {
			df.AddRSI(optimizedParams.RSIPeriod)
		}
		if optimizedParams.MACDEnable {
			df.AddMACD(optimizedParams.MACDFastPeriod, optimizedParams.MACDSlowPeriod, optimizedParams.MACDSignalPeriod)
		}

		df.BackTest(optimizedParams)
		profit := df.BacktestEvents.Profit
		t.Logf("profit: %f", profit)
	})
}

func deleteTradeParamsAll(tx DB) error {
	cmd := fmt.Sprintf("DELETE FROM trade_params")
	_, err := tx.Exec(cmd)
	return err
}
