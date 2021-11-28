package slack

import (
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
	"github.com/slack-go/slack"
)

type slackNotificationRepository struct {
	client       *Client
	timeLocation *time.Location
}

func NewSlackNotificationRepository(client *Client, timeLocation *time.Location) repository.NotificationRepository {
	return &slackNotificationRepository{
		client:       client,
		timeLocation: timeLocation,
	}
}

func (snr *slackNotificationRepository) NotifyOfTradingSuccess(event model.SignalEvent) error {
	timeString := event.Time().In(snr.timeLocation).Format("2006-01-02 15:04:05")

	msg := buildTextMessage(
		fmt.Sprintf("%s *%s*: %s", EmojiCoin, event.Side(), event.ProductCode()),
		fmt.Sprintf("At: %s", timeString),
		fmt.Sprintf("Price: %f", event.Price()),
		fmt.Sprintf("Size: %f", event.Size()),
	)

	option := slack.MsgOptionText(msg, true)
	_, _, err := snr.client.client.PostMessage(snr.client.channelId, option)
	return err
}

func (snr *slackNotificationRepository) NotifyOfTradingFailure(productCode string, err error) error {
	msg := buildTextMessage(
		fmt.Sprintf("%s（%s）", EmojiDizzyFace, productCode),
		"エラーが生じました",
		"```",
		err.Error(),
		"```",
	)

	option := slack.MsgOptionText(msg, true)
	_, _, err = snr.client.client.PostMessage(snr.client.channelId, option)
	return err
}

func buildTextMessage(lines ...string) string {
	return strings.Join(lines, "\n")
}
