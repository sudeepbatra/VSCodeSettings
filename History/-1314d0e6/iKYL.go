package rules

type ChikouCrossoverRule struct {
	Chikou []float64
	Value  []float64
}

func NewChikouCloseCrossoverRule(chikou, value []float64) Rule {
	return ChikouCrossoverRule{
		Chikou: chikou,
		Value:  value,
	}
}

func (cccr ChikouCrossoverRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(cccr.Value) {
		return false
	}

	latestChikou := cccr.Chikou[lastIndex]
	latestClose := cccr.Value[lastIndex]

	// Check if the Chikou Span crosses above the latest close price
	if latestChikou > latestClose && cccr.Chikou[lastIndex-1] <= cccr.Value[lastIndex-1] {
		return true
	}

	return false
}
