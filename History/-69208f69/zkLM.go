package common

func MaxInt(values ...int) int {
	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
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
