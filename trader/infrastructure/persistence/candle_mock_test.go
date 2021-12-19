package persistence_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
)

func TestCandleMock(t *testing.T) {
	cr := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)

	t.Run("find all candle", func(t *testing.T) {
		candles, err := cr.FindAll(config.ProductCode, config.CandleDuration, 10)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(candles) > 10 {
			t.Fatal("len(candles) > 10")
		}

		_, err = cr.FindAll(config.ProductCode, config.CandleDuration, -1)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
