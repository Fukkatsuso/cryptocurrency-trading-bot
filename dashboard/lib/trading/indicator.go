package trading

func MinMax(inReal []float64) (float64, float64) {
	min := inReal[0]
	max := inReal[0]
	for _, price := range inReal {
		if min > price {
			min = price
		}
		if max < price {
			max = price
		}
	}
	return min, max
}

func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func IchimokuCloud(inReal []float64) ([]float64, []float64, []float64, []float64, []float64) {
	length := len(inReal)
	tenkan := make([]float64, Min(9, length))
	kijun := make([]float64, Min(26, length))
	senkouA := make([]float64, Min(26, length))
	senkouB := make([]float64, Min(52, length))
	chikou := make([]float64, Min(26, length))

	for i := range inReal {
		if i >= 9 {
			min, max := MinMax(inReal[i-9 : i])
			tenkan = append(tenkan, (min+max)/2)
		}
		if i >= 26 {
			min, max := MinMax(inReal[i-26 : i])
			kijun = append(kijun, (min+max)/2)
			senkouA = append(senkouA, (tenkan[i]+kijun[i])/2)
			chikou = append(chikou, inReal[i-26])
		}
		if i >= 52 {
			min, max := MinMax(inReal[i-52 : i])
			senkouB = append(senkouB, (min+max)/2)
		}
	}
	return tenkan, kijun, senkouA, senkouB, chikou
}
