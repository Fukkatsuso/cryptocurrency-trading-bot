package slack

import (
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type slackNotificationMockRepository struct {
	timeLocation *time.Location
}

func NewSlackNotificationMockRepository(timeLocation *time.Location) repository.NotificationRepository {
	return &slackNotificationMockRepository{
		timeLocation: timeLocation,
	}
}

func (snr *slackNotificationMockRepository) NotifyOfTradingSuccess(event model.SignalEvent) error {
	timeString := event.Time().In(snr.timeLocation).Format("2006-01-02 15:04:05")

	msg := buildTextMessage(
		fmt.Sprintf("%s *%s*: %s", EmojiCoin, event.Side(), event.ProductCode()),
		fmt.Sprintf("At: %s", timeString),
		fmt.Sprintf("Price: %f", event.Price()),
		fmt.Sprintf("Size: %f", event.Size()),
	)

	fmt.Println(msg)
	return nil
}

func (snr *slackNotificationMockRepository) NotifyOfTradingFailure(productCode string, err error) error {
	msg := buildTextMessage(
		fmt.Sprintf("%s（%s）", EmojiDizzyFace, productCode),
		"エラーが生じました",
		"```",
		err.Error(),
		"```",
	)

	fmt.Println(msg)
	return nil
}
