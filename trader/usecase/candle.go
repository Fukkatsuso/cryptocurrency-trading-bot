package usecase

import (
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
)

type CandleUsecase interface {
	UpdateCandle(productCode string) error
}

type candleUsecase struct {
	candleService    service.CandleService
	tickerRepository repository.TickerRepository
}

func NewCandleUsecase(cs service.CandleService, tr repository.TickerRepository) *candleUsecase {
	return &candleUsecase{
		candleService:    cs,
		tickerRepository: tr,
	}
}

func (cu *candleUsecase) UpdateCandle(productCode string) error {
	// ticker取得
	ticker, err := cu.tickerRepository.Fetch(productCode)
	if err != nil {
		return err
	}

	// ticker -> candle
	candle := cu.candleService.TickerToCandle(*ticker)
	if candle == nil {
		return errors.New("Failed to convert ticker into candle")
	}

	// 最新のcandle
	currentCandle, err := cu.candleService.FindByTime(productCode, candle.Time().Time())
	if err != nil {
		return err
	}

	// candleを更新して保存
	newCandle := cu.candleService.Update(currentCandle, candle)
	if newCandle == nil {
		return errors.New("Failed to update candle")
	}
	return cu.candleService.Save(*newCandle)
}
