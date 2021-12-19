package repository

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
)

type TickerRepository interface {
	Fetch(productCode string) (*model.Ticker, error)
}
