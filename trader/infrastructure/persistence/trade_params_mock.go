package persistence

import (
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type tradeParamsMockRepository struct {
	latest map[string]*model.TradeParams
}

func NewTradeParamsMockRepository() repository.TradeParamsRepository {
	return &tradeParamsMockRepository{
		latest: map[string]*model.TradeParams{},
	}
}

func (tr *tradeParamsMockRepository) Save(tp model.TradeParams) error {
	tr.latest[tp.ProductCode()] = &tp
	return nil
}

func (tr *tradeParamsMockRepository) Find(productCode string) (*model.TradeParams, error) {
	if params, ok := tr.latest[productCode]; ok {
		return params, nil
	}
	return nil, errors.New("cannot find trade_params")
}
