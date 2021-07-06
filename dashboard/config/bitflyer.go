package config

import (
	"os"
)

var (
	APIKey      string
	APISecret   string
	ProductCode string
	TradeHour   int
)

func init() {
	APIKey = os.Getenv("BITFLYER_API_KEY")
	APISecret = os.Getenv("BITFLYER_API_SECRET")
	ProductCode = os.Getenv("PRODUCT_CODE")
	TradeHour = 9
}
