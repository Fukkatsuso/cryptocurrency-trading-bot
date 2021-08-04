package model

import (
	"fmt"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
)

func TestTradeParams(t *testing.T) {
	tx := NewTransaction(config.DSN())
	defer tx.Rollback()

	err := deleteTradeParamsAll(tx)
	if err != nil {
		t.Fatal("failed to exec deleteTradeParamsAll")
	}

	var tradeParams *TradeParams

	t.Run("create trade_params", func(t *testing.T) {
		tradeParams = &TradeParams{
			TradeEnable: true,
			ProductCode: config.ProductCode,
			Size:        0.01,
		}

		err := tradeParams.Create(tx, config.TradeParamTableName)
		if err != nil {
			t.Fatal("Failed to Create TradeParams:", tradeParams, err.Error())
		}
	})

	t.Run("get trade_params", func(t *testing.T) {
		getTradeParams := GetTradeParams(tx, config.TradeParamTableName, config.ProductCode)

		if getTradeParams == nil {
			t.Fatal("Failed to Get TradeParams")
		}
		if *getTradeParams != *tradeParams {
			t.Fatalf("%v != %v", *getTradeParams, *tradeParams)
		}
	})
}

func deleteTradeParamsAll(tx DB) error {
	cmd := fmt.Sprintf("DELETE FROM %s", config.TradeParamTableName)
	_, err := tx.Exec(cmd)
	return err
}
