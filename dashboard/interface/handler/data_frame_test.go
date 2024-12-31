package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestDataFrameHandler(t *testing.T) {
	tx := persistence.NewSQLiteTransaction(config.DSN())
	defer tx.Rollback()

	candleRepository := persistence.NewCandleMockRepository(config.CandleTableName, config.TimeFormat, config.ProductCode, config.CandleDuration)
	signalEventRepository := persistence.NewSignalEventRepository(tx, config.TimeFormat)

	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewMRBaseDataFrameService(indicatorService)

	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)

	dataFrameHandler := handler.NewDataFrameHandler(dataFrameUsecase)

	t.Run("get", func(t *testing.T) {
		ts := httptest.NewServer(dataFrameHandler.Get(config.ProductCode))
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		query := req.URL.Query()
		query.Add("size", "0.01")
		query.Add("sma", "true")
		query.Add("smaPeriod1", "7")
		query.Add("smaPeriod2", "14")
		query.Add("smaPeriod3", "50")
		query.Add("ema", "true")
		query.Add("emaPeriod1", "7")
		query.Add("emaPeriod2", "14")
		query.Add("emaPeriod3", "50")
		query.Add("bbands", "true")
		query.Add("bbandsN", "20")
		query.Add("bbandsK", "2")
		query.Add("ichimoku", "true")
		query.Add("rsi", "true")
		query.Add("rsiPeriod", "14")
		query.Add("rsiBuyThread", "30")
		query.Add("rsiSellThread", "70")
		query.Add("macd", "true")
		query.Add("macdPeriod1", "12")
		query.Add("macdPeriod2", "26")
		query.Add("macdPeriod3", "9")
		query.Add("stopLimitPercent", "0.75")
		query.Add("limit", "365")
		req.URL.RawQuery = query.Encode()

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("resp.StatusCode != http.StatusOK")
		}

		respBody, _ := ioutil.ReadAll(resp.Body)

		var df dto.DataFrame
		err = json.Unmarshal(respBody, &df)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
