package handler_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

func TestCandleHandler(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleRepository(tx, config.CandleTableName, config.TimeFormat)
	tickerRepository := bitflyer.NewBitflyerTickerMockRepository()

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)

	candleUsecase := usecase.NewCandleUsecase(candleService, tickerRepository)

	candleHandler := handler.NewCandleHandler(candleUsecase)

	t.Run("update candle", func(t *testing.T) {
		ts := httptest.NewServer(candleHandler.UpdateCandle(config.ProductCode))
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
