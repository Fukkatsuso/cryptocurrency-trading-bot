package model

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/trading"
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
	if period < len(df.Candles) {
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
	if period < len(df.Candles) {
		df.EMAs = append(df.EMAs, EMA{
			Period: period,
			Values: talib.Ema(df.Closes(), period),
		})
		return true
	}
	return false
}

func (df *DataFrame) AddBBands(n int, k float64) bool {
	if n <= len(df.Candles) {
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
	if tenkanN <= len(df.Candles) {
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
	if period < len(df.Candles) {
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
	if len(df.Candles) > 1 &&
		inFastPeriod < len(df.Candles) && inSlowPeriod < len(df.Candles) && inSignalPeriod < len(df.Candles) {
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

func (df *DataFrame) BackTestEMA(period1, period2 int, size float64) *SignalEvents {
	lenCandles := len(df.Candles)
	if lenCandles <= period1 || lenCandles <= period2 {
		return nil
	}
	signalEvents := NewSignalEvents()
	emaValue1 := talib.Ema(df.Closes(), period1)
	emaValue2 := talib.Ema(df.Closes(), period2)

	for i := 1; i < lenCandles; i++ {
		if i < period1 || i < period2 {
			continue
		}
		// ゴールデンクロス
		if emaValue1[i-1] < emaValue2[i-1] && emaValue1[i] >= emaValue2[i] {
			signalEvents.Buy(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
		// デッドクロス
		if emaValue1[i-1] > emaValue2[i-1] && emaValue1[i] <= emaValue2[i] {
			signalEvents.Sell(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
	}
	return signalEvents
}

func (df *DataFrame) OptimizeEMA(period1, period2 int, size float64) (float64, int, int) {
	performance := float64(0)
	bestPeriod1 := period1
	bestPeriod2 := period2

	for period1 = 5; period1 < 11; period1++ {
		for period2 = 12; period2 < 20; period2++ {
			signalEvents := df.BackTestEMA(period1, period2, size)
			if signalEvents == nil {
				continue
			}
			profit := signalEvents.EstimateProfit()
			if performance < profit {
				performance = profit
				bestPeriod1 = period1
				bestPeriod2 = period2
			}
		}
	}
	return performance, bestPeriod1, bestPeriod2
}

func (df *DataFrame) BackTestBBands(n int, k float64, size float64) *SignalEvents {
	lenCandles := len(df.Candles)
	if lenCandles <= n {
		return nil
	}

	signalEvents := &SignalEvents{}
	bbUp, _, bbDown := talib.BBands(df.Closes(), n, k, k, 0)
	for i := 1; i < lenCandles; i++ {
		if i < n {
			continue
		}
		if bbDown[i-1] > df.Candles[i-1].Close && bbDown[i] <= df.Candles[i].Close {
			signalEvents.Buy(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
		if bbUp[i-1] < df.Candles[i-1].Close && bbUp[i] >= df.Candles[i].Close {
			signalEvents.Sell(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
	}
	return signalEvents
}

func (df *DataFrame) OptimizeBBands(n int, k float64, size float64) (float64, int, float64) {
	performance := float64(0)
	bestN := n
	bestK := k

	for n := 10; n <= 30; n++ {
		for k := 1.8; k <= 2.2; k += 0.1 {
			signalEvents := df.BackTestBBands(n, k, size)
			if signalEvents == nil {
				continue
			}
			profit := signalEvents.EstimateProfit()
			if performance < profit {
				performance = profit
				bestN = n
				bestK = k
			}
		}
	}
	return performance, bestN, bestK
}

func (df *DataFrame) BackTestIchimoku(size float64) *SignalEvents {
	lenCandles := len(df.Candles)
	if lenCandles <= 52 {
		return nil
	}

	var signalEvents SignalEvents
	tenkan, kijun, senkouA, senkouB, chikou := trading.IchimokuCloud(df.Closes())
	for i := 1; i < lenCandles; i++ {
		// 三役好転
		if chikou[i-1] < df.Candles[i-1].High && chikou[i] >= df.Candles[i].High &&
			senkouA[i] < df.Candles[i].Low && senkouB[i] < df.Candles[i].Low &&
			tenkan[i] > kijun[i] {
			signalEvents.Buy(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
		// 三役逆転
		if chikou[i-1] > df.Candles[i-1].Low && chikou[i] <= df.Candles[i].Low &&
			senkouA[i] > df.Candles[i].High && senkouB[i] > df.Candles[i].High &&
			tenkan[i] < kijun[i] {
			signalEvents.Sell(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
	}
	return &signalEvents
}

func (df *DataFrame) OptimizeIchimoku(size float64) float64 {
	signalEvents := df.BackTestIchimoku(size)
	if signalEvents == nil {
		return 0
	}
	performance := signalEvents.EstimateProfit()
	return performance
}

func (df *DataFrame) BackTestRSI(period int, buyThread, sellThread float64, size float64) *SignalEvents {
	lenCandles := len(df.Candles)
	if lenCandles <= period {
		return nil
	}

	signalEvents := NewSignalEvents()
	values := talib.Rsi(df.Closes(), period)
	for i := 1; i < lenCandles; i++ {
		if values[i-1] == 0 || values[i-1] == 100 {
			continue
		}
		if values[i-1] < buyThread && values[i] >= buyThread {
			signalEvents.Buy(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
		if values[i-1] > sellThread && values[i] <= sellThread {
			signalEvents.Sell(df.ProductCode, df.Candles[i].Time, df.Candles[i].Close, size)
		}
	}
	return signalEvents
}

func (df *DataFrame) OptimizeRSI(period int, buyThread, sellThread float64, size float64) (float64, int, float64, float64) {
	performance := float64(0)
	bestPeriod := period
	bestBuyThread, bestSellThread := buyThread, sellThread

	for period := 3; period < 30; period++ {
		for buyThread := float64(20); buyThread <= 40; buyThread++ {
			for sellThread := float64(60); sellThread <= 80; sellThread++ {
				signalEvents := df.BackTestRSI(period, buyThread, sellThread, size)
				if signalEvents == nil {
					continue
				}
				profit := signalEvents.EstimateProfit()
				if performance < profit {
					performance = profit
					bestPeriod = period
					bestBuyThread = buyThread
					bestSellThread = sellThread
				}
			}
		}
	}
	return performance, bestPeriod, bestBuyThread, bestSellThread
}
