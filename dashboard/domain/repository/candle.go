package repository

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

type CandleRepository interface {
	Save(candle model.Candle) error
	FindByCandleTime(productCode string, duration time.Duration, timeTime model.CandleTime) (*model.Candle, error)
	FindAll(productCode string, duration time.Duration, limit int64) ([]model.Candle, error)
}
