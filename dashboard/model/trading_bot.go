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

func (bot *TradingBot) Buy(timeTime time.Time, timeFormat string) (string, error) {
	if !bot.SignalEvents.CanBuy(timeTime) {
		return "", errors.New("[Buy] can't buy due to signal_event's history")
	}

	// 所持中の現金
	_, availableCurrency := bot.APIClient.GetAvailableBalance(bot.CoinCode, bot.CurrencyCode)
	// 現在の価格
	ticker, err := bot.APIClient.GetTicker(bot.ProductCode)
	if err != nil {
		return "", errors.New(fmt.Sprintf("[Buy] couldn't get ticker: %s", err.Error()))
	}

	// 通貨購入サイズ
	size := bot.TradeParams.Size

	// お金が足りないときは購入しない
	needCurrency := ticker.BestAsk * size
	if availableCurrency < needCurrency {
		return "", errors.New(fmt.Sprintf("[Buy] you don't have enough money. available: %f, need: %f", availableCurrency, needCurrency))
	}

	// 注文発行
	order := &bitflyer.Order{
		ProductCode:     bot.ProductCode,
		ChildOrderType:  bitflyer.ChildOrderTypeMarket,
		Side:            bitflyer.OrderSideBuy,
		Size:            size,
		MinuteToExpires: bot.MinuteToExpires,
		TimeInForce:     bitflyer.TimeInForceGTC,
	}
	fmt.Println("[Order]", order)
	resp, err := bot.APIClient.SendOrder(order)
	if err != nil {
		fmt.Println("[Buy]", err)
		return "", errors.New(fmt.Sprintf("[Buy] failed in SendOrder(): %s", err.Error()))
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Buy]", "order send, but child_order_acceptance_id is none")
		return "", errors.New("[Buy] order send, but child_order_acceptance_id is none")
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID := resp.ChildOrderAcceptanceID
	completedOrder := bot.WaitUntilOrderComplete(childOrderAcceptanceID, timeTime)
	fmt.Printf("[Buy] order completed: %v", completedOrder)

	// 注文失敗した場合
	if completedOrder == nil {
		return childOrderAcceptanceID, errors.New("[Buy] failed to complete order")
	}

	// 注文成功した場合
	// SignalEventsに注文記録を追加
	signalEvent := bot.SignalEvents.Buy(bot.ProductCode, timeTime, completedOrder.Price, completedOrder.Size)
	if signalEvent == nil {
		return childOrderAcceptanceID, errors.New("[Sell] order send, but signal_event is nil")
	}
	// SingalEventをDBに保存
	saved := signalEvent.Save(bot.DBClient, timeFormat)

	if !saved {
		return childOrderAcceptanceID, errors.New(fmt.Sprintf("[Buy] couldn't save signal_event: %v", signalEvent))
	}
	return childOrderAcceptanceID, nil
}

func (bot *TradingBot) Sell(timeTime time.Time, timeFormat string) (string, error) {
	if !bot.SignalEvents.CanSell(timeTime) {
		return "", errors.New("[Sell] can't sell due to signal_event's history")
	}

	// 所持中の仮想通貨
	availableCoin, _ := bot.APIClient.GetAvailableBalance(bot.CoinCode, bot.CurrencyCode)

	// 通貨売却サイズ
	size := bot.TradeParams.Size
	// パラメータに設定したサイズよりも保有量が足りないときは保有量だけ使う
	if availableCoin < size {
		size = availableCoin
	}

	// 注文発行
	order := &bitflyer.Order{
		ProductCode:     bot.ProductCode,
		ChildOrderType:  bitflyer.ChildOrderTypeMarket,
		Side:            bitflyer.OrderSideSell,
		Size:            size,
		MinuteToExpires: bot.MinuteToExpires,
		TimeInForce:     bitflyer.TimeInForceGTC,
	}
	fmt.Println("[Order]", order)
	resp, err := bot.APIClient.SendOrder(order)
	if err != nil {
		fmt.Println("[Sell]", err)
		return "", errors.New(fmt.Sprintf("[Sell] failed in SendOrder(): %s", err.Error()))
	}
	if resp.ChildOrderAcceptanceID == "" {
		fmt.Println("[Sell]", "order send, but child_order_acceptance_id is none")
		return "", errors.New("[Sell] order send, but child_order_acceptance_id is none")
	}

	// 注文が完了するまで待つ
	childOrderAcceptanceID := resp.ChildOrderAcceptanceID
	completedOrder := bot.WaitUntilOrderComplete(childOrderAcceptanceID, timeTime)
	fmt.Printf("[Sell] order completed: %v", completedOrder)

	// 注文失敗した場合
	if completedOrder == nil {
		return childOrderAcceptanceID, errors.New("[Sell] failed to complete order")
	}

	// 注文成功した場合
	// SignalEventsに注文記録を追加
	signalEvent := bot.SignalEvents.Sell(bot.ProductCode, timeTime, completedOrder.Price, completedOrder.Size)
	if signalEvent == nil {
		return childOrderAcceptanceID, errors.New("[Sell] order send, but signal_event is nil")
	}
	// SingalEventをDBに保存
	saved := signalEvent.Save(bot.DBClient, timeFormat)

	if !saved {
		return childOrderAcceptanceID, errors.New(fmt.Sprintf("[Sell] couldn't save signal_event: %v", signalEvent))
	}
	return childOrderAcceptanceID, nil
}

// 取引手数料
// 過去30日間の取引量によって変わる
// とりあえずETH/JPYで最大の0.15%としておく
func TradingFee(productCode string, amount float64) float64 {
	return amount * 0.0015
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
				if order.ChildOrderState == bitflyer.OrderStateCompleted {
					if order.Side == bitflyer.OrderSideBuy {
						return &order
					}
					if order.Side == bitflyer.OrderSideSell {
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
		nowTime := time.Now().UTC()
		_, err := bot.Buy(nowTime, timeFormat)
		if err != nil {
			return err
		}
	}

	currentPrice := df.Candles[now].Close
	if sellPoint > 0 || ShouldCutLoss(bot.SignalEvents, currentPrice, params.StopLimitPercent) {
		nowTime := time.Now().UTC()
		_, err := bot.Sell(nowTime, timeFormat)
		if err != nil {
			return err
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
