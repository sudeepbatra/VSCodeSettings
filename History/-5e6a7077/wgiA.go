package rules

type crossDownRule struct {
	upper []float64
	lower []float64
}

func NewCrossDownRule(upper, lower []float64) Rule {
	return crossDownRule{
		upper: upper,
		lower: lower,
	}
}

func (cr crossDownRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		return true
	}

	return false
}