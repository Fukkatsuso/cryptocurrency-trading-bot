package repository

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type BalanceRepository interface {
	FetchAll() ([]model.Balance, error)
	FetchByCurrencyCode(currencyCode string) (*model.Balance, error)
}
