package common

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func Max(slice []float64) float64 {
	maxValue := slice[0]
	for _, value := range slice {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}

func Min(slice []float64) float64 {
	minValue := slice[0]
	for _, value := range slice {
		if value < minValue {
			minValue = value
		}
	}

	return minValue
}
