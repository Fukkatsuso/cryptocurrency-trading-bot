package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestNewTicker(t *testing.T) {
	table := []struct {
		productCode     string
		state           string
		timestamp       string
		tickID          int
		bestBid         float64
		bestAsk         float64
		bestBidSize     float64
		bestAskSize     float64
		totalBidDepth   float64
		totalAskDepth   float64
		marketBidSize   float64
		marketAskSize   float64
		ltp             float64
		volume          float64
		volumeByProduct float64
		ok              bool
	}{
		{
			productCode:     config.ProductCode,
			state:           "RUNNING",
			timestamp:       "2006-01-02T15:04:05",
			tickID:          1,
			bestBid:         10,
			bestAsk:         10,
			bestBidSize:     20,
			bestAskSize:     20,
			totalBidDepth:   30,
			totalAskDepth:   30,
			marketBidSize:   40,
			marketAskSize:   40,
			ltp:             50,
			volume:          60,
			volumeByProduct: 10,
			ok:              true,
		},
		{
			productCode:     "",
			state:           "",
			timestamp:       "",
			tickID:          -1,
			bestBid:         -1,
			bestAsk:         -1,
			bestBidSize:     -1,
			bestAskSize:     -1,
			totalBidDepth:   -1,
			totalAskDepth:   -1,
			marketBidSize:   -1,
			marketAskSize:   -1,
			ltp:             -1,
			volume:          -1,
			volumeByProduct: -1,
			ok:              false,
		},
	}

	for _, tc := range table {
		ticker := model.NewTicker(
			tc.productCode,
			tc.state,
			tc.timestamp,
			tc.tickID,
			tc.bestBid,
			tc.bestAsk,
			tc.bestBidSize,
			tc.bestAskSize,
			tc.totalBidDepth,
			tc.totalAskDepth,
			tc.marketBidSize,
			tc.marketAskSize,
			tc.ltp,
			tc.volume,
			tc.volumeByProduct,
		)

		if ticker == nil && tc.ok {
			t.Fatal("NewTicker() returns nil")
		} else if ticker != nil && !tc.ok {
			t.Fatal("NewTicker() returns not nil")
		}
	}
}
