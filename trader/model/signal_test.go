package model

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
)

func TestSignalEventSave(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	// 2100/01/01 00:00:00.00 UTC
	timeDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	event := &SignalEvent{
		Time:        timeDate,
		ProductCode: config.ProductCode,
		Side:        string(bitflyer.OrderSideBuy),
		Price:       1000,
		Size:        0.01,
	}

	ok := event.Save(tx, config.TimeFormat)
	if !ok {
		t.Fatal("Failed to Save SignalEvent:", event)
	}
}

func TestSignalEvents(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	// テスト前の準備として全削除する
	err := deleteSignalEventAll(tx)
	if err != nil {
		t.Fatal("failed to exec deleteSignalEventAll")
	}

	var events *SignalEvents

	t.Run("get zero signal_event", func(t *testing.T) {
		events = GetSignalEvents(tx, config.ProductCode)

		if events == nil {
			t.Fatal("Failed to GetSignalEvents")
		}

		if len(events.Signals) != 0 {
			t.Fatalf("wrong number of SignalEvents. Expected 0, but %d", len(events.Signals))
		}
	})

	t.Run("save some signal_event", func(t *testing.T) {
		// 2021/01/01 00:00:00.00 UTC
		timeDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		event := &SignalEvent{
			Time:        timeDate,
			ProductCode: config.ProductCode,
			Side:        string(bitflyer.OrderSideBuy),
			Price:       1000,
			Size:        0.01,
		}
		ok := event.Save(tx, config.TimeFormat)
		if !ok {
			t.Fatal("Failed to Save SignalEvent:", event)
		}

		// 2021/01/02 00:00:00.00 UTC
		timeDate = time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
		event = &SignalEvent{
			Time:        timeDate,
			ProductCode: config.ProductCode,
			Side:        string(bitflyer.OrderSideSell),
			Price:       1500,
			Size:        0.01,
		}
		ok = event.Save(tx, config.TimeFormat)
		if !ok {
			t.Fatal("Failed to Save SignalEvent:", event)
		}
	})

	t.Run("get signal_events", func(t *testing.T) {
		events = GetSignalEvents(tx, config.ProductCode)

		if events == nil {
			t.Fatal("Failed to GetSignalEvents")
		}

		if len(events.Signals) != 2 {
			t.Fatalf("wrong number of SignalEvents. Expected 2, but %d", len(events.Signals))
		}
	})

	t.Run("get signal_events after time", func(t *testing.T) {
		timeDate := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)
		events := GetSignalEventsAfterTime(tx, config.ProductCode, timeDate, config.TimeFormat)

		if events == nil {
			t.Fatal("Failed to GetSignalEventsAfterTime")
		}

		if len(events.Signals) != 1 {
			t.Fatalf("wrong number of SignalEvents. Expected 1, but %d", len(events.Signals))
		}
	})

	t.Run("test CanBuy, Buy", func(t *testing.T) {
		// can buy
		now := time.Now().UTC()
		if !events.CanBuy(now) {
			t.Fatal("events should enable to buy:", events)
		}
		ok := events.Buy(config.ProductCode, now, 1200, 0.01)
		if !ok || len(events.Signals) != 3 {
			t.Fatal("Failed to add buy event")
		}

		// cannot buy
		now = time.Now().UTC()
		if events.CanBuy(now) {
			t.Fatal("events should disable to buy:", events)
		}
		ok = events.Buy(config.ProductCode, now, 1200, 0.01)
		if ok || len(events.Signals) != 3 {
			t.Fatal("Failed to disable to add buy event")
		}
	})

	t.Run("test CanSell, Sell", func(t *testing.T) {
		// can sell
		now := time.Now().UTC()
		if !events.CanSell(now) {
			t.Fatal("events should enable to sell:", events)
		}
		ok := events.Sell(config.ProductCode, now, 2000, 0.01)
		if !ok || len(events.Signals) != 4 {
			t.Fatal("Failed to add sell event")
		}

		// cannot sell
		now = time.Now().UTC()
		if events.CanSell(now) {
			t.Fatal("events should disable to sell:", events)
		}
		ok = events.Sell(config.ProductCode, now, 2000, 0.01)
		if ok || len(events.Signals) != 4 {
			t.Fatal("Failed to disable to add sell event")
		}
	})

	t.Run("test EstimateProfit", func(t *testing.T) {
		profit := events.EstimateProfit()
		if profit != 13 {
			t.Fatalf("Wrong profit estimated. Expected 13 but %f", profit)
		}
	})
}

func deleteSignalEventAll(tx DB) error {
	cmd := "DELETE FROM signal_events"
	_, err := tx.Exec(cmd)
	return err
}
