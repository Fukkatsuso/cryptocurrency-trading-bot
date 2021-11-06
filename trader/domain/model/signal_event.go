package model

import (
	"time"
)

type SignalEvent struct {
	time        time.Time
	productCode string
	side        OrderSide
	price       float64
	size        float64
}

func NewSignalEvent(timeTime time.Time, productCode string, side OrderSide, price float64, size float64) *SignalEvent {
	if productCode == "" {
		return nil
	}

	if side == "" {
		return nil
	}

	if price <= 0 {
		return nil
	}

	if size <= 0 {
		return nil
	}

	timeTime = timeTime.In(time.UTC)

	return &SignalEvent{
		time:        timeTime,
		productCode: productCode,
		side:        side,
		price:       price,
		size:        size,
	}
}

func (s *SignalEvent) Time() time.Time {
	return s.time
}

func (s *SignalEvent) ProductCode() string {
	return s.productCode
}

func (s *SignalEvent) Side() OrderSide {
	return s.side
}

func (s *SignalEvent) Price() float64 {
	return s.price
}

func (s *SignalEvent) Size() float64 {
	return s.size
}

type SignalEvents struct {
	signals []SignalEvent
	profit  float64
}

func NewSignalEvents(signals []SignalEvent) *SignalEvents {
	if signals == nil {
		return nil
	}

	return &SignalEvents{
		signals: signals,
		profit:  0,
	}
}

func (s *SignalEvents) LastSignal() *SignalEvent {
	lenSignals := len(s.signals)
	if lenSignals == 0 {
		return nil
	}

	return &s.signals[lenSignals-1]
}

func (s *SignalEvents) Signals() []SignalEvent {
	return s.signals
}

func (s *SignalEvents) Profit() float64 {
	return s.profit
}

func (s *SignalEvents) CanBuyAt(timeTime time.Time) bool {
	lastSignal := s.LastSignal()
	if lastSignal == nil {
		return true
	}

	canBuy := lastSignal.side == OrderSideSell &&
		lastSignal.time.Before(timeTime)
	return canBuy
}

func (s *SignalEvents) CanSellAt(timeTime time.Time) bool {
	lastSignal := s.LastSignal()
	if lastSignal == nil {
		return false
	}

	canSell := lastSignal.side == OrderSideBuy &&
		lastSignal.time.Before(timeTime)
	return canSell
}

func (s *SignalEvents) AddBuySignal(signal SignalEvent) bool {
	if signal.side != OrderSideBuy {
		return false
	}

	if !s.CanBuyAt(signal.time) {
		return false
	}

	s.signals = append(s.signals, signal)
	return true
}

func (s *SignalEvents) AddSellSignal(signal SignalEvent) bool {
	if signal.side != OrderSideSell {
		return false
	}

	if !s.CanSellAt(signal.time) {
		return false
	}

	s.signals = append(s.signals, signal)
	return true
}

// 買って売ってを繰り返した履歴データから，利益を推定
func (s *SignalEvents) EstimateProfit() float64 {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for _, signal := range s.signals {
		if signal.side == OrderSideBuy {
			total -= signal.price * signal.size
			isHolding = true
		}
		if signal.side == OrderSideSell {
			total += signal.price * signal.size
			beforeSell = total
			isHolding = false
		}
	}

	var profit float64
	if isHolding {
		profit = beforeSell
	} else {
		profit = total
	}

	s.profit = profit
	return profit
}

// 損切りすべきか判断する
// 最近の買い注文の後
func (s *SignalEvents) ShouldCutLoss(currentPrice, stopLimitPercent float64) bool {
	if s == nil {
		return false
	}

	lastSignal := s.LastSignal()
	if lastSignal == nil ||
		lastSignal.Side() != OrderSideBuy {
		return false
	}

	stopLimit := lastSignal.Price() * stopLimitPercent
	return currentPrice < stopLimit
}
