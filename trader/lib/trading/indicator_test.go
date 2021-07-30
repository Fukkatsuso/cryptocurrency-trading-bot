package trading

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	inReal := []float64{
		1, 2, 3, 4, 5, -1, -2, 0, 0.1, -0.1, 10,
	}

	expectedMin, expectedMax := float64(-2), float64(10)
	min, max := MinMax(inReal)

	if min != expectedMin {
		t.Fatalf("MinMax returns different min: %f != %f", min, expectedMin)
	}
	if max != expectedMax {
		t.Fatalf("MinMax returns different max: %f != %f", max, expectedMax)
	}
}

func TestMin(t *testing.T) {
	x, y := 1, 100

	expectedMin := 1
	min := Min(x, y)

	if min != expectedMin {
		t.Fatalf("MinMax returns different min: %d != %d", min, expectedMin)
	}
}

func TestIchimokuCloud(t *testing.T) {
	inReal := []float64{
		263986, 267338, 279845, 267536, 247428.5, 192056.5, 222448.5, 222774.5, 245200.5, 257669,
	}

	// tenkan, kijun, senkouA, senkouB, chikou
	_, _, _, _, _ = IchimokuCloud(inReal)
}
