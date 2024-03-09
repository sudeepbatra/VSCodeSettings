package rules

type OverIndicatorRule struct {
	first     []float64
	second    []float64
	traceRule bool
}

func NewOverIndicatorRule(first, second []float64, traceRule bool) Rule {
	return &OverIndicatorRule{
		first:     first,
		second:    second,
		traceRule: traceRule,
	}
}

func (oir *OverIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) || index >= len(oir.second) {
		return false
	}

	satisfied := oir.first[index] > oir.second[index]
	if satisfied {

	}

	return satisfied
}
