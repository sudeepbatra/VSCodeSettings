package common

func NewCrossUpRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
		cmp:   1,
	}
}

func NewCrossDownRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
		cmp:   -1,
	}
}

type Rule interface {
	IsSatisfied(index int) bool
}

type crossRule struct {
	upper []float64
	lower []float64
	cmp   int
}

func (cr crossRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		return true
	} else if cr.upper[index-1] > cr.lower[index-1] && cr.upper[index] <= cr.lower[index] {
		return true
	}

	return false
}
