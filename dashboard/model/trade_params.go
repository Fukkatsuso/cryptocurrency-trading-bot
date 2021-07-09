package model

import (
	"fmt"
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
}

func GetTradeParams(db DB, tradeParamTableName, productCode string) *TradeParams {
	// 最後に作成されたパラメータを取得
	// productCodeで絞り込み，そのうちcreated_atが最新のレコードを探す
	cmd := fmt.Sprintf(`
            SELECT
                tp.enable,
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
                tp.macd_signal_period
            FROM
                %s AS tp
            WHERE
                tp.product_code = ?
                AND
                tp.created_at = (
                    SELECT
                        MAX(sub_tp.created_at)
                    FROM
                        %s AS sub_tp
                    WHERE
                        sub_tp.product_code = ?
                )`,
		tradeParamTableName, tradeParamTableName)
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
		buyPoint, sellPoint := 0, 0

		if params.EMAEnable && params.EMAPeriod1 <= i && params.EMAPeriod2 <= i {
			emaValue1Prev, emaValue1 := df.EMAs[0].Values[i-1], df.EMAs[0].Values[i-1]
			emaValue2Prev, emaValue2 := df.EMAs[1].Values[i-1], df.EMAs[1].Values[i-1]
			if emaValue1Prev < emaValue2Prev && emaValue1 >= emaValue2 {
				buyPoint++
			}
			if emaValue1Prev > emaValue2Prev && emaValue1 <= emaValue2 {
				sellPoint++
			}
		}

		if params.BBandsEnable && params.BBandsN <= i {
			bbandsUpPrev, bbandsUp := df.BBands.Up[i-1], df.BBands.Up[i]
			bbandsDownPrev, bbandsDown := df.BBands.Down[i-1], df.BBands.Down[i]
			if bbandsDownPrev > df.Candles[i-1].Close && bbandsDown <= df.Candles[i].Close {
				buyPoint++
			}
			if bbandsUpPrev < df.Candles[i-1].Close && bbandsUp >= df.Candles[i].Close {
				sellPoint++
			}
		}

		if params.IchimokuEnable {
			tenkan := df.IchimokuCloud.Tenkan[i]
			kijun := df.IchimokuCloud.Kijun[i]
			senkouA := df.IchimokuCloud.SenkouA[i]
			senkouB := df.IchimokuCloud.SenkouB[i]
			chikouPrev, chikou := df.IchimokuCloud.Chikou[i-1], df.IchimokuCloud.Chikou[i]
			if chikouPrev < df.Candles[i-1].High && chikou >= df.Candles[i].High &&
				senkouA < df.Candles[i].Low && senkouB < df.Candles[i].Low &&
				tenkan > kijun {
				buyPoint++
			}
			if chikouPrev > df.Candles[i-1].Low && chikou <= df.Candles[i].Low &&
				senkouA > df.Candles[i].High && senkouB > df.Candles[i].High &&
				tenkan < kijun {
				sellPoint++
			}
		}

		if params.RSIEnable && df.RSI.Values[i-1] != 0 && df.RSI.Values[i-1] != 100 {
			rsiPrev, rsi := df.RSI.Values[i-1], df.RSI.Values[i]
			if rsiPrev < params.RSISellThread && rsi >= params.RSIBuyThread {
				buyPoint++
			}
			if rsiPrev > params.RSISellThread && rsi <= params.RSISellThread {
				sellPoint++
			}
		}

		if params.MACDEnable {
			macdPrev, macd := df.MACD.MACD[i-1], df.MACD.MACD[i]
			signalPrev, signal := df.MACD.MACDSignal[i-1], df.MACD.MACDSignal[i]
			if macd < 0 && signal < 0 && macdPrev < signalPrev && macd >= signal {
				buyPoint++
			}
			if macd > 0 && signal > 0 && macdPrev > signalPrev && macd <= signal {
				sellPoint++
			}
		}

		if buyPoint > sellPoint {
			events.Buy(params.ProductCode, df.Candles[i].Time, df.Candles[i].Close, params.Size)
		}
		if sellPoint > buyPoint {
			events.Sell(params.ProductCode, df.Candles[i].Time, df.Candles[i].Close, params.Size)
		}
	}
	df.BacktestEvents = events

	df.BacktestEvents.EstimateProfit()
}
