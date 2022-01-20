package handler_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler/dto"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestBalance(t *testing.T) {
	balanceRepository := bitflyer.NewBitFlyerBalanceMockRepository()

	balanceUsecase := usecase.NewBalanceUsecase(balanceRepository)

	balanceHandler := handler.NewBalanceHandler(balanceUsecase)

	t.Run("get balance", func(t *testing.T) {
		ts := httptest.NewServer(balanceHandler.Get())
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL, nil)
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

		respBody, _ := ioutil.ReadAll(resp.Body)

		var balance []dto.Balance
		err = json.Unmarshal(respBody, &balance)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
