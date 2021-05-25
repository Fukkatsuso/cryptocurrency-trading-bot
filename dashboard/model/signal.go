package model

import (
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
	Profit  float64       `json:"profit"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func (s *SignalEvents) Buy(productCode string, time time.Time, price, size float64) bool {
	signalEvent := SignalEvent{
		ProductCode: productCode,
		Time:        time,
		Side:        "BUY",
		Price:       price,
		Size:        size,
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

func (s *SignalEvents) Sell(productCode string, time time.Time, price, size float64) bool {
	signalEvent := SignalEvent{
		ProductCode: productCode,
		Time:        time,
		Side:        "SELL",
		Price:       price,
		Size:        size,
	}
	s.Signals = append(s.Signals, signalEvent)
	return true
}

// 買って売ってを繰り返した履歴データから，利益が出るか検証
func (s *SignalEvents) EstimateProfit() {
	total := 0.0
	beforeSell := 0.0
	isHolding := false
	for i, signalEvent := range s.Signals {
		if i == 0 && signalEvent.Side == "SELL" {
			continue
		}
		if signalEvent.Side == "BUY" {
			total -= signalEvent.Price * signalEvent.Size
			isHolding = true
		}
		if signalEvent.Side == "SELL" {
			total += signalEvent.Price * signalEvent.Size
			isHolding = false
			beforeSell = total
		}
	}
	if isHolding {
		s.Profit = beforeSell
	}
	s.Profit = total
}
