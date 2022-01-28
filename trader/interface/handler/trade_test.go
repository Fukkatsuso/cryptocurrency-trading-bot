package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/slack"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

func TestTradeHandler(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)
	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)
	balanceRepository := bitflyer.NewBitFlyerBalanceMockRepository()
	tickerRepository := bitflyer.NewBitflyerTickerMockRepository()
	orderRepository := bitflyer.NewBitflyerOrderMockRepository()
	notificationRepository := slack.NewSlackNotificationMockRepository(config.LocalTime)

	signalEventService := service.NewSignalEventService(signalEventRepository)
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewMRBaseDataFrameService(indicatorService)
	tradeParamsService := service.NewTradeParamsService(tradeParamsRepository, dataFrameService)
	tradeService := service.NewTradeService(balanceRepository, tickerRepository, orderRepository, signalEventRepository, candleService, dataFrameService, tradeParamsService)
	notificationService := service.NewNotificationService(notificationRepository)

	tradeUsecase := usecase.NewTradeUsecase(signalEventService, tradeService, notificationService)

	tradeHandler := handler.NewTradeHandler(tradeUsecase)

	// 取引パラメータを用意しておく
	params := model.NewBasicTradeParams(config.ProductCode, 0.01)
	tradeParamsRepository.Save(*params)

	t.Run("trade", func(t *testing.T) {
		ts := httptest.NewServer(tradeHandler.Trade(config.ProductCode, 365))
		defer ts.Close()

		rec := httptest.NewRecorder()

		resp, err := http.Post(ts.URL, "text/plain", rec.Body)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("resp.StatusCode != http.StatusOK")
		}

		respBody, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(respBody))
	})
}
