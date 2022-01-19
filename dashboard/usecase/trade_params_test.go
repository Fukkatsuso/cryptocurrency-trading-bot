package usecase_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func TestTradeParams(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	tradeParamsRepository := persistence.NewTradeParamsRepository(tx)

	tradeParamsUsecase := usecase.NewTradeParamsUsecase(tradeParamsRepository)

	t.Run("save trade_params", func(t *testing.T) {
		params := model.NewBasicTradeParams(config.ProductCode, 0.01)
		err := tradeParamsUsecase.Save(*params)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("get trade_params", func(t *testing.T) {
		params, err := tradeParamsUsecase.Get(config.ProductCode)
		if err != nil {
			t.Fatal(err.Error())
		}

		compareParams := model.NewBasicTradeParams(config.ProductCode, 0.01)
		if *params != *compareParams {
			t.Fatalf("%+v != %+v", *params, *compareParams)
		}
	})
}
