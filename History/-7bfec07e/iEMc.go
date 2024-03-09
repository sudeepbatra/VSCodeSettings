package common

func NewCrossUpRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
	}
}

func NewCrossDownRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
	}
}

type Rule interface {
	IsSatisfied(index int) bool
}

type crossRule struct {
	upper []float64
	lower []float64
}

func (cr crossRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	// Check for cross-up (lower crosses above upper)
	if cr.lower[index-1] < cr.upper[index-1] && cr.lower[index] >= cr.upper[index] {
		return true
	}

	// Check for cross-down (upper crosses below lower)
	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		return true
	}

	return false
}
