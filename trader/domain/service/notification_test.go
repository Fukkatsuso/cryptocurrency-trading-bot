package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/slack"
)

func TestNotificationService(t *testing.T) {
	notificationRepository := slack.NewSlackNotificationMockRepository(config.LocalTime)
	notificationService := service.NewNotificationService(notificationRepository)

	t.Run("notify of trading success", func(t *testing.T) {
		event := model.NewSignalEvent(time.Now(), config.ProductCode, model.OrderSideBuy, 1000, 0.1)
		err := notificationService.NotifyOfTradingSuccess(*event)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("notify of trading failed", func(t *testing.T) {
		msg := errors.New("test of NotifyOfTradingFailure")
		err := notificationService.NotifyOfTradingFailed(config.ProductCode, msg)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
