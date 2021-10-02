package service

import (
	"errors"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type SignalEventService interface {
	Save(event model.SignalEvent) error
	FindAll(productCode string) (*model.SignalEvents, error)
	FindAllAfterTime(productCode string, timeTime time.Time) (*model.SignalEvents, error)
}

type signalEventService struct {
	signalEventRepository repository.SignalEventRepository
}

func NewSignalEventService(sr repository.SignalEventRepository) SignalEventService {
	return &signalEventService{
		signalEventRepository: sr,
	}
}

func (ss *signalEventService) Save(event model.SignalEvent) error {
	return ss.signalEventRepository.Save(event)
}

func (ss *signalEventService) FindAll(productCode string) (*model.SignalEvents, error) {
	signals, err := ss.signalEventRepository.FindAll(productCode)
	if err != nil {
		return nil, err
	}

	signalEvents := model.NewSignalEvents(signals)
	if signalEvents == nil {
		return nil, errors.New("Failed to create a new SignalEvents instance")
	}

	return signalEvents, nil
}

func (ss *signalEventService) FindAllAfterTime(productCode string, timeTime time.Time) (*model.SignalEvents, error) {
	signals, err := ss.signalEventRepository.FindAllAfterTime(productCode, timeTime)
	if err != nil {
		return nil, err
	}

	signalEvents := model.NewSignalEvents(signals)
	if signalEvents == nil {
		return nil, errors.New("Failed to create a new SignalEvents instance")
	}

	return signalEvents, nil
}
