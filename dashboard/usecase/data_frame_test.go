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
	tx := persistence.NewSQLiteTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)

	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)

	t.Run("get", func(t *testing.T) {
		params := model.NewBasicTradeParams(config.ProductCode, 0.01)
		_, err := dataFrameUsecase.Get(params, 1000, true)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
