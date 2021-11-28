package service

import (
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type CandleService interface {
	TickerToCandle(ticker model.Ticker) *model.Candle
	Update(oldCandle, newCandle *model.Candle) *model.Candle
	Save(candle model.Candle) error
	FindByTime(productCode string, timeTime time.Time) (*model.Candle, error)
	FindAll(productCode string, limit int64) ([]model.Candle, error)
}

// 日足
type candleServicePerDay struct {
	localTime        *time.Location
	tradeHour        int
	candleRepository repository.CandleRepository
}

func NewCandleServicePerDay(lt *time.Location, th int, cr repository.CandleRepository) CandleService {
	return &candleServicePerDay{
		localTime:        lt,
		tradeHour:        th,
		candleRepository: cr,
	}
}

func (cs *candleServicePerDay) Duration() time.Duration {
	return 24 * time.Hour
}

func (cs *candleServicePerDay) TickerToCandle(ticker model.Ticker) *model.Candle {
	price := ticker.MidPrice()

	tickerTime := model.NewCandleTimeByString(ticker.Timestamp())
	candleTime := tickerTime.TruncateHour(cs.localTime, cs.tradeHour)

	return model.NewCandle(ticker.ProductCode(), cs.Duration(), candleTime, price, price, price, price, ticker.Volume())
}

func (cs *candleServicePerDay) Update(oldCandle, newCandle *model.Candle) *model.Candle {
	if oldCandle == nil || newCandle == nil {
		return newCandle
	}

	if oldCandle.ProductCode() != newCandle.ProductCode() {
		return nil
	}

	if oldCandle.Duration() != newCandle.Duration() {
		return nil
	}

	if !oldCandle.Time().Equal(newCandle.Time()) {
		return newCandle
	}

	high := oldCandle.High()
	if high < newCandle.High() {
		high = newCandle.High()
	}

	low := oldCandle.Low()
	if low > newCandle.Low() {
		low = newCandle.Low()
	}

	return model.NewCandle(oldCandle.ProductCode(), oldCandle.Duration(), oldCandle.Time(), oldCandle.Open(), newCandle.Close(), high, low, newCandle.Volume())
}

func (cs *candleServicePerDay) Save(candle model.Candle) error {
	return cs.candleRepository.Save(candle)
}

func (cs *candleServicePerDay) FindByTime(productCode string, timeTime time.Time) (*model.Candle, error) {
	candleTime := model.NewCandleTime(timeTime)
	return cs.candleRepository.FindByCandleTime(productCode, cs.Duration(), candleTime)
}

func (cs *candleServicePerDay) FindAll(productCode string, limit int64) ([]model.Candle, error) {
	return cs.candleRepository.FindAll(productCode, cs.Duration(), limit)
}
