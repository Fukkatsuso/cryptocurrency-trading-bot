package model

import (
	"fmt"
	"time"

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

func (c *Candle) Create(db DB, candleTableName, timeFormat string) error {
	cmd := fmt.Sprintf("INSERT INTO %s (time, open, close, high, low, volume) VALUES (?, ?, ?, ?, ?, ?)", candleTableName)
	_, err := db.Exec(cmd, c.Time.Format(timeFormat), c.Open, c.Close, c.High, c.Low, c.Volume)
	if err != nil {
		return err
	}
	return err
}

func (c *Candle) Save(db DB, candleTableName, timeFormat string) error {
	cmd := fmt.Sprintf("UPDATE %s SET open = ?, close = ?, high = ?, low = ?, volume = ? WHERE time = ?", candleTableName)
	_, err := db.Exec(cmd, c.Open, c.Close, c.High, c.Low, c.Volume, c.Time.Format(timeFormat))
	if err != nil {
		return err
	}
	return err
}

func GetCandle(db DB, candleTableName, timeFormat, productCode string,
	duration time.Duration, dateTime time.Time) *Candle {
	cmd := fmt.Sprintf("SELECT time, open, close, high, low, volume FROM %s WHERE time = ?", candleTableName)
	row := db.QueryRow(cmd, dateTime.Format(timeFormat))
	var candle Candle
	err := row.Scan(&candle.Time, &candle.Open, &candle.Close, &candle.High, &candle.Low, &candle.Volume)
	if err != nil {
		return nil
	}
	return NewCandle(productCode, duration, candle.Time, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
}

// UTCでのtime.Timeを返す
func DateTimeUTC(timeString string) time.Time {
	dateTime, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		fmt.Println("[DateTimeUTC]", err)
	}
	return dateTime.In(time.UTC)
}

// hour時で1日を区切りTruncate
func TruncateDateTime(dateTime time.Time, hour int) time.Time {
	dateTime = dateTime.Truncate(time.Hour)

	// [0, hour)時の場合，日付を1日戻す
	if dateTime.Hour() < hour {
		dateTime = dateTime.Add(-24 * time.Hour)
	}

	// Hourをhourに揃える
	offset := time.Duration(hour - dateTime.Hour())
	dateTime = dateTime.Add(offset * time.Hour)

	return dateTime
}

func CreateCandleWithDuration(db DB, candleTableName, timeFormat string, localTime *time.Location, tradeHour int,
	ticker *bitflyer.Ticker, productCode string, duration time.Duration) error {
	// LocalTimeでTruncateしたものをUTCに戻す
	dateTime := DateTimeUTC(ticker.Timestamp)
	truncateDateTime := TruncateDateTime(dateTime.In(localTime), tradeHour).In(time.UTC)
	price := ticker.GetMidPrice()
	currentCandle := GetCandle(db, candleTableName, timeFormat,
		productCode, duration, truncateDateTime)

	if currentCandle == nil {
		candle := NewCandle(productCode, duration, truncateDateTime, price, price, price, price, ticker.Volume)
		return candle.Create(db, candleTableName, timeFormat)
	}

	if currentCandle.High <= price {
		currentCandle.High = price
	} else if currentCandle.Low >= price {
		currentCandle.Low = price
	}
	currentCandle.Volume = ticker.Volume
	currentCandle.Close = price
	return currentCandle.Save(db, candleTableName, timeFormat)
}

func GetAllCandle(db DB, candleTableName, timeFormat string,
	productCode string, duration time.Duration, limit int) (candles []Candle, err error) {
	cmd := fmt.Sprintf("SELECT * FROM (SELECT time, open, close, high, low, volume FROM %s ORDER BY time DESC LIMIT ?) AS candle ORDER BY time ASC", candleTableName)
	rows, err := db.Query(cmd, limit)
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
