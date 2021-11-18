package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type candleRepository struct {
	db              DB
	candleTableName string
	timeFormat      string
}

func NewCandleRepository(db DB, candleTableName, timeFormat string) repository.CandleRepository {
	if candleTableName == "" {
		return nil
	}

	return &candleRepository{
		db:              db,
		candleTableName: candleTableName,
		timeFormat:      timeFormat,
	}
}

func (cr candleRepository) Save(candle model.Candle) error {
	cmd := fmt.Sprintf(`
        INSERT INTO %s
            (time, open, close, high, low, volume)
        VALUES
            (?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            open = VALUES(open),
            close = VALUES(close),
            high = VALUES(high),
            low = VALUES(low),
            volume = VALUES(volume)
        `,
		cr.candleTableName,
	)
	_, err := cr.db.Exec(cmd, candle.Time().Format(cr.timeFormat), candle.Open(), candle.Close(), candle.High(), candle.Low(), candle.Volume())
	return err
}

func (cr candleRepository) FindByCandleTime(productCode string, duration time.Duration, candleTime model.CandleTime) (*model.Candle, error) {
	cmd := fmt.Sprintf(`
        SELECT
            open, close, high, low, volume
        FROM
            %s
        WHERE
            time = ?
        `,
		cr.candleTableName,
	)
	row := cr.db.QueryRow(cmd, candleTime.Format(cr.timeFormat))

	var candleOpen, candleClose, candleHigh, candleLow, candleVolume float64
	err := row.Scan(&candleOpen, &candleClose, &candleHigh, &candleLow, &candleVolume)
	// 発見できなかったらそのままnilを返す
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	candle := model.NewCandle(productCode, duration, candleTime, candleOpen, candleClose, candleHigh, candleLow, candleVolume)
	if candle == nil {
		return nil, errors.New(fmt.Sprint("invalid candle:", productCode, duration, candleTime, candleOpen, candleClose, candleHigh, candleLow, candleVolume))
	}
	return candle, nil
}

func (cr candleRepository) FindAll(productCode string, duration time.Duration, limit int64) ([]model.Candle, error) {
	cmd := fmt.Sprintf(`
        SELECT
            *
        FROM (
            SELECT
                time, open, close, high, low, volume
            FROM
                %s
            ORDER BY
                time DESC
            LIMIT ?
        ) AS candle
        ORDER BY
            time ASC
        `,
		cr.candleTableName,
	)
	rows, err := cr.db.Query(cmd, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candles := make([]model.Candle, 0)
	for rows.Next() {
		var timeTime time.Time
		var candleOpen, candleClose, candleHigh, candleLow, candleVolume float64
		err := rows.Scan(&timeTime, &candleOpen, &candleClose, &candleHigh, &candleLow, &candleVolume)
		if err != nil {
			return nil, err
		}

		candleTime := model.NewCandleTime(timeTime)
		candle := model.NewCandle(productCode, duration, candleTime, candleOpen, candleClose, candleHigh, candleLow, candleVolume)
		if candle == nil {
			return nil, errors.New(fmt.Sprint("invalid candle:", productCode, duration, candleTime, candleOpen, candleClose, candleHigh, candleLow, candleVolume))
		}

		candles = append(candles, *candle)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return candles, nil
}
