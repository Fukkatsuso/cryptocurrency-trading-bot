package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestNewBalance(t *testing.T) {
	table := []struct {
		currencyCode string
		amount       float64
		available    float64
		ok           bool
	}{
		{
			currencyCode: config.ProductCode,
			amount:       1.0,
			available:    1.0,
			ok:           true,
		},
		{
			currencyCode: "",
			amount:       1.0,
			available:    1.0,
			ok:           false,
		},
		{
			currencyCode: config.ProductCode,
			amount:       -1.0,
			available:    1.0,
			ok:           false,
		},
		{
			currencyCode: config.ProductCode,
			amount:       1.0,
			available:    -1.0,
			ok:           false,
		},
	}

	for _, b := range table {
		balance := model.NewBalance(b.currencyCode, b.amount, b.available)
		if b.ok && balance == nil {
			t.Fatal("NewBalance() returns nil")
		} else if !b.ok && balance != nil {
			t.Fatal("NewBalance() returns not nil")
		}
	}
}
