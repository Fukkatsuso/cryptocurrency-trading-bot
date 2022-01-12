package usecase

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type TradeParamsUsecase interface {
	Get(productCode string) (*model.TradeParams, error)
	Save(params model.TradeParams) error
}

type tradeParamsUsecase struct {
	tradeParamsRepository repository.TradeParamsRepository
}

func NewTradeParamsUsecase(tr repository.TradeParamsRepository) TradeParamsUsecase {
	return &tradeParamsUsecase{
		tradeParamsRepository: tr,
	}
}

func (tu *tradeParamsUsecase) Get(productCode string) (*model.TradeParams, error) {
	return tu.tradeParamsRepository.Find(productCode)
}

func (tu *tradeParamsUsecase) Save(params model.TradeParams) error {
	return tu.tradeParamsRepository.Save(params)
}
