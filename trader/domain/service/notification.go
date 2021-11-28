package service

import (
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type NotificationService interface {
	NotifyOfTradingSuccess(event model.SignalEvent) error
	NotifyOfTradingFailed(productCode string, err error) error
}

type notificationService struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationService(nr repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepository: nr,
	}
}

func (ns *notificationService) NotifyOfTradingSuccess(event model.SignalEvent) error {
	return ns.notificationRepository.NotifyOfTradingSuccess(event)
}

func (ns *notificationService) NotifyOfTradingFailed(productCode string, err error) error {
	return ns.notificationRepository.NotifyOfTradingFailure(productCode, err)
}
