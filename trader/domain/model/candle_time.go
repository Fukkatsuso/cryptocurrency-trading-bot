package model

import (
	"fmt"
	"time"
)

type CandleTime time.Time

func NewCandleTime(timeTime time.Time) CandleTime {
	timeTime = timeTime.In(time.UTC)
	return CandleTime(timeTime)
}

// UTCでパースして返す
func NewCandleTimeByString(timeString string) CandleTime {
	timeTime, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		fmt.Println("[DateTimeUTC]", err)
		return NewCandleTime(time.Time{})
	}
	return NewCandleTime(timeTime)
}

func (candleTime CandleTime) Time() time.Time {
	return time.Time(candleTime)
}

func (candleTime CandleTime) Format(layout string) string {
	return candleTime.Time().Format(layout)
}

func (candleTime CandleTime) Equal(compareTime CandleTime) bool {
	return candleTime.Time().Equal(compareTime.Time())
}

// hour時を境に切り捨てた時間
func (candleTime CandleTime) TruncateHour(localTime *time.Location, hour int) CandleTime {
	truncateTime := candleTime.Time().In(localTime).Truncate(time.Hour)

	// [0, hour)時の場合，日付を1日戻す
	if truncateTime.Hour() < hour {
		truncateTime = truncateTime.Add(-24 * time.Hour)
	}

	// Hourをhourに揃える
	offset := time.Duration(hour - truncateTime.Hour())
	truncateTime = truncateTime.Add(offset * time.Hour)

	return NewCandleTime(truncateTime)
}
