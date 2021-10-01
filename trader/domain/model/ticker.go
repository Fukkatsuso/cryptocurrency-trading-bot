package model

type Ticker struct {
	productCode     string
	state           string
	timestamp       string
	tickID          int
	bestBid         float64
	bestAsk         float64
	bestBidSize     float64
	bestAskSize     float64
	totalBidDepth   float64
	totalAskDepth   float64
	marketBidSize   float64
	marketAskSize   float64
	ltp             float64
	volume          float64
	volumeByProduct float64
}

func NewTicker(productCode, state string, timestamp string, tickID int, bestBid, bestAsk, bestBidSize, bestAskSize, totalBidDepth, totalAskDepth, marketBidSize, marketAskSize, ltp, volume, volumeByProduct float64) *Ticker {
	if productCode == "" {
		return nil
	}

	if state == "" {
		return nil
	}

	if timestamp == "" {
		return nil
	}

	if tickID < 0 {
		return nil
	}

	if bestBid <= 0 {
		return nil
	}

	if bestAsk <= 0 {
		return nil
	}

	if bestBidSize < 0 {
		return nil
	}

	if bestAskSize < 0 {
		return nil
	}

	if totalBidDepth < 0 {
		return nil
	}

	if totalAskDepth < 0 {
		return nil
	}

	if marketBidSize < 0 {
		return nil
	}

	if marketAskSize < 0 {
		return nil
	}

	if ltp <= 0 {
		return nil
	}

	if volume < 0 {
		return nil
	}

	if volumeByProduct < 0 {
		return nil
	}

	return &Ticker{
		productCode:     productCode,
		state:           state,
		timestamp:       timestamp,
		tickID:          tickID,
		bestBid:         bestBid,
		bestAsk:         bestAsk,
		bestBidSize:     bestBid,
		bestAskSize:     bestAsk,
		totalBidDepth:   totalBidDepth,
		totalAskDepth:   totalAskDepth,
		marketBidSize:   marketBidSize,
		marketAskSize:   marketAskSize,
		ltp:             ltp,
		volume:          volume,
		volumeByProduct: volume,
	}
}

func (t *Ticker) Timestamp() string {
	return t.timestamp
}

func (t *Ticker) MidPrice() float64 {
	return (t.bestBid + t.bestAsk) / 2
}

func (t *Ticker) ProductCode() string {
	return t.productCode
}

func (t *Ticker) Volume() float64 {
	return t.volume
}
