package handler_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestTradeParams(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)

	tradeParamsUsecase := usecase.NewTradeParamsUsecase(tradeParamsRepository)

	tradeParamsHandler := handler.NewTradeParamsHandler(tradeParamsUsecase)

	// save dammy trade_params
	params := model.NewBasicTradeParams(config.ProductCode, 0.01)
	err := tradeParamsUsecase.Save(*params)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Run("get trade_params", func(t *testing.T) {
		ts := httptest.NewServer(tradeParamsHandler.HandlerFunc())
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		query := req.URL.Query()
		query.Add("productCode", config.ProductCode)
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

		var params dto.TradeParams
		err = json.Unmarshal(respBody, &params)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("post trade_params", func(t *testing.T) {
		ts := httptest.NewServer(tradeParamsHandler.HandlerFunc())
		defer ts.Close()

		// request body
		params := model.NewBasicTradeParams(config.ProductCode, 1)
		paramsDto := dto.ConvertTradeParams(params)
		reqBody, err := json.Marshal(paramsDto)
		if err != nil {
			t.Fatal(err.Error())
		}

		req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Fatal(err.Error())
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("resp.StatusCode != http.StatusOK")
		}
	})
}
