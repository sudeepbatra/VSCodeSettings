package rules

type crossUpRule struct {
	upper     []float64
	lower     []float64
	traceRule bool
}

func NewCrossUpRule(upper, lower []float64, traceRule bool) Rule {
	return crossUpRule{
		upper:     upper,
		lower:     lower,
		traceRule: traceRule,
	}
}

func (cr crossUpRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.lower[index-1] < cr.upper[index-1] && cr.lower[index] >= cr.upper[index] {
		return true
	}

	return false
}
