package common

import "math"

type IsLowestRule struct {
	values   []float64
	barCount int
}

func NewIsLowestRule(values []float64, barCount int) Rule {
	return &IsLowestRule{
		values:   values,
		barCount: barCount,
	}
}

func (ilr *IsLowestRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(ilr.values) || ilr.barCount <= 0 {
		return false, nil
	}

	lowest := calculateLowestValue(ilr.values, index, ilr.barCount)
	refValue := ilr.values[index]

	satisfied := !isNaN(refValue) && !isNaN(lowest) && refValue == lowest

	return satisfied, nil
}

func isNaN(val float64) bool {
	return math.IsNaN(val)
}

func calculateLowestValue(values []float64, currentIndex, barCount int) float64 {
	lowest := math.Inf(1) // Positive infinity as an initial value
	for i := max(0, currentIndex-barCount+1); i <= currentIndex; i++ {
		if values[i] < lowest {
			lowest = values[i]
		}
	}
	return lowest
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
