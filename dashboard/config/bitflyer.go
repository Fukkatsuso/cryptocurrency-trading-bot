package config

import (
	"os"
	"time"
)

var (
	APIKey         string
	APISecret      string
	ProductCode    string
	CandleDuration time.Duration
	TradeHour      int
)

func init() {
	APIKey = os.Getenv("BITFLYER_API_KEY")
	APISecret = os.Getenv("BITFLYER_API_SECRET")
	ProductCode = os.Getenv("PRODUCT_CODE")
	CandleDuration = 24 * time.Hour
	TradeHour = 9
}
