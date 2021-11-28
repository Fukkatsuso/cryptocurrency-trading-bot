package model_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestNewSignalEvent(t *testing.T) {
	table := []struct {
		time        time.Time
		productCode string
		side        model.OrderSide
		price       float64
		size        float64
		ok          bool
	}{
		{
			time:        time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideBuy,
			price:       1000,
			size:        0.01,
			ok:          true,
		},
		{
			time:        time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: "",
			side:        model.OrderSideBuy,
			price:       1000,
			size:        0.01,
			ok:          false,
		},
		{
			time:        time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        "",
			price:       1000,
			size:        0.01,
			ok:          false,
		},
		{
			time:        time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideBuy,
			price:       -1000,
			size:        0.01,
			ok:          false,
		},
		{
			time:        time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideBuy,
			price:       1000,
			size:        -0.01,
			ok:          false,
		},
	}

	for _, s := range table {
		signalEvent := model.NewSignalEvent(s.time, s.productCode, s.side, s.price, s.size)
		if s.ok && signalEvent == nil {
			t.Fatal("NewSignalEvent() returns nil")
		} else if !s.ok && signalEvent != nil {
			t.Fatal("NewSignalEvent() returns not nil")
		}
	}
}

func TestSignalEvents(t *testing.T) {
	table := []struct {
		time        time.Time
		productCode string
		side        model.OrderSide
		price       float64
		size        float64
	}{
		{
			time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideBuy,
			price:       1000,
			size:        1,
		},
		{
			time:        time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideSell,
			price:       2000,
			size:        1,
		},
		{
			time:        time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideBuy,
			price:       10000,
			size:        1,
		},
		{
			time:        time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
			productCode: config.ProductCode,
			side:        model.OrderSideSell,
			price:       20000,
			size:        1,
		},
	}

	signalEventList := make([]model.SignalEvent, 0)
	for _, s := range table {
		signalEvent := model.NewSignalEvent(s.time, s.productCode, s.side, s.price, s.size)
		signalEventList = append(signalEventList, *signalEvent)
	}

	var signalEvents *model.SignalEvents
	t.Run("NewSignalEvents", func(t *testing.T) {
		signalEvents = model.NewSignalEvents(signalEventList)
		if signalEvents == nil {
			t.Fatal("NewSignalEvents() returns nil")
		}
	})

	t.Run("LastSignal", func(t *testing.T) {
		lastEvent := signalEventList[len(signalEventList)-1]
		if *signalEvents.LastSignal() != lastEvent {
			t.Fatalf("%+v != %+v", *signalEvents.LastSignal(), lastEvent)
		}
	})

	buyTime := time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)

	t.Run("CanBuyAt", func(t *testing.T) {
		if !signalEvents.CanBuyAt(buyTime) {
			t.Fatal("CanBuyAt() returns false")
		}
	})

	t.Run("AddBuySignal", func(t *testing.T) {
		signalEvent := model.NewSignalEvent(buyTime, config.ProductCode, model.OrderSideBuy, 100000, 1)
		if !signalEvents.AddBuySignal(*signalEvent) {
			t.Fatal("AddBuySignal() returns false")
		}
	})

	sellTime := time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)

	t.Run("CanSellAt", func(t *testing.T) {
		if !signalEvents.CanSellAt(sellTime) {
			t.Fatal("CanSellAt() returns false")
		}
	})

	t.Run("AddSellSignal", func(t *testing.T) {
		signalEvent := model.NewSignalEvent(sellTime, config.ProductCode, model.OrderSideSell, 200000, 1)
		if !signalEvents.AddSellSignal(*signalEvent) {
			t.Fatal("AddSellSignal() returns false")
		}
	})

	t.Run("EstimateProfit", func(t *testing.T) {
		profit := signalEvents.EstimateProfit()
		if profit != signalEvents.Profit() {
			t.Fatalf("%v != %v", profit, signalEvents.Profit())
		}
		needProfit := float64(-1000*1 + 2000*1 - 10000*1 + 20000*1 - 100000*1 + 200000*1)
		if profit != needProfit {
			t.Fatalf("%v != %v", profit, needProfit)
		}
	})

	t.Run("ShouldCutLoss", func(t *testing.T) {
		time := time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)
		buyPrice := float64(1000000)
		signalEvent := model.NewSignalEvent(time, config.ProductCode, model.OrderSideBuy, buyPrice, 1)
		signalEvents.AddBuySignal(*signalEvent)

		currentPrice := buyPrice * 0.5
		stopLimitPercent := 0.6
		if !signalEvents.ShouldCutLoss(currentPrice, stopLimitPercent) {
			t.Fatal("ShouldCutLoss() returns false")
		}
	})
}
