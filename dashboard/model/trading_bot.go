package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/bitflyer"
)

type TradingBot struct {
	APIClient       *bitflyer.Client
	DBClient        DB
	ProductCode     string
	CoinCode        string
	CurrencyCode    string
	Duration        time.Duration
	PastPeriod      int
	SignalEvents    *SignalEvents
	TradeParams     *TradeParams
	MinuteToExpires int
}

func NewTradingBot(db DB, apiKey, apiSecret, productCode string, duration time.Duration, pastPeriod int) *TradingBot {
	// 取引所のAPIクライアント
	apiClient := bitflyer.NewClient(apiKey, apiSecret)

	codes := strings.Split(productCode, "_")

	bot := &TradingBot{
		APIClient:       apiClient,
		DBClient:        db,
		ProductCode:     productCode,
		CoinCode:        codes[0],
		CurrencyCode:    codes[1],
		Duration:        duration,
		PastPeriod:      pastPeriod,
		MinuteToExpires: 1,
	}
	return bot
}

func (bot *TradingBot) Buy(candle Candle, timeFormat string) (string, bool) {
	if !bot.SignalEvents.CanBuy(candle.Time) {
		return "", false
	}

	// 所持中の現金
	_, availableCurrency := bot.APIClient.GetAvailableBalance(bot.CoinCode, bot.CurrencyCode)
	// 現在の価格
	ticker, err := bot.APIClient.GetTicker(bot.ProductCode)
	if err != nil {
		return "", false
	}

	// 通貨購入サイズ
	// とりあえず0.01で固定
	size := 0.01

	// お金が足りないときは購入しない
	if availableCurrency < ticker.BestAsk*size {
		return "", false
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
		return "", false
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Buy]", "order send, but child_order_acceptance_id is none")
		return "", false
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID := resp.ChildOrderAcceptanceID
	completedOrder := bot.WaitUntilOrderComplete(childOrderAcceptanceID, candle.Time)

	// 注文失敗した場合
	if completedOrder == nil {
		return childOrderAcceptanceID, false
	}

	// 注文成功した場合
	// SignalEventsに注文記録を追加
	bot.SignalEvents.Buy(bot.ProductCode, candle.Time, order.AveragePrice, order.Size)
	// SingalEventをDBに保存
	signalEvent := SignalEvent{
		ProductCode: bot.ProductCode,
		Time:        candle.Time,
		Side:        "BUY",
		Price:       order.AveragePrice,
		Size:        order.Size,
	}
	saved := signalEvent.Save(bot.DBClient, timeFormat)

	return childOrderAcceptanceID, saved
}

func (bot *TradingBot) Sell(candle Candle, timeFormat string) (string, bool) {
	if !bot.SignalEvents.CanSell(candle.Time) {
		return "", false
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
		return "", false
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Sell]", "order send, but child_order_acceptance_id is none")
		return "", false
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID := resp.ChildOrderAcceptanceID
	completedOrder := bot.WaitUntilOrderComplete(childOrderAcceptanceID, candle.Time)

	// 注文失敗した場合
	if completedOrder == nil {
		return childOrderAcceptanceID, false
	}

	// 注文成功した場合
	// SignalEventsに注文記録を追加
	bot.SignalEvents.Sell(bot.ProductCode, candle.Time, order.AveragePrice, order.Size)
	// SingalEventをDBに保存
	signalEvent := SignalEvent{
		ProductCode: bot.ProductCode,
		Time:        candle.Time,
		Side:        "SELL",
		Price:       order.AveragePrice,
		Size:        order.Size,
	}
	saved := signalEvent.Save(bot.DBClient, timeFormat)

	return childOrderAcceptanceID, saved
}

func (bot *TradingBot) WaitUntilOrderComplete(childOrderAcceptanceID string, executeTime time.Time) *bitflyer.Order {
	params := map[string]string{
		"product_code":              bot.ProductCode,
		"child_order_acceptance_id": childOrderAcceptanceID,
	}

	// 最長2分待つ
	expire := time.After(2 * time.Minute)
	// 15秒ごとに注文状況をポーリング
	interval := time.Tick(15 * time.Second)

	return func() *bitflyer.Order {
		for {
			select {
			case <-expire:
				return nil
			case <-interval:
				orders, err := bot.APIClient.ListOrder(params)
				if err != nil {
					return nil
				}
				if len(orders) == 0 {
					return nil
				}
				order := orders[0]
				if order.ChildOrderState == "COMPLETED" {
					if order.Side == "BUY" {
						return &order
					}
					if order.Side == "SELL" {
						return &order
					}
					return nil
				}
			}
		}
	}()
}

func (bot *TradingBot) Trade(db DB, candleTableName, timeFormat string) error {
	params := bot.TradeParams
	if params == nil {
		return errors.New("[Trade] TradeParams is nil")
	}

	candles, _ := GetAllCandle(db, candleTableName, timeFormat, bot.ProductCode, bot.Duration, bot.PastPeriod)
	lenCandles := len(candles)

	df := DataFrame{
		ProductCode: bot.ProductCode,
		Candles:     candles,
	}

	if params.EMAEnable {
		ok1 := df.AddEMA(params.EMAPeriod1)
		ok2 := df.AddEMA(params.EMAPeriod2)
		params.EMAEnable = ok1 && ok2
	}

	if params.BBandsEnable {
		ok := df.AddBBands(params.BBandsN, params.BBandsK)
		params.BBandsEnable = ok
	}

	if params.IchimokuEnable {
		ok := df.AddIchimoku()
		params.IchimokuEnable = ok
	}

	if params.MACDEnable {
		ok := df.AddMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod)
		params.MACDEnable = ok
	}

	if params.RSIEnable {
		ok := df.AddRSI(params.RSIPeriod)
		params.RSIEnable = ok
	}

	now := lenCandles - 1
	buyPoint, sellPoint := df.Analyze(now, params)

	if buyPoint > 0 {
		childOrderAcceptanceID, isOrderCompleted := bot.Buy(df.Candles[now], timeFormat)
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] buy order is not completed, id=%s", childOrderAcceptanceID))
		}
	}

	currentPrice := df.Candles[now].Close
	if sellPoint > 0 || ShouldCutLoss(bot.SignalEvents, currentPrice, params.StopLimitPercent) {
		childOrderAcceptanceID, isOrderCompleted := bot.Sell(df.Candles[now], timeFormat)
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] sell order is not completed, id=%s", childOrderAcceptanceID))
		}
		bot.OptimizeTradeParams(candleTableName, timeFormat)
	}

	return nil
}

func (bot *TradingBot) OptimizeTradeParams(candleTableName, timeFormat string) {
	candles, _ := GetAllCandle(bot.DBClient, candleTableName, timeFormat, bot.ProductCode, bot.Duration, bot.PastPeriod)
	df := DataFrame{
		ProductCode: bot.ProductCode,
		Candles:     candles,
	}

	optimizedParams := df.OptimizeTradeParams(bot.TradeParams)
	if *optimizedParams != *bot.TradeParams {
		optimizedParams.Create(bot.DBClient)
	}
	bot.TradeParams = optimizedParams
}
