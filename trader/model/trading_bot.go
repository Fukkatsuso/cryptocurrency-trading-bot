package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
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

	var emaValues1 []float64
	var emaValues2 []float64
	if params.EMAEnable {
		ok1 := df.AddEMA(params.EMAPeriod1)
		ok2 := df.AddEMA(params.EMAPeriod2)
		params.EMAEnable = ok1 && ok2
	}
	if params.EMAEnable {
		emaValues1 = df.EMAs[0].Values
		emaValues2 = df.EMAs[1].Values
	}

	var bbUp []float64
	var bbDown []float64
	if params.BBandsEnable {
		ok := df.AddBBands(params.BBandsN, params.BBandsK)
		params.BBandsEnable = ok
	}
	if params.BBandsEnable {
		bbUp = df.BBands.Up
		bbDown = df.BBands.Down
	}

	var tenkan, kijun, senkouA, senkouB, chikou []float64
	if params.IchimokuEnable {
		ok := df.AddIchimoku()
		params.IchimokuEnable = ok
	}
	if params.IchimokuEnable {
		tenkan = df.IchimokuCloud.Tenkan
		kijun = df.IchimokuCloud.Kijun
		senkouA = df.IchimokuCloud.SenkouA
		senkouB = df.IchimokuCloud.SenkouB
		chikou = df.IchimokuCloud.Chikou
	}

	var outMACD, outMACDSignal []float64
	if params.MACDEnable {
		ok := df.AddMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod)
		params.MACDEnable = ok
	}
	if params.MACDEnable {
		outMACD = df.MACD.MACD
		outMACDSignal = df.MACD.MACDSignal
	}

	var rsiValues []float64
	if params.RSIEnable {
		ok := df.AddRSI(params.RSIPeriod)
		params.RSIEnable = ok
	}
	if params.RSIEnable {
		rsiValues = df.RSI.Values
	}

	buyPoint, sellPoint := 0, 0
	if params.EMAEnable && params.EMAPeriod1 < lenCandles && params.EMAPeriod2 < lenCandles {
		if emaValues1[lenCandles-2] < emaValues2[lenCandles-2] && emaValues1[lenCandles-1] >= emaValues2[lenCandles-1] {
			buyPoint++
		}
		if emaValues1[lenCandles-2] > emaValues2[lenCandles-2] && emaValues1[lenCandles-1] <= emaValues2[lenCandles-1] {
			sellPoint++
		}
	}

	if params.BBandsEnable && params.BBandsN < lenCandles {
		if bbDown[lenCandles-2] > df.Candles[lenCandles-2].Close && bbDown[lenCandles-1] <= df.Candles[lenCandles-1].Close {
			buyPoint++
		}
		if bbUp[lenCandles-2] < df.Candles[lenCandles-2].Close && bbUp[lenCandles-1] >= df.Candles[lenCandles-1].Close {
			sellPoint++
		}
	}

	if params.MACDEnable {
		if outMACD[lenCandles-1] < 0 && outMACDSignal[lenCandles-1] < 0 && outMACD[lenCandles-2] < outMACDSignal[lenCandles-2] && outMACD[lenCandles-1] >= outMACDSignal[lenCandles-1] {
			buyPoint++
		}
		if outMACD[lenCandles-1] > 0 && outMACDSignal[lenCandles-1] > 0 && outMACD[lenCandles-2] > outMACDSignal[lenCandles-2] && outMACD[lenCandles-1] <= outMACDSignal[lenCandles-1] {
			sellPoint++
		}
	}

	if params.IchimokuEnable {
		if chikou[lenCandles-2] < df.Candles[lenCandles-2].High && chikou[lenCandles-1] >= df.Candles[lenCandles-1].High &&
			senkouA[lenCandles-1] < df.Candles[lenCandles-1].Low && senkouB[lenCandles-1] < df.Candles[lenCandles-1].Low &&
			tenkan[lenCandles-1] > kijun[lenCandles-1] {
			buyPoint++
		}
		if chikou[lenCandles-2] > df.Candles[lenCandles-2].Low && chikou[lenCandles-1] <= df.Candles[lenCandles-1].Low &&
			senkouA[lenCandles-1] > df.Candles[lenCandles-1].High && senkouB[lenCandles-1] > df.Candles[lenCandles-1].High &&
			tenkan[lenCandles-1] < kijun[lenCandles-1] {
			sellPoint++
		}
	}

	if params.RSIEnable && rsiValues[lenCandles-2] != 0 && rsiValues[lenCandles-2] != 100 {
		if rsiValues[lenCandles-2] < params.RSISellThread && rsiValues[lenCandles-1] >= params.RSIBuyThread {
			buyPoint++
		}
		if rsiValues[lenCandles-2] > params.RSISellThread && rsiValues[lenCandles-1] <= params.RSISellThread {
			sellPoint++
		}
	}

	if buyPoint > 0 {
		childOrderAcceptanceID, isOrderCompleted := bot.Buy(df.Candles[lenCandles-1], timeFormat)
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] buy order is not completed, id=%s", childOrderAcceptanceID))
		}
	}

	if sellPoint > 0 {
		childOrderAcceptanceID, isOrderCompleted := bot.Sell(df.Candles[lenCandles-1], timeFormat)
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] sell order is not completed, id=%s", childOrderAcceptanceID))
		}
	}

	return nil
}
