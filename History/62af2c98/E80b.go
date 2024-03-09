package rules

type crossAboveThresholdRule struct {
	series    []float64
	threshold float64
}

func NewCrossAboveRule(series []float64, threshold float64) Rule {
	return crossAboveThresholdRule{
		series:    series,
		threshold: threshold,
	}
}

func (cr crossAboveThresholdRule) IsSatisfied(index int) (bool, error) {
	if index <= 0 || index >= len(cr.series) {
		return false, nil
	}

	if cr.series[index-1] < cr.threshold && cr.series[index] >= cr.threshold {
		return true, nil
	}

	return false, nil
}
