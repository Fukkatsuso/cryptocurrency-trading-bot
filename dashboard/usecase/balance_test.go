package usecase_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestBalance(t *testing.T) {
	balanceRepository := bitflyer.NewBitFlyerBalanceMockRepository()

	balanceUsecase := usecase.NewBalanceUsecase(balanceRepository)

	t.Run("get balance", func(t *testing.T) {
		balances, err := balanceUsecase.Get()
		if err != nil {
			t.Fatal(err.Error())
		}
		if balances == nil {
			t.Fatal("Get() returns nil")
		}
	})
}
