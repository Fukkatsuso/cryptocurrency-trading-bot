package slack_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/slack"
)

func TestSlackNotificationRepository(t *testing.T) {
	client := slack.NewClient(config.SlackBotToken, config.SlackChannelID)
	notificationRepository := slack.NewSlackNotificationRepository(client, config.LocalTime)

	t.Run("notify of trading success", func(t *testing.T) {
		event := model.NewSignalEvent(time.Now(), config.ProductCode, model.OrderSideBuy, 1000, 0.1)
		err := notificationRepository.NotifyOfTradingSuccess(*event)
		if err != nil {
			t.Skip(err)
		}
	})

	t.Run("notify of trading failure", func(t *testing.T) {
		msg := errors.New("test of NotifyOfTradingFailure")
		err := notificationRepository.NotifyOfTradingFailure(config.ProductCode, msg)
		if err != nil {
			t.Skip(err)
		}
	})
}
