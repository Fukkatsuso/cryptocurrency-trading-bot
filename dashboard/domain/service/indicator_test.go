package service_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestIndicatorService(t *testing.T) {
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	candles, err := candleRepository.FindAll(config.ProductCode, config.CandleDuration, -1)
	if err != nil {
		t.Fatal(err.Error())
	}

	df := model.NewDataFrame(config.ProductCode, candles, nil)
	inReal := df.Closes()
	lenCandle := len(candles)

	indicatorService := service.NewIndicatorService()

	t.Run("EMA", func(t *testing.T) {
		emaFast := model.NewEMA(inReal, 7)
		emaSlow := model.NewEMA(inReal, 14)

		buy := indicatorService.BuySignalOfEMA(emaFast, emaSlow, lenCandle-1)
		t.Logf("BuySignalOfEMA: %t", buy)

		sell := indicatorService.SellSignalOfEMA(emaFast, emaSlow, lenCandle-1)
		t.Logf("SellSignalOfEMA: %t", sell)
	})

	t.Run("BBands", func(t *testing.T) {
		bbands := model.NewBBands(inReal, 20, 2)

		buy := indicatorService.BuySignalOfBBands(bbands, candles, lenCandle-1)
		t.Logf("BuySignalOfBBands: %t", buy)

		sell := indicatorService.SellSignalOfBBands(bbands, candles, lenCandle-1)
		t.Logf("SellSignalOfBBands: %t", sell)
	})

	t.Run("Ichimoku Cloud", func(t *testing.T) {
		ichimoku := model.NewIchimokuCloud(inReal)

		buy := indicatorService.BuySignalOfIchimoku(ichimoku, candles, lenCandle-1)
		t.Logf("BuySignalOfIchimoku: %t", buy)

		sell := indicatorService.SellSignalOfIchimoku(ichimoku, candles, lenCandle-1)
		t.Logf("SellSignalOfIchimoku: %t", sell)
	})

	t.Run("RSI", func(t *testing.T) {
		rsi := model.NewRSI(inReal, 14)

		buy := indicatorService.BuySignalOfRSI(rsi, 30, lenCandle-1)
		t.Logf("BuySignalOfRSI: %t", buy)

		sell := indicatorService.SellSignalOfRSI(rsi, 70, lenCandle-1)
		t.Logf("SellSignalOfRSI: %t", sell)
	})

	t.Run("MACD", func(t *testing.T) {
		macd := model.NewMACD(inReal, 12, 26, 9)

		buy := indicatorService.BuySignalOfMACD(macd, lenCandle-1)
		t.Logf("BuySignalOfMACD: %t", buy)

		sell := indicatorService.SellSignalOfMACD(macd, lenCandle-1)
		t.Logf("SellSignalOfMACD: %t", sell)
	})

	t.Run("Is Boxed Range?", func(t *testing.T) {
		term := 50
		rsiPeriod := 14

		// 疑似レンジ相場
		boxed := make([]float64, term)
		for i := range boxed {
			boxed[i] = float64(10001 - 2*(i&2))
		}
		rsi := model.NewRSI(boxed, rsiPeriod)
		t.Log(rsi.Values())
		if !indicatorService.IsBoxedRange(rsi, 7, term-1) {
			t.Fatal("not boxed range")
		}

		// 疑似トレンド相場
		trend := make([]float64, term)
		for i := range trend {
			trend[i] = float64(10000 + i)
		}
		rsi = model.NewRSI(trend, rsiPeriod)
		t.Log(rsi.Values())
		if indicatorService.IsBoxedRange(rsi, 7, term-1) {
			t.Fatal("boxed range")
		}
	})
}
