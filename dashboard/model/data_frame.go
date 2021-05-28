package model

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/lib/trading"
	"github.com/markcheno/go-talib"
)

type DataFrame struct {
	ProductCode    string         `json:"productCode"`
	Candles        []Candle       `json:"candles"`
	SMAs           []SMA          `json:"smas,omitempty"`
	EMAs           []EMA          `json:"emas,omitempty"`
	BBands         *BBands        `json:"bbands,omitempty"`
	IchimokuCloud  *IchimokuCloud `json:"ichimoku,omitempty"`
	RSI            *RSI           `json:"rsi,omitempty"`
	MACD           *MACD          `json:"macd,omitempty"`
	BacktestEvents *SignalEvents  `json:"backtestEvents,omitempty"`
}

// 単純移動平均
type SMA struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// 指数平滑移動平均線
type EMA struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// ボリンジャーバンド
type BBands struct {
	N    int       `json:"n,omitempty"`
	K    float64   `json:"k,omitempty"`
	Up   []float64 `json:"up,omitempty"`
	Mid  []float64 `json:"mid,omitempty"`
	Down []float64 `json:"down,omitempty"`
}

// 一目均衡表
type IchimokuCloud struct {
	Tenkan  []float64 `json:"tenkan,omitempty"`
	Kijun   []float64 `json:"kijun,omitempty"`
	SenkouA []float64 `json:"senkoua,omitempty"`
	SenkouB []float64 `json:"senkoub,omitempty"`
	Chikou  []float64 `json:"chikou,omitempty"`
}

// Relative Strength Index
type RSI struct {
	Period int       `json:"period,omitenpty"`
	Values []float64 `json:"values,omitempty"`
}

// Moving Average Convergence/Divergence: 移動平均・収束拡散
type MACD struct {
	FastPeriod   int       `json:"fast_period,omitempty"`
	SlowPeriod   int       `json:"slow_period,omitempty"`
	SignalPeriod int       `json:"signal_period,omitempty"`
	MACD         []float64 `json:"macd,omitempty"`
	MACDSignal   []float64 `json:"macd_signal,omitempty"`
	MACDHist     []float64 `json:"macd_hist,omitempty"`
}

func (df *DataFrame) Times() []time.Time {
	s := make([]time.Time, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Time
	}
	return s
}

func (df *DataFrame) Opens() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Open
	}
	return s
}

func (df *DataFrame) Closes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Close
	}
	return s
}

func (df *DataFrame) Highs() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.High
	}
	return s
}

func (df *DataFrame) Lows() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Low
	}
	return s
}

func (df *DataFrame) Volumes() []float64 {
	s := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		s[i] = candle.Volume
	}
	return s
}

func (df *DataFrame) AddSMA(period int) bool {
	if df.SMAs == nil {
		df.SMAs = make([]SMA, 0)
	}
	if len(df.Candles) > period {
		df.SMAs = append(df.SMAs, SMA{
			Period: period,
			Values: talib.Sma(df.Closes(), period),
		})
		return true
	}
	return false
}

func (df *DataFrame) AddEMA(period int) bool {
	if df.EMAs == nil {
		df.EMAs = make([]EMA, 0)
	}
	if len(df.Candles) > period {
		df.EMAs = append(df.EMAs, EMA{
			Period: period,
			Values: talib.Ema(df.Closes(), period),
		})
		return true
	}
	return false
}

func (df *DataFrame) AddBBands(n int, k float64) bool {
	if n <= len(df.Closes()) {
		up, mid, down := talib.BBands(df.Closes(), n, k, k, 0)
		df.BBands = &BBands{
			N:    n,
			K:    k,
			Up:   up,
			Mid:  mid,
			Down: down,
		}
		return true
	}
	return false
}

func (df *DataFrame) AddIchimoku() bool {
	tenkanN := 9
	if len(df.Closes()) >= tenkanN {
		tenkan, kijun, senkouA, senkouB, chikou := trading.IchimokuCloud(df.Closes())
		df.IchimokuCloud = &IchimokuCloud{
			Tenkan:  tenkan,
			Kijun:   kijun,
			SenkouA: senkouA,
			SenkouB: senkouB,
			Chikou:  chikou,
		}
		return true
	}
	return false
}

func (df *DataFrame) AddRSI(period int) bool {
	if len(df.Candles) > period {
		values := talib.Rsi(df.Closes(), period)
		df.RSI = &RSI{
			Period: period,
			Values: values,
		}
		return true
	}
	return false
}

func (df *DataFrame) AddMACD(inFastPeriod, inSlowPeriod, inSignalPeriod int) bool {
	if len(df.Candles) > 1 {
		outMACD, outMACDSignal, outMACDHist := talib.Macd(df.Closes(), inFastPeriod, inSlowPeriod, inSignalPeriod)
		df.MACD = &MACD{
			FastPeriod:   inFastPeriod,
			SlowPeriod:   inSlowPeriod,
			SignalPeriod: inSignalPeriod,
			MACD:         outMACD,
			MACDSignal:   outMACDSignal,
			MACDHist:     outMACDHist,
		}
		return true
	}
	return false
}

type TradeParams struct {
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
