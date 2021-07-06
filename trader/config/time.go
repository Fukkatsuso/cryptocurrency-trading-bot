package config

import (
	"time"
)

var (
	LocalTime *time.Location
)

func init() {
	LocalTime = time.FixedZone("Asia/Tokyo", 9*60*60)
}
