package model

import (
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/lib/bitflyer"
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

func (c *Candle) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, high, low, volume) VALUES (?, ?, ?, ?, ?, ?)", config.CandleTableName)
	_, err := config.DB.Exec(cmd, c.Time.Format(config.TimeFormat), c.Open, c.Close, c.High, c.Low, c.Volume)
	if err != nil {
		return err
	}
	return err
}

func (c *Candle) Save() error {
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?", config.CandleTableName)
	_, err := config.DB.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(config.TimeFormat))
	if err != nil {
		return err
	}
	return err
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

// tickerのTimestampをフォーマット付きLocalTimeに変換
func DateTimeLocal(timeString string) time.Time {
	dateTimeUTC, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		fmt.Println("[DateTimeLocal]", err)
	}

	return dateTimeUTC.In(config.LocalTime)
}

// hour時で1日を区切ってTruncate
func TruncateDateTimeLocal(timeString string, hour int) time.Time {
	dateTime := DateTimeLocal(timeString).Truncate(time.Hour)

	// [0, hour)時の場合，日付を1日戻す
	if dateTime.Hour() < hour {
		dateTime = dateTime.Add(-24 * time.Hour)
	}

	// Hourをhourに揃える
	offset := time.Duration(hour - dateTime.Hour())
	return dateTime.Add(offset * time.Hour)
}

func CreateCandleWithDuration(ticker *bitflyer.Ticker, productCode string, duration time.Duration) error {
	truncateDateTime := TruncateDateTimeLocal(ticker.Timestamp, config.TradeHour)
	price := ticker.GetMidPrice()
	currentCandle := GetCandle(productCode, duration, truncateDateTime)

	if currentCandle == nil {
		candle := NewCandle(productCode, duration, truncateDateTime, price, price, price, price, ticker.Volume)
		return candle.Create()
	}

	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume += ticker.Volume
	currentCandle.Close = price
	return currentCandle.Save()
}
