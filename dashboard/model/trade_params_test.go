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

	var tradeParams *TradeParams

	t.Run("create trade_params", func(t *testing.T) {
		tradeParams = &TradeParams{
			TradeEnable: true,
			ProductCode: config.ProductCode,
			Size:        0.01,
		}

		err := tradeParams.Create(tx, config.TradeParamTableName)
		if err != nil {
			t.Fatal("Failed to Create TradeParams:", tradeParams, err.Error())
		}
	})

	t.Run("get trade_params", func(t *testing.T) {
		getTradeParams := GetTradeParams(tx, config.TradeParamTableName, config.ProductCode)

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

	emaPerformance, emaPeriod1, emaPeriod2 := df.OptimizeEMA(params.EMAPeriod1, params.EMAPeriod2, params.Size)
	t.Logf("ema: %f (period1=%d, period2=%d)", emaPerformance, emaPeriod1, emaPeriod2)

	bbandsPerformance, bbandsN, bbandsK := df.OptimizeBBands(params.BBandsN, params.BBandsK, params.Size)
	t.Logf("bbands: %f (n=%d, k=%f)", bbandsPerformance, bbandsN, bbandsK)

	ichimokuPerformance := df.OptimizeIchimoku(params.Size)
	t.Logf("ichimoku: %f", ichimokuPerformance)

	rsiPerformance, rsiPeriod, rsiBuyThread, rsiSellThread := df.OptimizeRSI(params.RSIPeriod, params.RSIBuyThread, params.RSISellThread, params.Size)
	t.Logf("rsi: %f (period=%d, buyThread=%f sellThread=%f)", rsiPerformance, rsiPeriod, rsiBuyThread, rsiSellThread)

	macdPerformance, macdFastPeriod, macdSlowPeriod, macdSignalPeriod := df.OptimizeMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod, params.Size)
	t.Logf("macd: %f (fastPeriod=%d, slowPeriod=%d, signalPeriod=%d)", macdPerformance, macdFastPeriod, macdSlowPeriod, macdSignalPeriod)
}

func deleteTradeParamsAll(tx DB) error {
	cmd := fmt.Sprintf("DELETE FROM %s", config.TradeParamTableName)
	_, err := tx.Exec(cmd)
	return err
}
