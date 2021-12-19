package bitflyer

import (
	"errors"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/repository"
)

type bitflyerTickerMockRepository struct {
}

func NewBitflyerTickerMockRepository() repository.TickerRepository {
	return &bitflyerTickerMockRepository{}
}

func (btr *bitflyerTickerMockRepository) Fetch(productCode string) (*model.Ticker, error) {
	// 元データ: ETH_JPY @2021-11-09T11:31:11.797
	ticker := Ticker{
		ProductCode:     productCode,
		State:           BoardStateRunning,
		Timestamp:       time.Now().UTC().Format(TimestampFormat),
		TickID:          2946055,
		BestBid:         540284,
		BestAsk:         540437,
		BestBidSize:     540284,
		BestAskSize:     540437,
		TotalBidDepth:   3506.8420941,
		TotalAskDepth:   1972.7725585,
		MarketBidSize:   0,
		MarketAskSize:   0,
		Ltp:             540284,
		Volume:          12507.4724701,
		VolumeByProduct: 12507.4724701,
	}

	domainModelTicker := ticker.toDomainModelTicker()
	if domainModelTicker == nil {
		return nil, errors.New("invalid ticker fetched")
	}

	return domainModelTicker, nil
}
