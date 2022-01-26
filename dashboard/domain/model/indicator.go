package model

import "github.com/markcheno/go-talib"

// 単純移動平均
type SMA struct {
	period int
	values []float64
}

func NewSMA(inReal []float64, period int) *SMA {
	if period <= 0 || len(inReal) <= period {
		return nil
	}

	values := talib.Sma(inReal, period)
	if values == nil {
		return nil
	}

	return &SMA{
		period: period,
		values: values,
	}
}

func (sma *SMA) Period() int {
	return sma.period
}

func (sma *SMA) Values() []float64 {
	return sma.values
}

// 指数平滑移動平均線
type EMA struct {
	period int
	values []float64
}

func NewEMA(inReal []float64, period int) *EMA {
	if period <= 0 || len(inReal) <= period {
		return nil
	}

	values := talib.Ema(inReal, period)
	if values == nil {
		return nil
	}

	return &EMA{
		period: period,
		values: values,
	}
}

func (ema *EMA) Period() int {
	return ema.period
}

func (ema *EMA) Values() []float64 {
	return ema.values
}

// ボリンジャーバンド
type BBands struct {
	n    int
	k    float64
	up   []float64
	mid  []float64
	down []float64
}

func NewBBands(inReal []float64, n int, k float64) *BBands {
	if n <= 0 || len(inReal) < n {
		return nil
	}

	if k <= 0 {
		return nil
	}

	up, mid, down := talib.BBands(inReal, n, k, k, 0)

	return &BBands{
		n:    n,
		k:    k,
		up:   up,
		mid:  mid,
		down: down,
	}
}

func (bbands *BBands) N() int {
	return bbands.n
}

func (bbands *BBands) K() float64 {
	return bbands.k
}

func (bbands *BBands) Up() []float64 {
	return bbands.up
}

func (bbands *BBands) Mid() []float64 {
	return bbands.mid
}

func (bbands *BBands) Down() []float64 {
	return bbands.down
}

// 一目均衡表
type IchimokuCloud struct {
	tenkan  []float64
	kijun   []float64
	senkouA []float64
	senkouB []float64
	chikou  []float64
}

func minMax(inReal []float64) (float64, float64) {
	min := inReal[0]
	max := inReal[0]
	for _, price := range inReal {
		if min > price {
			min = price
		}
		if max < price {
			max = price
		}
	}
	return min, max
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func NewIchimokuCloud(inReal []float64) *IchimokuCloud {
	length := len(inReal)
	if length < 52 {
		return nil
	}

	tenkan := make([]float64, min(9, length))
	kijun := make([]float64, min(26, length))
	senkouA := make([]float64, min(26, length))
	senkouB := make([]float64, min(52, length))
	chikou := make([]float64, min(26, length))

	for i := range inReal {
		if i >= 9 {
			min, max := minMax(inReal[i-9 : i])
			tenkan = append(tenkan, (min+max)/2)
		}
		if i >= 26 {
			min, max := minMax(inReal[i-26 : i])
			kijun = append(kijun, (min+max)/2)
			senkouA = append(senkouA, (tenkan[i]+kijun[i])/2)
			chikou = append(chikou, inReal[i-26])
		}
		if i >= 52 {
			min, max := minMax(inReal[i-52 : i])
			senkouB = append(senkouB, (min+max)/2)
		}
	}

	return &IchimokuCloud{
		tenkan:  tenkan,
		kijun:   kijun,
		senkouA: senkouA,
		senkouB: senkouB,
		chikou:  chikou,
	}
}

func (ichimokuCloud *IchimokuCloud) Tenkan() []float64 {
	return ichimokuCloud.tenkan
}

func (ichimokuCloud *IchimokuCloud) Kijun() []float64 {
	return ichimokuCloud.kijun
}

func (ichimokuCloud *IchimokuCloud) SenkouA() []float64 {
	return ichimokuCloud.senkouA
}

func (ichimokuCloud *IchimokuCloud) SenkouB() []float64 {
	return ichimokuCloud.senkouB
}

func (ichimokuCloud *IchimokuCloud) Chikou() []float64 {
	return ichimokuCloud.chikou
}

// Relative Strength Index
type RSI struct {
	period int
	values []float64
}

func NewRSI(inReal []float64, period int) *RSI {
	if period <= 0 || len(inReal) <= period {
		return nil
	}

	values := talib.Rsi(inReal, period)

	return &RSI{
		period: period,
		values: values,
	}
}

func (rsi *RSI) Period() int {
	return rsi.period
}

func (rsi *RSI) Values() []float64 {
	return rsi.values
}

// Moving Average Convergence/Divergence: 移動平均・収束拡散
type MACD struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
	macd         []float64
	macdSignal   []float64
	macdHist     []float64
}

func NewMACD(inReal []float64, inFastPeriod, inSlowPeriod, inSignalPeriod int) *MACD {
	if len(inReal) <= 0 {
		return nil
	}

	if inFastPeriod <= 0 || len(inReal) <= inFastPeriod {
		return nil
	}

	if inSlowPeriod <= 0 || len(inReal) <= inSlowPeriod {
		return nil
	}

	if inSignalPeriod <= 0 || len(inReal) <= inSignalPeriod {
		return nil
	}

	outMACD, outMACDSignal, outMACDHist := talib.Macd(inReal, inFastPeriod, inSlowPeriod, inSignalPeriod)

	return &MACD{
		fastPeriod:   inFastPeriod,
		slowPeriod:   inSlowPeriod,
		signalPeriod: inSignalPeriod,
		macd:         outMACD,
		macdSignal:   outMACDSignal,
		macdHist:     outMACDHist,
	}
}

func (macd *MACD) FastPeriod() int {
	return macd.fastPeriod
}

func (macd *MACD) SlowPeriod() int {
	return macd.slowPeriod
}

func (macd *MACD) SignalPeriod() int {
	return macd.signalPeriod
}

func (macd *MACD) Macd() []float64 {
	return macd.macd
}

func (macd *MACD) MacdSignal() []float64 {
	return macd.macdSignal
}

func (macd *MACD) MacdHist() []float64 {
	return macd.macdHist
}

// 平均足
type AverageCandle struct {
	opens  []float64
	closes []float64
	highs  []float64
	lows   []float64
}

func NewAverageCandle(candles []Candle) *AverageCandle {
	lenCandle := len(candles)

	opens := make([]float64, lenCandle)
	closes := make([]float64, lenCandle)
	highs := make([]float64, lenCandle)
	lows := make([]float64, lenCandle)
	for i, candle := range candles {
		// open
		if i == 0 {
			opens[i] = candle.Open()
		} else {
			opens[i] = (candles[i-1].Open() + candles[i-1].Close()) / 2.0
		}
		// close
		closes[i] = (candle.Open() + candle.Close() + candle.High() + candle.Low()) / 4.0
		// high
		highs[i] = candle.High()
		// low
		lows[i] = candle.Low()
	}

	return &AverageCandle{
		opens:  opens,
		closes: closes,
		highs:  highs,
		lows:   lows,
	}
}

func (ac *AverageCandle) Opens() []float64 {
	return ac.opens
}

func (ac *AverageCandle) Closes() []float64 {
	return ac.closes
}

func (ac *AverageCandle) High() []float64 {
	return ac.highs
}

func (ac *AverageCandle) Lows() []float64 {
	return ac.lows
}
