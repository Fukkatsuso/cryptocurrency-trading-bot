package model_test

import (
	"testing"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestNewCandleTime(t *testing.T) {
	table := []struct {
		time time.Time
	}{
		{
			time: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time: time.Date(2021, time.January, 2, 9, 0, 0, 0, time.UTC),
		},
		{
			time: time.Date(2021, time.January, 3, 12, 0, 0, 0, time.UTC),
		},
		{
			time: time.Date(2021, time.January, 4, 23, 59, 59, 990000000, time.UTC),
		},
	}

	for _, c := range table {
		candleTime := model.NewCandleTime(c.time)
		if !candleTime.Time().Equal(c.time) {
			t.Fatalf("%v != %v", candleTime.Time(), c.time)
		}
	}
}

func TestNewCandleTimeByString(t *testing.T) {
	table := []struct {
		timeString string
		time       time.Time
	}{
		{
			timeString: "2021-01-01T00:00:00.00",
			time:       time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			timeString: "2021-01-02T09:00:00.00",
			time:       time.Date(2021, time.January, 2, 9, 0, 0, 0, time.UTC),
		},
		{
			timeString: "2021-01-03T12:00:00.00",
			time:       time.Date(2021, time.January, 3, 12, 0, 0, 0, time.UTC),
		},
		{
			timeString: "2021-01-04T23:59:59.99",
			time:       time.Date(2021, time.January, 4, 23, 59, 59, 990000000, time.UTC),
		},
	}

	for _, c := range table {
		candleTime := model.NewCandleTimeByString(c.timeString)
		if !candleTime.Time().Equal(c.time) {
			t.Fatalf("%v != %v", candleTime.Time(), c.time)
		}
	}
}

func TestCandleTimeTruncateHour(t *testing.T) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	table := []struct {
		time         time.Time
		location     *time.Location
		hour         int
		expectedTime time.Time
	}{
		// UTC
		{
			time:         time.Date(2021, time.January, 1, 2, 3, 4, 5, time.UTC),
			location:     time.UTC,
			hour:         0,
			expectedTime: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			location:     time.UTC,
			hour:         1,
			expectedTime: time.Date(2020, time.December, 31, 1, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
			location:     time.UTC,
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 10, 0, 0, 0, time.UTC),
			location:     time.UTC,
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 9, 0, 0, 0, time.UTC),
			location:     time.UTC,
			hour:         15,
			expectedTime: time.Date(2020, time.December, 31, 15, 0, 0, 0, time.UTC),
		},
		// jst
		{
			time:         time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			location:     jst,
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 1, 0, 0, 0, time.UTC),
			location:     jst,
			hour:         9,
			expectedTime: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			time:         time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			location:     jst,
			hour:         15,
			expectedTime: time.Date(2020, time.December, 31, 6, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range table {
		candleTime := model.NewCandleTime(c.time)
		truncatedTime := candleTime.TruncateHour(c.location, c.hour)
		expectedCandleTime := model.NewCandleTime(c.expectedTime)
		if !truncatedTime.Equal(expectedCandleTime) {
			t.Fatalf("%v != %v", truncatedTime.Time(), expectedCandleTime.Time())
		}
	}
}
