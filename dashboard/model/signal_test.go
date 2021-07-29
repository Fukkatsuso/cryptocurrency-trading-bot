package model

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
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
