package dto

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

type DataFrame struct {
	ProductCode    string         `json:"productCode"`
	Candles        []Candle       `json:"candles"`
	Events         *SignalEvents  `json:"events"`
	SMAs           []SMA          `json:"smas,omitempty"`
	EMAs           []EMA          `json:"emas,omitempty"`
	BBands         *BBands        `json:"bbands,omitempty"`
	IchimokuCloud  *IchimokuCloud `json:"ichimoku,omitempty"`
	RSI            *RSI           `json:"rsi,omitempty"`
	MACD           *MACD          `json:"macd,omitempty"`
	BacktestEvents *SignalEvents  `json:"backtestEvents,omitempty"`
}

func ConvertDataFrame(df *model.DataFrame) DataFrame {
	candles := make([]Candle, 0)
	for _, c := range df.Candles() {
		candles = append(candles, ConvertCandle(c))
	}

	events := ConvertSignalEvents(df.Events())

	smas := make([]SMA, 0)
	for _, s := range df.SMAs() {
		smas = append(smas, ConvertSMA(s))
	}

	emas := make([]EMA, 0)
	for _, e := range df.EMAs() {
		emas = append(emas, ConvertEMA(e))
	}

	bbands := ConvertBBands(df.BBands())

	ichimoku := ConvertIchimokuCloud(df.IchimokuCloud())

	rsi := ConvertRSI(df.RSI())

	macd := ConvertMACD(df.MACD())

	backTestEvents := ConvertSignalEvents(df.BacktestEvents())

	dto := DataFrame{
		ProductCode:    df.ProductCode(),
		Candles:        candles,
		Events:         events,
		SMAs:           smas,
		EMAs:           emas,
		BBands:         bbands,
		IchimokuCloud:  ichimoku,
		RSI:            rsi,
		MACD:           macd,
		BacktestEvents: backTestEvents,
	}

	return dto
}

type Candle struct {
	ProductCode string        `json:"productCode"`
	Duration    time.Duration `json:"duration"`
	Time        time.Time     `json:"time"`
	Open        float64       `json:"open"`
	Close       float64       `json:"close"`
	High        float64       `json:"high"`
	Low         float64       `json:"low"`
	Volume      float64       `json:"volume"`
}

func ConvertCandle(c model.Candle) Candle {
	return Candle{
		ProductCode: c.ProductCode(),
		Duration:    c.Duration(),
		Time:        c.Time().Time(),
		Open:        c.Open(),
		Close:       c.Close(),
		High:        c.High(),
		Low:         c.Low(),
		Volume:      c.Volume(),
	}
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
	Profit  float64       `json:"profit"`
}

func ConvertSignalEvents(s *model.SignalEvents) *SignalEvents {
	if s == nil {
		return nil
	}

	signals := make([]SignalEvent, 0)
	for _, s := range s.Signals() {
		signals = append(signals, ConvertSignalEvent(s))
	}
	return &SignalEvents{
		Signals: signals,
		Profit:  s.Profit(),
	}
}

type SignalEvent struct {
	Time        time.Time       `json:"time"`
	ProductCode string          `json:"productCode"`
	Side        model.OrderSide `json:"side"`
	Price       float64         `json:"price"`
	Size        float64         `json:"size"`
}

func ConvertSignalEvent(s model.SignalEvent) SignalEvent {
	return SignalEvent{
		Time:        s.Time(),
		ProductCode: s.ProductCode(),
		Side:        s.Side(),
		Price:       s.Price(),
		Size:        s.Size(),
	}
}

type SMA struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

func ConvertSMA(sma model.SMA) SMA {
	return SMA{
		Period: sma.Period(),
		Values: sma.Values(),
	}
}

type EMA struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

func ConvertEMA(ema model.EMA) EMA {
	return EMA{
		Period: ema.Period(),
		Values: ema.Values(),
	}
}

type BBands struct {
	N    int       `json:"n,omitempty"`
	K    float64   `json:"k,omitempty"`
	Up   []float64 `json:"up,omitempty"`
	Mid  []float64 `json:"mid,omitempty"`
	Down []float64 `json:"down,omitempty"`
}

func ConvertBBands(bbands *model.BBands) *BBands {
	if bbands == nil {
		return nil
	}

	return &BBands{
		N:    bbands.N(),
		K:    bbands.K(),
		Up:   bbands.Up(),
		Mid:  bbands.Mid(),
		Down: bbands.Down(),
	}
}

type IchimokuCloud struct {
	Tenkan  []float64 `json:"tenkan,omitempty"`
	Kijun   []float64 `json:"kijun,omitempty"`
	SenkouA []float64 `json:"senkoua,omitempty"`
	SenkouB []float64 `json:"senkoub,omitempty"`
	Chikou  []float64 `json:"chikou,omitempty"`
}

func ConvertIchimokuCloud(ic *model.IchimokuCloud) *IchimokuCloud {
	if ic == nil {
		return nil
	}

	return &IchimokuCloud{
		Tenkan:  ic.Tenkan(),
		Kijun:   ic.Kijun(),
		SenkouA: ic.SenkouA(),
		SenkouB: ic.SenkouB(),
		Chikou:  ic.Chikou(),
	}
}

type RSI struct {
	Period int       `json:"period,omitenpty"`
	Values []float64 `json:"values,omitempty"`
}

func ConvertRSI(rsi *model.RSI) *RSI {
	if rsi == nil {
		return nil
	}

	return &RSI{
		Period: rsi.Period(),
		Values: rsi.Values(),
	}
}

type MACD struct {
	FastPeriod   int       `json:"fastPeriod,omitempty"`
	SlowPeriod   int       `json:"slowPeriod,omitempty"`
	SignalPeriod int       `json:"signalPeriod,omitempty"`
	MACD         []float64 `json:"macd,omitempty"`
	MACDSignal   []float64 `json:"macdSignal,omitempty"`
	MACDHist     []float64 `json:"macdHist,omitempty"`
}

func ConvertMACD(macd *model.MACD) *MACD {
	if macd == nil {
		return nil
	}

	return &MACD{
		FastPeriod:   macd.FastPeriod(),
		SlowPeriod:   macd.SlowPeriod(),
		SignalPeriod: macd.SignalPeriod(),
		MACD:         macd.Macd(),
		MACDSignal:   macd.MacdSignal(),
		MACDHist:     macd.MacdHist(),
	}
}

type Balance struct {
	CurrencyCode string  `json:"currencyCode"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func ConvertBalance(balance *model.Balance) *Balance {
	if balance == nil {
		return nil
	}

	return &Balance{
		CurrencyCode: balance.CurrencyCode(),
		Amount:       balance.Amount(),
		Available:    balance.Available(),
	}
}
