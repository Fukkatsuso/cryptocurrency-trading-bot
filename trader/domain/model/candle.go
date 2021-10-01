package model

import "time"

type Candle struct {
	productCode string
	duration    time.Duration
	time        CandleTime
	open        float64
	close       float64
	high        float64
	low         float64
	volume      float64
}

func NewCandle(productCode string, duration time.Duration, candleTime CandleTime, open, close, high, low, volume float64) *Candle {
	if productCode == "" {
		return nil
	}

	if duration <= 0 {
		return nil
	}

	if open <= 0 {
		return nil
	}

	if close <= 0 {
		return nil
	}

	if high <= 0 {
		return nil
	}

	if low <= 0 {
		return nil
	}

	if volume <= 0 {
		return nil
	}

	if high < low {
		return nil
	}

	if high < open || high < close {
		return nil
	}

	if low > open || low > close {
		return nil
	}

	return &Candle{
		productCode: productCode,
		duration:    duration,
		time:        candleTime,
		open:        open,
		close:       close,
		high:        high,
		low:         low,
		volume:      volume,
	}
}

func (candle *Candle) ProductCode() string {
	return candle.productCode
}

func (candle *Candle) Duration() time.Duration {
	return candle.duration
}

func (candle *Candle) Time() CandleTime {
	return candle.time
}

func (candle *Candle) Open() float64 {
	return candle.open
}

func (candle *Candle) Close() float64 {
	return candle.close
}

func (candle *Candle) High() float64 {
	return candle.high
}

func (candle *Candle) Low() float64 {
	return candle.low
}

func (candle *Candle) Volume() float64 {
	return candle.volume
}
