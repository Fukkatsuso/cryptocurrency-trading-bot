package model

type DataFrame struct {
	productCode    string
	candles        []Candle
	events         *SignalEvents
	smas           []SMA
	emas           []EMA
	bbands         *BBands
	ichimokuCloud  *IchimokuCloud
	rsi            *RSI
	macd           *MACD
	backtestEvents *SignalEvents
}

func NewDataFrame(productCode string, candles []Candle, events *SignalEvents) *DataFrame {
	if productCode == "" {
		return nil
	}

	return &DataFrame{
		productCode: productCode,
		candles:     candles,
		events:      events,
	}
}

func (df *DataFrame) ProductCode() string {
	return df.productCode
}

func (df *DataFrame) Candles() []Candle {
	return df.candles
}

func (df *DataFrame) Times() []CandleTime {
	s := make([]CandleTime, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.Time()
	}
	return s
}

func (df *DataFrame) Opens() []float64 {
	s := make([]float64, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.Open()
	}
	return s
}

func (df *DataFrame) Closes() []float64 {
	s := make([]float64, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.Close()
	}
	return s
}

func (df *DataFrame) Highs() []float64 {
	s := make([]float64, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.High()
	}
	return s
}

func (df *DataFrame) Lows() []float64 {
	s := make([]float64, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.Low()
	}
	return s
}

func (df *DataFrame) Volumes() []float64 {
	s := make([]float64, len(df.candles))
	for i, candle := range df.candles {
		s[i] = candle.Volume()
	}
	return s
}

func (df *DataFrame) Events() *SignalEvents {
	return df.events
}

func (df *DataFrame) SMAs() []SMA {
	return df.smas
}

func (df *DataFrame) EMAs() []EMA {
	return df.emas
}

func (df *DataFrame) BBands() *BBands {
	return df.bbands
}

func (df *DataFrame) IchimokuCloud() *IchimokuCloud {
	return df.ichimokuCloud
}

func (df *DataFrame) RSI() *RSI {
	return df.rsi
}

func (df *DataFrame) MACD() *MACD {
	return df.macd
}

func (df *DataFrame) BacktestEvents() *SignalEvents {
	return df.backtestEvents
}

func (df *DataFrame) AddSMA(period int) bool {
	if df.smas == nil {
		df.smas = make([]SMA, 0)
	}

	sma := NewSMA(df.Closes(), period)
	if sma == nil {
		return false
	}

	df.smas = append(df.smas, *sma)
	return true
}

func (df *DataFrame) AddEMA(period int) bool {
	if df.emas == nil {
		df.emas = make([]EMA, 0)
	}

	ema := NewEMA(df.Closes(), period)
	if ema == nil {
		return false
	}

	df.emas = append(df.emas, *ema)
	return true
}

func (df *DataFrame) AddBBands(n int, k float64) bool {
	bbands := NewBBands(df.Closes(), n, k)
	if bbands == nil {
		return false
	}

	df.bbands = bbands
	return true
}

func (df *DataFrame) AddIchimoku() bool {
	ichimoku := NewIchimokuCloud(df.Closes())
	if ichimoku == nil {
		return false
	}

	df.ichimokuCloud = ichimoku
	return true
}

func (df *DataFrame) AddRSI(period int) bool {
	rsi := NewRSI(df.Closes(), period)
	if rsi == nil {
		return false
	}

	df.rsi = rsi
	return true
}

func (df *DataFrame) AddMACD(inFastPeriod, inSlowPeriod, inSignalPeriod int) bool {
	macd := NewMACD(df.Closes(), inFastPeriod, inSlowPeriod, inSignalPeriod)
	if macd == nil {
		return false
	}

	df.macd = macd
	return true
}

func (df *DataFrame) AddBacktestEvents(events *SignalEvents) {
	df.backtestEvents = events
}

// レンジ相場かどうか判定する
// 一定期間RSIが40-60の間を推移していれば，レンジ相場
func (df *DataFrame) IsBoxedRange(period, at int) bool {
	values := df.rsi.values
	for i := at; i >= 0 && at-i < period; i-- {
		if values[i] <= 40 || 60 <= values[i] {
			return false
		}
	}
	return true
}
