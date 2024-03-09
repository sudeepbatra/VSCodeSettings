package common

type IsFallingRule struct {
	values      []float64
	barCount    int
	minStrength float64
}

func NewIsFallingRule(values []float64, barCount int) Rule {
	return &IsFallingRule{
		values:      values,
		barCount:    barCount,
		minStrength: 1.0,
	}
}

func NewIsFallingRuleWithStrength(values []float64, barCount int, minStrength float64) Rule {
	if minStrength >= 1 {
		minStrength = 0.99
	}

	return &IsFallingRule{
		values:      values,
		barCount:    barCount,
		minStrength: minStrength,
	}
}

func (ifr *IsFallingRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(ifr.values) || ifr.barCount <= 0 {
		return false, nil
	}

	count := 0

	for i := maxInt(0, index-ifr.barCount+1); i <= index; i++ {
		if i > 0 && ifr.values[i] < ifr.values[i-1] {
			count++
		}
	}

	ratio := float64(count) / float64(ifr.barCount)

	satisfied := ratio >= ifr.minStrength

	return satisfied, nil
}
