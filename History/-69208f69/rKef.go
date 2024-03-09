package common

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func MaxFloat64(slice []float64) float64 {
	maxValue := slice[0]
	for _, value := range slice {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}

func MinFloat64(slice []float64) float64 {
	minValue := slice[0]
	for _, value := range slice {
		if value < minValue {
			minValue = value
		}
	}

	return minValue
}
