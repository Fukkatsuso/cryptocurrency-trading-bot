package model_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
)

func TestNewCandle(t *testing.T) {
	table := []struct {
		productCode string
		duration    time.Duration
		time        model.CandleTime
		open        float64
		close       float64
		high        float64
		low         float64
		volume      float64
		ok          bool
	}{
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          true,
		},
		{
			productCode: "",
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    -1,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        -3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       -3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        -3500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         -2500.0,
			volume:      100000.0,
			ok:          false,
		},
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         2500.0,
			volume:      -100000.0,
			ok:          false,
		},
		// high < low
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        1500.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		// high < open
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        2800.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		// high < close
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3100.0,
			low:         2500.0,
			volume:      100000.0,
			ok:          false,
		},
		// low > open
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         3100.0,
			volume:      100000.0,
			ok:          false,
		},
		// low > close
		{
			productCode: config.ProductCode,
			duration:    config.CandleDuration,
			time:        model.NewCandleTime(time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)),
			open:        3000.0,
			close:       3200.0,
			high:        3500.0,
			low:         3300.0,
			volume:      100000.0,
			ok:          false,
		},
	}

	for _, c := range table {
		candle := model.NewCandle(c.productCode, c.duration, c.time, c.open, c.close, c.high, c.low, c.volume)
		if c.ok && candle == nil {
			t.Fatal("NewCandle() returns nil")
		} else if !c.ok && candle != nil {
			t.Fatal("NewCandle() returns not nil")
		}
	}
}
