package model

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/bitflyer"
)

func TestCandle(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	// 2100/01/01 00:00:00.00 UTC
	timeDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	open, close, high, low, volume := 3000.0, 3200.0, 3500.0, 2500.0, 100000.0
	candle := NewCandle(config.ProductCode, config.CandleDuration, timeDate, open, close, high, low, volume)

	t.Run("create candle", func(t *testing.T) {
		err := candle.Create(tx, config.CandleTableName, config.TimeFormat)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("get candle", func(t *testing.T) {
		getCandle := GetCandle(tx, config.CandleTableName, config.TimeFormat, candle.ProductCode, candle.Duration, candle.Time)
		if *getCandle != *candle {
			t.Fatalf("%v != %v", *getCandle, *candle)
		}
	})

	t.Run("save candle", func(t *testing.T) {
		// +500.0
		candle.Open += 500.0
		candle.Close += 500.0
		candle.High += 500.0
		candle.Low += 500.0
		candle.Volume += 500.0

		err := candle.Save(tx, config.CandleTableName, config.TimeFormat)
		if err != nil {
			t.Fatal(err.Error())
		}

		// compare
		getCandle := GetCandle(tx, config.CandleTableName, config.TimeFormat, candle.ProductCode, candle.Duration, candle.Time)
		if *getCandle != *candle {
			t.Fatalf("%v != %v", *getCandle, *candle)
		}
	})

	t.Run("create candle with duration", func(t *testing.T) {
		ticker := bitflyer.Ticker{
			ProductCode:     config.ProductCode,
			State:           bitflyer.BoardStateRunning,
			Timestamp:       "2100-01-02T00:00:00.00",
			TickID:          39481886,
			BestBid:         238804,
			BestAsk:         238941,
			BestBidSize:     0.5,
			BestAskSize:     5.14,
			TotalBidDepth:   4426.816638,
			TotalAskDepth:   2641.5711685,
			MarketBidSize:   0,
			MarketAskSize:   0,
			Ltp:             238883,
			Volume:          24647.7435024,
			VolumeByProduct: 24647.7435024,
		}

		err := CreateCandleWithDuration(tx, config.CandleTableName, config.TimeFormat, config.LocalTime, config.TradeHour, &ticker, config.ProductCode, config.CandleDuration)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("get all candle", func(t *testing.T) {
		limit := 365

		candles, err := GetAllCandle(tx, config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration, limit)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(candles) > limit {
			t.Fatalf("get %d (> %d) candles", len(candles), limit)
		}
	})
}

func TestDateTimeUTC(t *testing.T) {
	cases := []struct {
		timeString   string
		expectedTime time.Time
	}{
		{
			timeString:   "2021-01-01T00:00:00.00",
			expectedTime: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			timeString:   "2021-01-02T09:00:00.00",
			expectedTime: time.Date(2021, time.January, 2, 9, 0, 0, 0, time.UTC),
		},
		{
			timeString:   "2021-01-03T12:00:00.00",
			expectedTime: time.Date(2021, time.January, 3, 12, 0, 0, 0, time.UTC),
		},
		{
			timeString:   "2021-01-04T23:59:59.99",
			expectedTime: time.Date(2021, time.January, 4, 23, 59, 59, 990000000, time.UTC),
		},
	}

	for _, c := range cases {
		dateTime := DateTimeUTC(c.timeString)
		if !dateTime.Equal(c.expectedTime) {
			t.Fatalf("%v != %v", dateTime, c.expectedTime)
		}
	}
}

func TestTruncateDateTime(t *testing.T) {
	cases := []struct {
		dateTime     time.Time
		hour         int
		expectedTime time.Time
	}{
		{
			dateTime:     time.Date(2021, time.January, 1, 2, 3, 4, 5, time.UTC),
			hour:         0,
			expectedTime: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			dateTime:     time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			hour:         1,
			expectedTime: time.Date(2020, time.December, 31, 1, 0, 0, 0, time.UTC),
		},
		{
			dateTime:     time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			dateTime:     time.Date(2021, time.January, 1, 10, 0, 0, 0, time.UTC),
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			dateTime:     time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
			hour:         15,
			expectedTime: time.Date(2020, time.December, 31, 15, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range cases {
		dateTime := TruncateDateTime(c.dateTime, c.hour)
		if !dateTime.Equal(c.expectedTime) {
			t.Fatalf("%v != %v", dateTime, c.expectedTime)
		}
	}
}
