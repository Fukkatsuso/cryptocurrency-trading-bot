package usecase_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

func TestCandleUsecase(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleRepository(tx, config.CandleTableName, config.TimeFormat)
	tickerRepository := bitflyer.NewBitflyerTickerMockRepository()

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)

	candleUsecase := usecase.NewCandleUsecase(candleService, tickerRepository)

	t.Run("update candle", func(t *testing.T) {
		err := candleUsecase.UpdateCandle(config.ProductCode)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
