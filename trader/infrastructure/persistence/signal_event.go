package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type signalEventRepository struct {
	db         DB
	timeFormat string
}

func NewSignalEventRepository(db DB, timeFormat string) repository.SignalEventRepository {
	return &signalEventRepository{
		db:         db,
		timeFormat: timeFormat,
	}
}

func (sr *signalEventRepository) Save(signal model.SignalEvent) error {
	cmd := `
        INSERT INTO signal_events
            (time, product_code, side, price, size)
        VALUES
            (?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            time = VALUES(time)
        `
	_, err := sr.db.Exec(cmd, signal.Time().Format(sr.timeFormat), signal.ProductCode(), signal.Side(), signal.Price(), signal.Size())

	return err
}

func (sr *signalEventRepository) FindAll(productCode string) ([]model.SignalEvent, error) {
	cmd := `
        SELECT
            *
        FROM signal_events
        WHERE
            product_code = ?
        ORDER BY
            time ASC
        `
	rows, err := sr.db.Query(cmd, productCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signalEvents := []model.SignalEvent{}
	for rows.Next() {
		var timeTime time.Time
		var productCode string
		var side model.OrderSide
		var price, size float64
		err := rows.Scan(&timeTime, &productCode, &side, &price, &size)
		if err != nil {
			return nil, err
		}

		signalEvent := model.NewSignalEvent(timeTime, productCode, side, price, size)
		if signalEvent == nil {
			return nil, errors.New(fmt.Sprint("invalid signal_event:", timeTime, productCode, side, price, size))
		}

		signalEvents = append(signalEvents, *signalEvent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return signalEvents, nil
}

func (sr *signalEventRepository) FindAllAfterTime(productCode string, timeTime time.Time) ([]model.SignalEvent, error) {
	cmd := `
        SELECT
            *
        FROM
            signal_events
        WHERE
            product_code = ? AND
            time >= ?
        ORDER BY
            time ASC
        `
	rows, err := sr.db.Query(cmd, productCode, timeTime.Format(sr.timeFormat))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	signalEvents := []model.SignalEvent{}
	for rows.Next() {
		var timeTime time.Time
		var productCode string
		var side model.OrderSide
		var price, size float64
		err := rows.Scan(&timeTime, &productCode, &side, &price, &size)
		if err != nil {
			return nil, err
		}

		signalEvent := model.NewSignalEvent(timeTime, productCode, side, price, size)
		if signalEvent == nil {
			return nil, errors.New(fmt.Sprint("invalid signal_event:", timeTime, productCode, side, price, size))
		}

		signalEvents = append(signalEvents, *signalEvent)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return signalEvents, nil
}
