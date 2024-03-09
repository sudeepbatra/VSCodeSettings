package common

type crossRule struct {
	upper []float64
	lower []float64
}

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

func (cr crossRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.lower[index-1] < cr.upper[index-1] && cr.lower[index] >= cr.upper[index] {
		return true
	}

	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		return true
	}

	return false
}
