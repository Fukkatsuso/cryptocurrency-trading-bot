package repository

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

type SignalEventRepository interface {
	Save(signal model.SignalEvent) error
	FindAll(productCode string) ([]model.SignalEvent, error)
	FindAllAfterTime(productCode string, timeTime time.Time) ([]model.SignalEvent, error)
}
