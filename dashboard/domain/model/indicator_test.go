package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
)

func TestSMA(t *testing.T) {
	inReal := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var sma *model.SMA

	sma = model.NewSMA(inReal, 3)
	if sma == nil {
		t.Fatal("NewSMA() returns nil")
	}

	sma = model.NewSMA(inReal, -1)
	if sma != nil {
		t.Fatal("NewSMA() returns not nil")
	}

	sma = model.NewSMA(inReal, 20)
	if sma != nil {
		t.Fatal("NewSMA() returns not nil")
	}
}

func TestEMA(t *testing.T) {
	inReal := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var ema *model.EMA

	ema = model.NewEMA(inReal, 3)
	if ema == nil {
		t.Fatal("NewEMA() returns nil")
	}
	if ema.Period() != 3 {
		t.Fatalf("%d != %d", ema.Period(), 3)
	}

	ema = model.NewEMA(inReal, -1)
	if ema != nil {
		t.Fatal("NewEMA() returns not nil")
	}

	ema = model.NewEMA(inReal, 20)
	if ema != nil {
		t.Fatal("NewEMA() returns not nil")
	}
}

func TestBBands(t *testing.T) {
	inReal := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var bbands *model.BBands

	bbands = model.NewBBands(inReal, 3, 2.0)
	if bbands == nil {
		t.Fatal("NewBBands() returns nil")
	}
	if bbands.N() != 3 {
		t.Fatalf("%d != %d", bbands.N(), 3)
	}
	if bbands.K() != 2.0 {
		t.Fatalf("%f != %f", bbands.K(), 2.0)
	}

	bbands = model.NewBBands(inReal, -1, 2.0)
	if bbands != nil {
		t.Fatal("NewBBands() returns not nil")
	}

	bbands = model.NewBBands(inReal, 20, 2.0)
	if bbands != nil {
		t.Fatal("NewBBands() returns nil")
	}

	bbands = model.NewBBands(inReal, 3, -1.0)
	if bbands != nil {
		t.Fatal("NewBBands() returns nil")
	}
}

func TestIchimokuCloud(t *testing.T) {
	var inReal []float64
	var ichimokuCloud *model.IchimokuCloud

	inReal = newSequence(52)
	ichimokuCloud = model.NewIchimokuCloud(inReal)
	if ichimokuCloud == nil {
		t.Fatal("NewIchimokuCloud() returns nil")
	}

	inReal = newSequence(51)
	ichimokuCloud = model.NewIchimokuCloud(inReal)
	if ichimokuCloud != nil {
		t.Fatal("NewIchimokuCloud() returns not nil")
	}
}

func newSequence(length int) []float64 {
	seq := make([]float64, length)
	for i := 0; i < length; i++ {
		seq[i] = float64(i + 1)
	}
	return seq
}

func TestRSI(t *testing.T) {
	inReal := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var rsi *model.RSI

	rsi = model.NewRSI(inReal, 3)
	if rsi == nil {
		t.Fatal("NewRSI() returns nil")
	}

	rsi = model.NewRSI(inReal, -1)
	if rsi != nil {
		t.Fatal("NewRSI() returns not nil")
	}

	rsi = model.NewRSI(inReal, 20)
	if rsi != nil {
		t.Fatal("NewRSI() returns not nil")
	}
}

func TestMACD(t *testing.T) {
	inReal := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var macd *model.MACD

	macd = model.NewMACD(inReal, 3, 7, 5)
	if macd == nil {
		t.Fatal("NewMACD() returns nil")
	}

	macd = model.NewMACD(inReal, -1, 7, 5)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}

	macd = model.NewMACD(inReal, 20, 7, 5)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}

	macd = model.NewMACD(inReal, 3, -1, 5)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}

	macd = model.NewMACD(inReal, 3, 20, 5)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}

	macd = model.NewMACD(inReal, 3, 7, -1)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}

	macd = model.NewMACD(inReal, 3, 7, 20)
	if macd != nil {
		t.Fatal("NewMACD() returns not nil")
	}
}
