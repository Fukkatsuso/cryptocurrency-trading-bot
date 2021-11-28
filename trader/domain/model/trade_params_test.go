package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
)

func TestTradeParams(t *testing.T) {
	params := model.NewTradeParams(
		true,
		config.ProductCode,
		0.1,
		true,
		7,
		14,
		50,
		true,
		7,
		14,
		50,
		true,
		20,
		2,
		true,
		true,
		14,
		30,
		70,
		true,
		12,
		26,
		9,
		0.75,
	)
	if params == nil {
		t.Fatal("NewTradeParams() returns nil")
	}

	t.Run("disable indicator", func(t *testing.T) {
		params.EnableSMA(false)
		if params.SMAEnable() {
			t.Fatal("EnableSMA(false) should disable sma")
		}

		params.EnableEMA(false)
		if params.EMAEnable() {
			t.Fatal("EnableEMA(false) should disable ema")
		}

		params.EnableBBands(false)
		if params.BBandsEnable() {
			t.Fatal("EnableBBands(false) should disable bbands")
		}

		params.EnableIchimoku(false)
		if params.IchimokuEnable() {
			t.Fatal("EnableIchimoku(false) should disable ichimoku_cloud")
		}

		params.EnableRSI(false)
		if params.RSIEnable() {
			t.Fatal("EnableRSI(false) should disable rsi")
		}

		params.EnableMACD(false)
		if params.MACDEnable() {
			t.Fatal("EnableMACD(false) should disable macd")
		}
	})
}
