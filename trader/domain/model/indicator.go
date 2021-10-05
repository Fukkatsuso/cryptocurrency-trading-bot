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
	if tenkanN := 9; tenkanN > length {
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
