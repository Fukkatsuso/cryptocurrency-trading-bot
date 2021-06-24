package model

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/bitflyer"
)

type TradingBot struct {
	APIClient       *bitflyer.Client
	ProductCode     string
	Duration        time.Duration
	PastPeriod      int
	SignalEvents    *SignalEvents
	TradeParams     *TradeParams
	MinuteToExpires int
}

func NewTradingBot(productCode string, duration time.Duration, pastPeriod int) *TradingBot {
	// 取引所のAPIクライアント
	apiClient := bitflyer.NewClient(config.APIKey, config.APISecret)

	bot := &TradingBot{
		APIClient:       apiClient,
		ProductCode:     productCode,
		Duration:        duration,
		PastPeriod:      pastPeriod,
		MinuteToExpires: 1,
	}
	return bot
}
