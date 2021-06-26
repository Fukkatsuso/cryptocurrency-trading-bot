package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/bitflyer"
)

type TradingBot struct {
	APIClient       *bitflyer.Client
	ProductCode     string
	CoinCode        string
	CurrencyCode    string
	Duration        time.Duration
	PastPeriod      int
	SignalEvents    *SignalEvents
	TradeParams     *TradeParams
	MinuteToExpires int
}

func NewTradingBot(productCode string, duration time.Duration, pastPeriod int) *TradingBot {
	// 取引所のAPIクライアント
	apiClient := bitflyer.NewClient(config.APIKey, config.APISecret)

	codes := strings.Split(productCode, "_")

	bot := &TradingBot{
		APIClient:       apiClient,
		ProductCode:     productCode,
		CoinCode:        codes[0],
		CurrencyCode:    codes[1],
		Duration:        duration,
		PastPeriod:      pastPeriod,
		MinuteToExpires: 1,
	}
	return bot
}

func (bot *TradingBot) Buy(candle Candle) (childOrderAcceptanceID string, isOrderCompleted bool) {
	if !bot.SignalEvents.CanBuy(candle.Time) {
		return
	}

	// 所持中の現金
	_, availableCurrency := bot.APIClient.GetAvailableBalance(bot.CoinCode, bot.CurrencyCode)
	// 現在の価格
	ticker, err := bot.APIClient.GetTicker(bot.ProductCode)
	if err != nil {
		return
	}

	// 通貨購入サイズ
	// とりあえず0.01で固定
	size := 0.01

	// お金が足りないときは購入しない
	if availableCurrency < ticker.BestAsk*size {
		return
	}

	// 注文発行
	order := &bitflyer.Order{
		ProductCode:     bot.ProductCode,
		ChildOrderType:  "MARKET",
		Side:            "BUY",
		Size:            size,
		MinuteToExpires: bot.MinuteToExpires,
		TimeInForce:     "GTC",
	}
	fmt.Println("[Order]", order)
	resp, err := bot.APIClient.SendOrder(order)
	if err != nil {
		fmt.Println("[Buy]", err)
		return
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Buy]", "order send, but child_order_acceptance_id is none")
		return
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID = resp.ChildOrderAcceptanceID
	isOrderCompleted = bot.WaitUntilOrderComplete(childOrderAcceptanceID, candle.Time)

	return childOrderAcceptanceID, isOrderCompleted
}

func (bot *TradingBot) Sell(candle Candle) (childOrderAcceptanceID string, isOrderCompleted bool) {
	if !bot.SignalEvents.CanSell(candle.Time) {
		return
	}

	// 所持中の仮想通貨
	availableCoin, _ := bot.APIClient.GetAvailableBalance(bot.CoinCode, bot.CurrencyCode)

	// 通貨売却サイズ
	size := availableCoin

	// 注文発行
	order := &bitflyer.Order{
		ProductCode:     bot.ProductCode,
		ChildOrderType:  "MARKET",
		Side:            "SELL",
		Size:            size,
		MinuteToExpires: bot.MinuteToExpires,
		TimeInForce:     "GTC",
	}
	fmt.Println("[Order]", order)
	resp, err := bot.APIClient.SendOrder(order)
	if err != nil {
		fmt.Println("[Sell]", err)
		return
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Sell]", "order send, but child_order_acceptance_id is none")
		return
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID = resp.ChildOrderAcceptanceID
	isOrderCompleted = bot.WaitUntilOrderComplete(childOrderAcceptanceID, candle.Time)

	return childOrderAcceptanceID, isOrderCompleted
}

func (bot *TradingBot) WaitUntilOrderComplete(childOrderAcceptanceID string, executeTime time.Time) bool {
	params := map[string]string{
		"product_code":              bot.ProductCode,
		"child_order_acceptance_id": childOrderAcceptanceID,
	}

	// 最長2分待つ
	expire := time.After(2 * time.Minute)
	// 15秒ごとに注文状況をポーリング
	interval := time.Tick(15 * time.Second)

	return func() bool {
		for {
			select {
			case <-expire:
				return false
			case <-interval:
				orders, err := bot.APIClient.ListOrder(params)
				if err != nil {
					return false
				}
				if len(orders) == 0 {
					return false
				}
				order := orders[0]
				if order.ChildOrderState == "COMPLETED" {
					if order.Side == "BUY" {
						couldBuy := bot.SignalEvents.Buy(bot.ProductCode, executeTime, order.AveragePrice, order.Size)
						if !couldBuy {
							fmt.Printf("[Buy] could not buy: childOrderAcceptanceID=%s order=%+v\n", childOrderAcceptanceID, order)
						}
						return couldBuy
					}
					if order.Side == "SELL" {
						couldSell := bot.SignalEvents.Sell(bot.ProductCode, executeTime, order.AveragePrice, order.Size)
						if !couldSell {
							fmt.Printf("[Buy] could not sell: childOrderAcceptanceID=%s order=%+v\n", childOrderAcceptanceID, order)
						}
						return couldSell
					}
					return false
				}
			}
		}
	}()
}
