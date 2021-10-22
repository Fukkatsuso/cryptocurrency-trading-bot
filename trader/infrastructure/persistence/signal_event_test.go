package persistence_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
)

// 日時は2100年1月1日以降かつ昇順
func newSignalEvents() []model.SignalEvent {
	table := []struct {
		time        time.Time
		productCode string
		side        string
		price       float64
		size        float64
	}{
		{
			time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC),
			config.ProductCode,
			string(model.OrderSideBuy),
			1000.0,
			0.01,
		},
		{
			time.Date(2100, 1, 2, 0, 0, 0, 0, time.UTC),
			config.ProductCode,
			string(model.OrderSideSell),
			1500.0,
			0.01,
		},
	}

	signalEvents := make([]model.SignalEvent, 0)
	for _, t := range table {
		signalEvent := model.NewSignalEvent(t.time, t.productCode, model.OrderSide(t.side), t.price, t.size)
		if signalEvent == nil {
			continue
		}
		signalEvents = append(signalEvents, *signalEvent)
	}
	return signalEvents
}

func TestSignalEvent(t *testing.T) {
	tx := NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)

	signalEvents := newSignalEvents()

	t.Run("save signal_event", func(t *testing.T) {
		for _, signalEvent := range signalEvents {
			err := signalEventRepository.Save(signalEvent)
			if err != nil {
				t.Fatal(err.Error())
			}
		}
	})

	t.Run("find all signal_event", func(t *testing.T) {
		ss, err := signalEventRepository.FindAll(config.ProductCode)
		if err != nil {
			t.Fatal(err.Error())
		}
		if ss == nil || len(ss) < len(signalEvents) {
			t.Fatal("FingAll() returns incomplete data")
		}
	})

	t.Run("find signal_event after time", func(t *testing.T) {
		criteriaTime := signalEvents[0].Time().Add(time.Second)
		ss, err := signalEventRepository.FindAllAfterTime(config.ProductCode, criteriaTime)
		if err != nil {
			t.Fatal(err.Error())
		}
		if ss == nil || len(ss) != len(signalEvents)-1 {
			t.Fatal("FindAllAfterTime() returns incomplete data")
		}
	})
}
