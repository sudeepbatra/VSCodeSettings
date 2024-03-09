package rules

type crossAboveThresholdRule struct {
	series    []float64
	threshold float64
	traceRule bool
}

func NewCrossAboveRule(series []float64, threshold float64, traceRule bool) Rule {
	return crossAboveThresholdRule{
		series:    series,
		threshold: threshold,
		traceRule: traceRule,
	}
}

func (cr crossAboveThresholdRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.series) {
		return false
	}

	if cr.series[index-1] < cr.threshold && cr.series[index] >= cr.threshold {
		return true
	}

	return false
}
