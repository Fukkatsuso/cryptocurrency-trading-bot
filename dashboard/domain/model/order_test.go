package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestOrder(t *testing.T) {
	t.Run("new buy order", func(t *testing.T) {
		var order *model.Order

		order = model.NewBuyOrder(config.ProductCode, 1)
		if order == nil {
			t.Fatal("NewBuyOrder() returns nil")
		}

		order = model.NewBuyOrder("", 1)
		if order != nil {
			t.Fatal("NewBuyOrder() returns not nil")
		}

		order = model.NewBuyOrder(config.ProductCode, -1)
		if order != nil {
			t.Fatal("NewBuyOrder() returns not nil")
		}
	})

	t.Run("new sell order", func(t *testing.T) {
		var order *model.Order

		order = model.NewSellOrder(config.ProductCode, 1)
		if order == nil {
			t.Fatal("NewSellOrder() returns nil")
		}

		order = model.NewSellOrder("", 1)
		if order != nil {
			t.Fatal("NewSellOrder() returns not nil")
		}

		order = model.NewSellOrder(config.ProductCode, -1)
		if order != nil {
			t.Fatal("NewSellOrder() returns not nil")
		}
	})
}
