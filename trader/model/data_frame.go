package model

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/trading"
	"github.com/markcheno/go-talib"
)

type DataFrame struct {
	ProductCode   string         `json:"productCode"`
	Candles       []Candle       `json:"candles"`
	SMAs          []SMA          `json:"smas,omitempty"`
	EMAs          []EMA          `json:"emas,omitempty"`
	BBands        *BBands        `json:"bbands,omitempty"`
	IchimokuCloud *IchimokuCloud `json:"ichimoku,omitempty"`
	RSI           *RSI           `json:"rsi,omitempty"`
	MACD          *MACD          `json:"macd,omitempty"`
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
