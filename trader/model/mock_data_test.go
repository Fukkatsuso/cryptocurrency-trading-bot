package model

import (
	"testing"
)

func TestCandleMockData(t *testing.T) {
	_, err := CandleMockData()
	if err != nil {
		t.Fatal(err.Error())
	}
}
