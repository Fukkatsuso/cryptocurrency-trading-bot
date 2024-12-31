package service_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestCandleServicePerDay(t *testing.T) {
	tx := persistence.NewSQLiteTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleRepository(tx, config.CandleTableName, config.TimeFormat)
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)

	var candle *model.Candle

	t.Run("ticker to candle", func(t *testing.T) {
		ticker := model.NewTicker(
			config.ProductCode,
			"RUNNING",
			"2100-01-01T08:28:46.02",
			1000000000,
			519000,
			519281,
			0.0104,
			0.2,
			4311.5029572,
			1749.5828469,
			0,
			0,
			519298,
			14766.24084,
			14766.24084,
		)

		candle = candleService.TickerToCandle(*ticker)
		if candle == nil {
			t.Fatal("TickerToCandle() returns nil")
		}

		// "分"以下が切り捨てられているか
		time := candle.Time().Time()
		if time.Minute() != 0 ||
			time.Second() != 0 ||
			time.Nanosecond() != 0 {
			t.Fatal("candle.Time() is not truncated")
		}
	})

	t.Run("update candle", func(t *testing.T) {
		var newCandle *model.Candle
		high, low := candle.High(), candle.Low()

		high += 1000
		low -= 1000
		newCandle = model.NewCandle(
			candle.ProductCode(),
			candle.Duration(),
			candle.Time(),
			candle.Open(),
			candle.Close(),
			high,
			low,
			candle.Volume(),
		)
		candle = candleService.Update(candle, newCandle)
		if candle == nil {
			t.Fatal("Update() returns nil")
		}
		if candle.High() != high {
			t.Fatalf("candle.High() != %f", high)
		}
		if candle.Low() != low {
			t.Fatalf("candle.Low() != %f", low)
		}

		high -= 1000
		low += 1000
		newCandle = model.NewCandle(
			candle.ProductCode(),
			candle.Duration(),
			candle.Time(),
			candle.Open(),
			candle.Close(),
			high,
			low,
			candle.Volume(),
		)
		candle = candleService.Update(candle, newCandle)
		if candle == nil {
			t.Fatal("Update() returns nil")
		}
		if candle.High() <= high {
			t.Fatalf("candle.High() <= %f", high)
		}
		if candle.Low() >= low {
			t.Fatalf("candle.Low() >= %f", low)
		}
	})

	t.Run("save candle", func(t *testing.T) {
		err := candleService.Save(*candle)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("find by time", func(t *testing.T) {
		time := candle.Time().Time()
		_, err := candleService.FindByTime(config.ProductCode, time)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("find all candle", func(t *testing.T) {
		candles, err := candleService.FindAll(config.ProductCode, 10)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(candles) > 10 {
			t.Fatal("len(candles) > 10")
		}
	})
}
