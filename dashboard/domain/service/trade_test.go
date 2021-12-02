package service_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestTradeService(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	balanceRepository := bitflyer.NewBitFlyerBalanceMockRepository()
	tickerRepository := bitflyer.NewBitflyerTickerMockRepository()
	orderRepository := bitflyer.NewBitflyerOrderMockRepository()
	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)
	tradeParamsService := service.NewTradeParamsService(tradeParamsRepository, dataFrameService)
	tradeService := service.NewTradeService(balanceRepository, tickerRepository, orderRepository, signalEventRepository, candleService, dataFrameService, tradeParamsService)

	events := make([]model.SignalEvent, 0)
	signalEvents := model.NewSignalEvents(events)

	productCode := config.ProductCode
	tradeSize := 0.01

	t.Run("buy", func(t *testing.T) {
		nowTime := time.Now().UTC()
		err := tradeService.Buy(signalEvents, productCode, tradeSize, nowTime)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("sell", func(t *testing.T) {
		nowTime := time.Now().UTC()
		err := tradeService.Sell(signalEvents, productCode, tradeSize, nowTime)
		if err != nil {
			t.Fatal(err)
		}
	})

	// 正常系: Trade()の実行時点でTradeParamsが存在する
	params := model.NewTradeParams(
		true,
		productCode,
		tradeSize,
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
		err := tradeService.Trade(productCode, 365)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
