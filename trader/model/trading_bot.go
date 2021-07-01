package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/trading"
	"github.com/markcheno/go-talib"
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

func (bot *TradingBot) Trade(db *sql.DB, candleTableName, timeFormat string) error {
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
		emaValues1 = talib.Ema(df.Closes(), params.EMAPeriod1)
		emaValues2 = talib.Ema(df.Closes(), params.EMAPeriod2)
	}

	var bbUp []float64
	var bbDown []float64
	if params.BBandsEnable {
		bbUp, _, bbDown = talib.BBands(df.Closes(), params.BBandsN, params.BBandsK, params.BBandsK, 0)
	}

	var tenkan, kijun, senkouA, senkouB, chikou []float64
	if params.IchimokuEnable {
		tenkan, kijun, senkouA, senkouB, chikou = trading.IchimokuCloud(df.Closes())
	}

	var outMACD, outMACDSignal []float64
	if params.MACDEnable {
		outMACD, outMACDSignal, _ = talib.Macd(df.Closes(), params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod)
	}

	var rsiValues []float64
	if params.RSIEnable {
		rsiValues = talib.Rsi(df.Closes(), params.RSIPeriod)
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
		childOrderAcceptanceID, isOrderCompleted := bot.Buy(df.Candles[lenCandles-1])
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] buy order is not completed, id=%s", childOrderAcceptanceID))
		}
	}

	if sellPoint > 0 {
		childOrderAcceptanceID, isOrderCompleted := bot.Sell(df.Candles[lenCandles-1])
		if !isOrderCompleted {
			return errors.New(fmt.Sprintf("[Trade] sell order is not completed, id=%s", childOrderAcceptanceID))
		}
	}

	return nil
}
