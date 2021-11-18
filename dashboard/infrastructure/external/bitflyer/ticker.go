package bitflyer

import (
	"encoding/json"
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

// 板の状態
type BoardState string

const (
	BoardStateRunning      BoardState = "RUNNING"       // 通常稼働中
	BoardStateClosed       BoardState = "CLOSED"        // 取引停止中
	BoardStateStarting     BoardState = "STARTING"      // 再起動中
	BoardStatePreopen      BoardState = "PREOPEN"       // 板寄せ中
	BoardStateCircuitBreak BoardState = "CIRCUIT BREAK" // サーキットブレイク発動中
	BoardStateAwaitingSQ   BoardState = "AWAITING SQ"   // Lightning Futures の取引終了後 SQ（清算値）の確定前
	BoardStateMatured      BoardState = "MATURED"       // Lightning Futures の満期に到達
)

const TimestampFormat = "2006-01-02T15:04:05"

type Ticker struct {
	ProductCode     string     `json:"product_code"`
	State           BoardState `json:"state"`
	Timestamp       string     `json:"timestamp"`
	TickID          int        `json:"tick_id"`
	BestBid         float64    `json:"best_bid"`
	BestAsk         float64    `json:"best_ask"`
	BestBidSize     float64    `json:"best_bid_size"`
	BestAskSize     float64    `json:"best_ask_size"`
	TotalBidDepth   float64    `json:"total_bid_depth"`
	TotalAskDepth   float64    `json:"total_ask_depth"`
	MarketBidSize   float64    `json:"market_bid_size"`
	MarketAskSize   float64    `json:"market_ask_size"`
	Ltp             float64    `json:"ltp"`
	Volume          float64    `json:"volume"`
	VolumeByProduct float64    `json:"volume_by_product"`
}

func (ticker *Ticker) toDomainModelTicker() *model.Ticker {
	return model.NewTicker(
		ticker.ProductCode,
		string(ticker.State),
		ticker.Timestamp,
		ticker.TickID,
		ticker.BestBid,
		ticker.BestAsk,
		ticker.BestBidSize,
		ticker.BestAskSize,
		ticker.TotalBidDepth,
		ticker.TotalAskDepth,
		ticker.MarketBidSize,
		ticker.MarketAskSize,
		ticker.Ltp,
		ticker.Volume,
		ticker.VolumeByProduct,
	)
}

type bitflyerTickerRepository struct {
	apiClient *Client
}

func NewBitflyerTickerRepository(apiClient *Client) repository.TickerRepository {
	return &bitflyerTickerRepository{
		apiClient: apiClient,
	}
}

func (btr *bitflyerTickerRepository) Fetch(productCode string) (*model.Ticker, error) {
	path := "ticker"
	query := map[string]string{"product_code": productCode}
	resp, err := btr.apiClient.doRequest("GET", path, query, nil)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return nil, err
	}
	if ticker.State != BoardStateRunning {
		return nil, errors.New("bitflyer is not running")
	}

	domainModelTicker := ticker.toDomainModelTicker()
	if domainModelTicker == nil {
		return nil, errors.New("invalid ticker fetched")
	}

	return domainModelTicker, nil
}
