package rules

type ChikouCrossoverRule struct {
	Chikou []float64
	value  []float64
}

func NewChikouCloseCrossoverRule(chikou, close []float64) Rule {
	return ChikouCrossoverRule{
		Chikou: chikou,
		Close:  close,
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
