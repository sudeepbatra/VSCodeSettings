package rules

import "github.com/sudeepbatra/alpha-hft/common"

type IsRisingRule struct {
	values      []float64
	barCount    int
	minStrength float64
}

func NewIsRisingRule(values []float64, barCount int) Rule {
	return &IsRisingRule{
		values:      values,
		barCount:    barCount,
		minStrength: 1.0,
	}
}

func NewIsRisingRuleWithStrength(values []float64, barCount int, minStrength float64) Rule {
	if minStrength >= 1 {
		minStrength = 0.99
	}

	return &IsRisingRule{
		values:      values,
		barCount:    barCount,
		minStrength: minStrength,
	}
}

func (irr *IsRisingRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(irr.values) || irr.barCount <= 0 {
		return false, nil
	}

	count := 0

	for i := common.MaxInt(0, index-irr.barCount+1); i <= index; i++ {
		if i > 0 && irr.values[i] > irr.values[i-1] {
			count++
		}
	}

	ratio := float64(count) / float64(irr.barCount)

	satisfied := ratio >= irr.minStrength

	return satisfied, nil
}