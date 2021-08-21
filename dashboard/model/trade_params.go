package model

import (
	"fmt"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/bitflyer"
)

type TradeParams struct {
	TradeEnable      bool
	ProductCode      string
	Size             float64
	SMAEnable        bool
	SMAPeriod1       int
	SMAPeriod2       int
	SMAPeriod3       int
	EMAEnable        bool
	EMAPeriod1       int
	EMAPeriod2       int
	EMAPeriod3       int
	BBandsEnable     bool
	BBandsN          int
	BBandsK          float64
	IchimokuEnable   bool
	RSIEnable        bool
	RSIPeriod        int
	RSIBuyThread     float64
	RSISellThread    float64
	MACDEnable       bool
	MACDFastPeriod   int
	MACDSlowPeriod   int
	MACDSignalPeriod int
	StopLimitPercent float64
}

func (tradeParams *TradeParams) Create(db DB) error {
	cmd := fmt.Sprintf(`
        INSERT INTO trade_params (
            trade_enable,
            product_code,
            size,
            sma_enable,
            sma_period1,
            sma_period2,
            sma_period3,
            ema_enable,
            ema_period1,
            ema_period2,
            ema_period3,
            bbands_enable,
            bbands_n,
            bbands_k,
            ichimoku_enable,
            rsi_enable,
            rsi_period,
            rsi_buy_thread,
            rsi_sell_thread,
            macd_enable,
            macd_fast_period,
            macd_slow_period,
            macd_signal_period,
            stop_limit_percent
        ) VALUES (
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?,
            ?
        )`,
	)

	_, err := db.Exec(cmd,
		tradeParams.TradeEnable,
		tradeParams.ProductCode,
		tradeParams.Size,
		tradeParams.SMAEnable,
		tradeParams.SMAPeriod1,
		tradeParams.SMAPeriod2,
		tradeParams.SMAPeriod3,
		tradeParams.EMAEnable,
		tradeParams.EMAPeriod1,
		tradeParams.EMAPeriod2,
		tradeParams.EMAPeriod3,
		tradeParams.BBandsEnable,
		tradeParams.BBandsN,
		tradeParams.BBandsK,
		tradeParams.IchimokuEnable,
		tradeParams.RSIEnable,
		tradeParams.RSIPeriod,
		tradeParams.RSIBuyThread,
		tradeParams.RSISellThread,
		tradeParams.MACDEnable,
		tradeParams.MACDFastPeriod,
		tradeParams.MACDSlowPeriod,
		tradeParams.MACDSignalPeriod,
		tradeParams.StopLimitPercent,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetTradeParams(db DB, productCode string) *TradeParams {
	// 最後に作成されたパラメータを取得
	// productCodeで絞り込み，そのうちcreated_atが最新のレコードを探す
	cmd := fmt.Sprintf(`
            SELECT
                tp.trade_enable,
                tp.size,
                tp.sma_enable,
                tp.sma_period1,
                tp.sma_period2,
                tp.sma_period3,
                tp.ema_enable,
                tp.ema_period1,
                tp.ema_period2,
                tp.ema_period3,
                tp.bbands_enable,
                tp.bbands_n,
                tp.bbands_k,
                tp.ichimoku_enable,
                tp.rsi_enable,
                tp.rsi_period,
                tp.rsi_buy_thread,
                tp.rsi_sell_thread,
                tp.macd_enable,
                tp.macd_fast_period,
                tp.macd_slow_period,
                tp.macd_signal_period,
                tp.stop_limit_percent
            FROM
                trade_params AS tp
            WHERE
                tp.product_code = ?
                AND
                tp.created_at = (
                    SELECT
                        MAX(sub_tp.created_at)
                    FROM
                        trade_params AS sub_tp
                    WHERE
                        sub_tp.product_code = ?
                )`,
	)
	row := db.QueryRow(cmd, productCode, productCode)

	var tradeParams TradeParams
	err := row.Scan(
		&tradeParams.TradeEnable,
		&tradeParams.Size,
		&tradeParams.SMAEnable,
		&tradeParams.SMAPeriod1,
		&tradeParams.SMAPeriod2,
		&tradeParams.SMAPeriod3,
		&tradeParams.EMAEnable,
		&tradeParams.EMAPeriod1,
		&tradeParams.EMAPeriod2,
		&tradeParams.EMAPeriod3,
		&tradeParams.BBandsEnable,
		&tradeParams.BBandsN,
		&tradeParams.BBandsK,
		&tradeParams.IchimokuEnable,
		&tradeParams.RSIEnable,
		&tradeParams.RSIPeriod,
		&tradeParams.RSIBuyThread,
		&tradeParams.RSISellThread,
		&tradeParams.MACDEnable,
		&tradeParams.MACDFastPeriod,
		&tradeParams.MACDSlowPeriod,
		&tradeParams.MACDSignalPeriod,
		&tradeParams.StopLimitPercent,
	)
	if err != nil {
		return nil
	}

	tradeParams.ProductCode = productCode
	return &tradeParams
}

func (df *DataFrame) BackTest(params *TradeParams) {
	if params == nil {
		return
	}

	events := NewSignalEvents()
	for i := 1; i < len(df.Candles); i++ {
		buyPoint, sellPoint := df.Analyze(i, params)

		if buyPoint > 0 {
			events.Buy(params.ProductCode, df.Candles[i].Time, df.Candles[i].Close, params.Size)
		}

		currentPrice := df.Candles[i].Close
		if sellPoint > 0 || ShouldCutLoss(events, currentPrice, params.StopLimitPercent) {
			events.Sell(params.ProductCode, df.Candles[i].Time, df.Candles[i].Close, params.Size)
		}
	}
	events.Profit = events.EstimateProfit()

	df.BacktestEvents = events
}

// 各指標の時点"at"で分析する
// buyPoint, sellPointを返す
func (df *DataFrame) Analyze(at int, params *TradeParams) (int, int) {
	buyPoint, sellPoint := 0, 0

	if params.EMAEnable &&
		params.EMAPeriod1 <= at && params.EMAPeriod2 <= at &&
		0 < at && at < len(df.EMAs[0].Values) && at < len(df.EMAs[1].Values) {
		emaValue1Prev, emaValue1 := df.EMAs[0].Values[at-1], df.EMAs[0].Values[at]
		emaValue2Prev, emaValue2 := df.EMAs[1].Values[at-1], df.EMAs[1].Values[at]
		if emaValue1Prev < emaValue2Prev && emaValue1 >= emaValue2 {
			buyPoint++
		}
		if emaValue1Prev > emaValue2Prev && emaValue1 <= emaValue2 {
			sellPoint++
		}
	}

	if params.BBandsEnable &&
		params.BBandsN <= at &&
		0 < at && at < len(df.BBands.Up) && at < len(df.BBands.Down) && at < len(df.Candles) {
		bbandsUpPrev, bbandsUp := df.BBands.Up[at-1], df.BBands.Up[at]
		bbandsDownPrev, bbandsDown := df.BBands.Down[at-1], df.BBands.Down[at]
		if bbandsDownPrev > df.Candles[at-1].Close && bbandsDown <= df.Candles[at].Close {
			buyPoint++
		}
		if bbandsUpPrev < df.Candles[at-1].Close && bbandsUp >= df.Candles[at].Close {
			sellPoint++
		}
	}

	if params.IchimokuEnable &&
		0 < at && at < len(df.IchimokuCloud.Tenkan) && at < len(df.IchimokuCloud.Kijun) && at < len(df.IchimokuCloud.SenkouA) && at < len(df.IchimokuCloud.SenkouB) && at < len(df.IchimokuCloud.Chikou) && at < len(df.Candles) {
		tenkan := df.IchimokuCloud.Tenkan[at]
		kijun := df.IchimokuCloud.Kijun[at]
		senkouA := df.IchimokuCloud.SenkouA[at]
		senkouB := df.IchimokuCloud.SenkouB[at]
		chikouPrev, chikou := df.IchimokuCloud.Chikou[at-1], df.IchimokuCloud.Chikou[at]
		if chikouPrev < df.Candles[at-1].High && chikou >= df.Candles[at].High &&
			senkouA < df.Candles[at].Low && senkouB < df.Candles[at].Low &&
			tenkan > kijun {
			buyPoint++
		}
		if chikouPrev > df.Candles[at-1].Low && chikou <= df.Candles[at].Low &&
			senkouA > df.Candles[at].High && senkouB > df.Candles[at].High &&
			tenkan < kijun {
			sellPoint++
		}
	}

	if params.RSIEnable &&
		0 < at && at < len(df.RSI.Values) &&
		df.RSI.Values[at-1] != 0 && df.RSI.Values[at-1] != 100 {
		rsiPrev, rsi := df.RSI.Values[at-1], df.RSI.Values[at]
		if rsiPrev < params.RSIBuyThread && rsi >= params.RSIBuyThread {
			buyPoint++
		}
		if rsiPrev > params.RSISellThread && rsi <= params.RSISellThread {
			sellPoint++
		}
	}

	if params.MACDEnable &&
		0 < at && at < len(df.MACD.MACD) && at < len(df.MACD.MACDSignal) {
		macdPrev, macd := df.MACD.MACD[at-1], df.MACD.MACD[at]
		signalPrev, signal := df.MACD.MACDSignal[at-1], df.MACD.MACDSignal[at]
		if macd < 0 && signal < 0 && macdPrev < signalPrev && macd >= signal {
			buyPoint++
		}
		if macd > 0 && signal > 0 && macdPrev > signalPrev && macd <= signal {
			sellPoint++
		}
	}

	return buyPoint, sellPoint
}

// 損切りすべきか判断する
// 最近の買い注文の後，
func ShouldCutLoss(events *SignalEvents, currentPrice, stopLimitPercent float64) bool {
	if events == nil {
		return false
	}

	signals := events.Signals
	if len(signals) == 0 {
		return false
	}

	lastSignal := signals[len(signals)-1]
	if lastSignal.Side != string(bitflyer.OrderSideBuy) {
		return false
	}

	stopLimit := lastSignal.Price * stopLimitPercent
	return currentPrice < stopLimit
}

func (df *DataFrame) OptimizeTradeParams(params *TradeParams) *TradeParams {
	_, emaPeriod1, emaPeriod2 := df.OptimizeEMA(params.EMAPeriod1, params.EMAPeriod2, params.Size)
	_, bbandsN, bbandsK := df.OptimizeBBands(params.BBandsN, params.BBandsK, params.Size)
	_, rsiPeriod, rsiBuyThread, rsiSellThread := df.OptimizeRSI(params.RSIPeriod, params.RSIBuyThread, params.RSISellThread, params.Size)
	_, macdFastPeriod, macdSlowPeriod, macdSignalPeriod := df.OptimizeMACD(params.MACDFastPeriod, params.MACDSlowPeriod, params.MACDSignalPeriod, params.Size)

	newParams := &TradeParams{
		TradeEnable:      params.TradeEnable,
		ProductCode:      params.ProductCode,
		Size:             params.Size,
		SMAEnable:        params.SMAEnable,
		SMAPeriod1:       params.SMAPeriod1,
		SMAPeriod2:       params.SMAPeriod2,
		SMAPeriod3:       params.SMAPeriod3,
		EMAEnable:        params.EMAEnable,
		EMAPeriod1:       emaPeriod1,
		EMAPeriod2:       emaPeriod2,
		EMAPeriod3:       params.EMAPeriod3,
		BBandsEnable:     params.BBandsEnable,
		BBandsN:          bbandsN,
		BBandsK:          bbandsK,
		IchimokuEnable:   params.IchimokuEnable,
		RSIEnable:        params.RSIEnable,
		RSIPeriod:        rsiPeriod,
		RSIBuyThread:     rsiBuyThread,
		RSISellThread:    rsiSellThread,
		MACDEnable:       params.MACDEnable,
		MACDFastPeriod:   macdFastPeriod,
		MACDSlowPeriod:   macdSlowPeriod,
		MACDSignalPeriod: macdSignalPeriod,
		StopLimitPercent: params.StopLimitPercent,
	}
	return newParams
}
