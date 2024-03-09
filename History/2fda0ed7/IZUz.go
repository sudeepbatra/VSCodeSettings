package rules

import (
	"math"

	"github.com/sudeepbatra/alpha-hft/common"
)

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

func (ihr *IsHighestRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(ihr.values) || ihr.barCount <= 0 {
		return false
	}

	highest := calculateHighestValue(ihr.values, index, ihr.barCount)

	return !math.IsNaN(ihr.values[index]) && !math.IsNaN(highest) && ihr.values[index] == highest
}

func calculateHighestValue(values []float64, currentIndex, barCount int) float64 {
	highest := math.Inf(-1)

	for i := common.MaxInt(0, currentIndex-barCount+1); i <= currentIndex; i++ {
		if values[i] > highest {
			highest = values[i]
		}
	}

	return highest
}
