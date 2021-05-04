package model

import (
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
)

type Candle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Time        time.Time     `json:"time"`
	Open        float64       `json:"open"`
	Close       float64       `json:"close"`
	High        float64       `json:"high"`
	Low         float64       `json:"low"`
	Volume      float64       `json:"volume"`
}

func NewCandle(productCode string, duration time.Duration, timeDate time.Time, open, close, high, low, volume float64) *Candle {
	return &Candle{
		productCode,
		duration,
		timeDate,
		open,
		close,
		high,
		low,
		volume,
	}
}

func GetCandle(productCode string, duration time.Duration, dateTime time.Time) *Candle {
	cmd := fmt.Sprintf("SELECT time, open, close, high, low, volume FROM %s WHERE time = ?", config.CandleTableName)
	row := config.DB.QueryRow(cmd, dateTime.Format(config.TimeFormat))
	var candle Candle
	err := row.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(productCode, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

func GetAllCandle(productCode string, duration time.Duration, limit int) (candles []Candle, err error) {
	cmd := fmt.Sprintf("SELECT * FROM (SELECT time, open, close, high, low, volume FROM %s ORDER BY time DESC LIMIT ?) AS candle ORDER BY time ASC", config.CandleTableName)
	rows, err := config.DB.Query(cmd, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	candles = make([]Candle, 0)
	for rows.Next() {
		var candle Candle
		candle.ProductCode = productCode
		candle.Duration = duration
		rows.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
		candles = append(candles, candle)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return candles, nil
}
