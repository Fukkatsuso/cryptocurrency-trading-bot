package model

import (
	"fmt"
	"strings"
	"time"
)

type SignalEvent struct {
	Time        time.Time `json:"time"`
	ProductCode string    `json:"product_code"`
	Side        string    `json:"side"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
}

func (s *SignalEvent) Save(db DB, timeFormat string) bool {
	cmd := "INSERT INTO signal_events (time, product_code, side, price, size) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(cmd, s.Time.Format(timeFormat), s.ProductCode, s.Side, s.Price, s.Size)
	if err != nil {
		fmt.Println("[Save]", err)
		// 重複エラーであれば問題ない
		if strings.Contains(err.Error(), "Duplicate entry") {
			return true
		}
		return false
	}
	return true
}

type SignalEvents struct {
	Signals []SignalEvent `json:"signals,omitempty"`
	Profit  float64       `json:"profit"`
}

func NewSignalEvents() *SignalEvents {
	return &SignalEvents{}
}

func GetSignalEventsByProductCode(db DB, productCode string) *SignalEvents {
	cmd := "SELECT * FROM signal_events WHERE product_code = ? ORDER BY time ASC"
	rows, err := db.Query(cmd, productCode)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var signalEvents SignalEvents
	for rows.Next() {
		var signalEvent SignalEvent
		rows.Scan(&signalEvent.Time, &signalEvent.ProductCode, &signalEvent.Side, &signalEvent.Price, &signalEvent.Size)
		signalEvents.Signals = append(signalEvents.Signals, signalEvent)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}

	return &signalEvents
}

func (s *SignalEvents) CanBuy(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return true
	}

	lastSignal := s.Signals[lenSignals-1]
	canBuy := lastSignal.Side == "SELL" && lastSignal.Time.Before(time)
	return canBuy
}

func (s *SignalEvents) CanSell(time time.Time) bool {
	lenSignals := len(s.Signals)
	if lenSignals == 0 {
		return false
	}

	lastSignal := s.Signals[lenSignals-1]
	canSell := lastSignal.Side == "BUY" && lastSignal.Time.Before(time)
	return canSell
}

func (s *SignalEvents) Buy(productCode string, time time.Time, price, size float64) bool {
	if !s.CanBuy(time) {
		return false
	}
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
	if !s.CanSell(time) {
		return false
	}
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
	for _, signalEvent := range s.Signals {
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
