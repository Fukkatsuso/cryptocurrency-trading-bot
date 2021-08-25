package controller

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/model"
)

func TestSlackNotifySignalEvent(t *testing.T) {
	signal := &model.SignalEvent{
		Time:        time.Now().UTC(),
		ProductCode: config.ProductCode,
		Side:        string(bitflyer.OrderSideBuy),
		Price:       1000,
		Size:        0.1,
	}
	err := SlackNotifySignalEvent(signal)
	if err != nil {
		t.Fatal(err.Error())
	}
}
