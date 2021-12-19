package repository

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
)

type NotificationRepository interface {
	NotifyOfTradingSuccess(event model.SignalEvent) error
	NotifyOfTradingFailure(productCode string, err error) error
}
