package common

import "math"

type IsHighestRule struct {
	values   []float64
	barCount int
}

func NewIsHighestRule(values []float64, barCount int) Rule {
	return &IsHighestRule{
		values:   values,
		barCount: barCount,
	}
}

func (ihr *IsHighestRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(ihr.values) || ihr.barCount <= 0 {
		return false, nil
	}

	highest := calculateHighestValue(ihr.values, index, ihr.barCount)

	satisfied := !isNaN(ihr.values[index]) && !isNaN(highest) && ihr.values[index] == highest

	return satisfied, nil
}

func isNaN(val float64) bool {
	return math.IsNaN(val)
}

func calculateHighestValue(values []float64, currentIndex, barCount int) float64 {
	highest := math.Inf(-1)

	for i := maxInt(0, currentIndex-barCount+1); i <= currentIndex; i++ {
		if values[i] > highest {
			highest = values[i]
		}
	}

	return highest
}