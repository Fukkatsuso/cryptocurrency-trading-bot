package bitflyer_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
)

func TestBitFlyerBalanceRepository(t *testing.T) {
	apiClient := bitflyer.NewClient(config.APIKey, config.APISecret)
	balanceRepository := bitflyer.NewBitFlyerBalanceRepository(apiClient)

	t.Run("fetch all", func(t *testing.T) {
		balances, err := balanceRepository.FetchAll()
		if err != nil {
			t.Skip(err.Error())
		}
		t.Log(balances)
	})

	t.Run("fetch by currency code", func(t *testing.T) {
		table := []string{
			"BTC",
			"ETH",
			"JPY",
		}

		for _, c := range table {
			balance, err := balanceRepository.FetchByCurrencyCode(c)
			if err != nil {
				t.Skip(c, err.Error())
			}
			t.Log(c, balance)
		}
	})
}
