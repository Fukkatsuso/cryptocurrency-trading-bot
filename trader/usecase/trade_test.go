package usecase_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/slack"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

func TestTradeUsecase(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat)
	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)
	balanceRepository := bitflyer.NewBitFlyerBalanceMockRepository()
	tickerRepository := bitflyer.NewBitflyerTickerMockRepository()
	orderRepository := bitflyer.NewBitflyerOrderMockRepository()
	notificationRepository := slack.NewSlackNotificationMockRepository(config.LocalTime)

	signalEventService := service.NewSignalEventService(signalEventRepository)
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)
	tradeParamsService := service.NewTradeParamsService(tradeParamsRepository, dataFrameService)
	tradeService := service.NewTradeService(balanceRepository, tickerRepository, orderRepository, signalEventRepository, candleService, dataFrameService, tradeParamsService)
	notificationService := service.NewNotificationService(notificationRepository)

	tradeUsecase := usecase.NewTradeUsecase(signalEventService, tradeService, notificationService)

	// 取引パラメータを用意しておく
	params := model.NewTradeParams(
		true,
		config.ProductCode,
		0.01,
		true,
		7,
		14,
		50,
		true,
		7,
		14,
		50,
		true,
		20,
		2,
		true,
		true,
		14,
		30,
		70,
		true,
		12,
		26,
		9,
		0.75,
	)
	tradeParamsRepository.Save(*params)

	t.Run("trade", func(t *testing.T) {
		err := tradeUsecase.Trade(config.ProductCode, 365)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
