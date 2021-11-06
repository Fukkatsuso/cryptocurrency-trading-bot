package service_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
)

func TestSignalEventService(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)
	signalEventService := service.NewSignalEventService(signalEventRepository)

	signalTime := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("save signal_event", func(t *testing.T) {
		event := model.NewSignalEvent(signalTime, config.ProductCode, model.OrderSideBuy, 100000, 0.1)
		err := signalEventService.Save(*event)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("find all signal_event", func(t *testing.T) {
		events, err := signalEventService.FindAll(config.ProductCode)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(events) < 1 {
			t.Fatal("len(events) < 1")
		}
	})

	t.Run("sind all after time", func(t *testing.T) {
		events, err := signalEventRepository.FindAllAfterTime(config.ProductCode, signalTime)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(events) != 1 {
			t.Fatal("len(events) != 1")
		}
	})
}
