package usecase_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestDataFrameUsecase(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat)
	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)

	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)

	t.Run("get", func(t *testing.T) {
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
		_, err := dataFrameUsecase.Get(params, 1000, true)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
