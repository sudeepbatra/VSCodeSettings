package common

type IsFallingRule struct {
	ref         []float64
	barCount    int
	minStrength float64
}

func NewIsFallingRule(ref []float64, barCount int) Rule {
	return &IsFallingRule{
		ref:         ref,
		barCount:    barCount,
		minStrength: 1.0,
	}
}

func NewIsFallingRuleWithStrength(ref []float64, barCount int, minStrength float64) Rule {
	if minStrength >= 1 {
		minStrength = 0.99
	}

	return &IsFallingRule{
		ref:         ref,
		barCount:    barCount,
		minStrength: minStrength,
	}
}

func (ifr *IsFallingRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(ifr.ref) || ifr.barCount <= 0 {
		return false, nil
	}

	count := 0

	for i := max(0, index-ifr.barCount+1); i <= index; i++ {
		if i > 0 && ifr.ref[i] < ifr.ref[i-1] {
			count++
		}
	}

	ratio := float64(count) / float64(ifr.barCount)

	satisfied := ratio >= ifr.minStrength

	return satisfied, nil
}
