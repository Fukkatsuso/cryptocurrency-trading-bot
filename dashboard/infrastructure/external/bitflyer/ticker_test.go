package bitflyer_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/bitflyer"
)

func TestBitflyerTickerRepository(t *testing.T) {
	apiClient := bitflyer.NewClient(config.APIKey, config.APISecret)
	tickerRepository := bitflyer.NewBitflyerTickerRepository(apiClient)

	t.Run("fetch", func(t *testing.T) {
		ticker, err := tickerRepository.Fetch(config.ProductCode)
		// 外部APIを利用するため，予期せずfetchできない場合がある
		if err != nil {
			t.Skip(err.Error())
		}
		t.Log(ticker)
	})
}
