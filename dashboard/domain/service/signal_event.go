package service

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type SignalEventService interface {
	Save(event model.SignalEvent) error
	FindAll(productCode string) ([]model.SignalEvent, error)
	FindAllAfterTime(productCode string, timeTime time.Time) ([]model.SignalEvent, error)
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

func (ss *signalEventService) FindAll(productCode string) ([]model.SignalEvent, error) {
	signals, err := ss.signalEventRepository.FindAll(productCode)
	if err != nil {
		return nil, err
	}

	return signals, nil
}

func (ss *signalEventService) FindAllAfterTime(productCode string, timeTime time.Time) ([]model.SignalEvent, error) {
	signals, err := ss.signalEventRepository.FindAllAfterTime(productCode, timeTime)
	if err != nil {
		return nil, err
	}

	return signals, nil
}
