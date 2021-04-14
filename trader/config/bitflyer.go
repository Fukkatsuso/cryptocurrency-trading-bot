package config

import "os"

var (
	APIKey      string
	APISecret   string
	ProductCode string
)

func init() {
	APIKey = os.Getenv("BITFLYER_API_KEY")
	APISecret = os.Getenv("BITFLYER_API_SECRET")
	ProductCode = os.Getenv("PRODUCT_CODE")
}
