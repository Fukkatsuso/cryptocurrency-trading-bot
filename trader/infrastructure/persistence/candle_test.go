package persistence_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
)

func newCandles() []model.Candle {
	table := []struct {
		time   time.Time
		open   float64
		close  float64
		high   float64
		low    float64
		volume float64
	}{
		{
			time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			3000.0,
			3200.0,
			3500.0,
			2500.0,
			100000.0,
		},
	}

	candles := make([]model.Candle, 0)
	for _, t := range table {
		candleTime := model.NewCandleTime(t.time)
		candle := model.NewCandle(config.ProductCode, config.CandleDuration, candleTime, t.open, t.close, t.high, t.low, t.volume)
		if candle == nil {
			continue
		}
		candles = append(candles, *candle)
	}
	return candles
}

func TestCandle(t *testing.T) {
	tx := NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleRepository(tx, config.CandleTableName, config.TimeFormat)

	candles := newCandles()

	t.Run("save candle", func(t *testing.T) {
		for _, candle := range candles {
			err := candleRepository.Save(candle)
			if err != nil {
				t.Fatal(err.Error())
			}
		}
	})

	t.Run("find all candle", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			cc, err := candleRepository.FindAll(config.ProductCode, config.CandleDuration, int64(i))
			if err != nil {
				t.Fatal(err.Error())
			}
			if len(cc) > i {
				t.Fatal("limit of FindAll() is invalid")
			}
		}
	})

	t.Run("find candle", func(t *testing.T) {
		c1 := candles[0]
		c2, err := candleRepository.FindByCandleTime(c1.ProductCode(), c1.Duration(), c1.Time())
		if err != nil {
			t.Fatal(err.Error())
		}
		if c1 != *c2 {
			t.Fatalf("%+v != %+v", c1, *c2)
		}
	})

	t.Run("update candle", func(t *testing.T) {
		c1 := candles[0]
		productCode := c1.ProductCode()
		duration := c1.Duration()
		candleTime := c1.Time()
		c2 := model.NewCandle(productCode, duration, candleTime, c1.Open()+500.0, c1.Close()+500.0, c1.High()+500.0, c1.Low()+500.0, c1.Volume()+500.0)

		err := candleRepository.Save(*c2)
		if err != nil {
			t.Fatal(err.Error())
		}

		c3, _ := candleRepository.FindByCandleTime(productCode, duration, candleTime)
		if *c2 != *c3 {
			t.Fatalf("%+v != %+v", *c2, *c3)
		}
	})
}
