package rules

type crossBelowThresholdRule struct {
	series    []float64
	threshold float64
}

func NewCrossBelowRule(series []float64, threshold float64) Rule {
	return crossBelowThresholdRule{
		series:    series,
		threshold: threshold,
	}
}

func (cr crossBelowThresholdRule) IsSatisfied(index int) (bool, error) {
	if index <= 0 || index >= len(cr.series) {
		return false, nil
	}

	if cr.series[index-1] >= cr.threshold && cr.series[index] < cr.threshold {
		return true, nil
	}

	return false, nil
}
