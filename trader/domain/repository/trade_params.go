package repository

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"

type TradeParamsRepository interface {
	Save(tp model.TradeParams) error
	Find(productCode string) (*model.TradeParams, error)
}
