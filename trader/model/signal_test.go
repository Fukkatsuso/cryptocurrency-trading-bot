package model

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
)

func TestSignalEventSave(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	// 2100/01/01 00:00:00.00 UTC
	timeDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	event := &SignalEvent{
		Time:        timeDate,
		ProductCode: config.ProductCode,
		Side:        "BUY",
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

	t.Run("save some signal_event", func(t *testing.T) {
		// 2021/01/01 00:00:00.00 UTC
		timeDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		event := &SignalEvent{
			Time:        timeDate,
			ProductCode: config.ProductCode,
			Side:        "BUY",
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
			Side:        "SELL",
			Price:       1500,
			Size:        0.01,
		}
		ok = event.Save(tx, config.TimeFormat)
		if !ok {
			t.Fatal("Failed to Save SignalEvent:", event)
		}
	})

	var events *SignalEvents

	t.Run("get signal_events by product_code", func(t *testing.T) {
		events = GetSignalEventsByProductCode(tx, config.ProductCode)

		if events == nil {
			t.Fatal("Failed to GetSignalEventsByProductCode")
		}

		if len(events.Signals) != 2 {
			t.Fatalf("wrong number of SignalEvents. Expected 2, but %d", len(events.Signals))
		}
	})
}

func deleteSignalEventAll(tx DB) error {
	cmd := "DELETE FROM signal_events"
	_, err := tx.Exec(cmd)
	return err
}
