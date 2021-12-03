package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestDataFrame(t *testing.T) {
	cr := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	candles, err := cr.FindAll(config.ProductCode, config.CandleDuration, -1)
	if err != nil {
		t.Fatal(err.Error())
	}

	signalEvents := model.NewSignalEvents(make([]model.SignalEvent, 0))
	df := model.NewDataFrame(config.ProductCode, candles, signalEvents)
	if df == nil {
		t.Fatal("NewDataFrame() returns nil")
	}

	t.Run("check length", func(t *testing.T) {
		lenCandles := len(df.Candles())
		if lenCandles != len(candles) {
			t.Fatal("lenCandles != len(candles)")
		}

		if len(df.Times()) != lenCandles {
			t.Fatal("len(df.Times()) != lenCandles")
		}

		if len(df.Opens()) != lenCandles {
			t.Fatal("len(df.Opens()) != lenCandles")
		}

		if len(df.Closes()) != lenCandles {
			t.Fatal("len(df.Closes()) != lenCandles")
		}

		if len(df.Highs()) != lenCandles {
			t.Fatal("len(df.Highs()) != lenCandles")
		}

		if len(df.Lows()) != lenCandles {
			t.Fatal("len(df.Lows()) != lenCandles")
		}

		if len(df.Volumes()) != lenCandles {
			t.Fatal("len(df.Volumes()) != lenCandles")
		}
	})

	t.Run("add indicator", func(t *testing.T) {
		df.AddSMA(7)
		df.AddSMA(14)
		df.AddSMA(50)

		df.AddEMA(7)
		df.AddEMA(14)
		df.AddEMA(50)

		df.AddBBands(20, 2)

		df.AddIchimoku()

		df.AddRSI(14)

		df.AddMACD(12, 26, 9)
	})

	t.Run("add signal_events", func(t *testing.T) {
		signalEvents := model.NewSignalEvents(make([]model.SignalEvent, 0))
		df.AddBacktestEvents(signalEvents)
	})
}
