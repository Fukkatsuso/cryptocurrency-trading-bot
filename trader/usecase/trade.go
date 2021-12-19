package usecase

import (
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
)

type TradeUsecase interface {
	Trade(productCode string, pastPeriod int) error
}

type tradeUsecase struct {
	signalEventService  service.SignalEventService
	tradeService        service.TradeService
	notificationService service.NotificationService
}

func NewTradeUsecase(ss service.SignalEventService, ts service.TradeService, ns service.NotificationService) TradeUsecase {
	return &tradeUsecase{
		signalEventService:  ss,
		tradeService:        ts,
		notificationService: ns,
	}
}

func (tu *tradeUsecase) Trade(productCode string, pastPeriod int) error {
	// 取引前の時刻
	// 通知発生基準にする
	beforeTradeTime := time.Now().UTC()

	err := tu.tradeService.Trade(productCode, pastPeriod)
	if err != nil {
		return err
	}

	// 今回実行した取引を取得
	events, err := tu.signalEventService.FindAllAfterTime(productCode, beforeTradeTime)
	if err != nil {
		return err
	}

	// 通知
	for _, event := range events {
		err := tu.notificationService.NotifyOfTradingSuccess(event)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
