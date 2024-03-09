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

func (cccr ChikouCrossoverRule) IsSatisfied(index int) bool {
	if index < 26 || index >= len(cccr.Value) {
		return false
	}

	latestChikou := cccr.Chikou[index]
	latestClose := cccr.Value[index]

	// Check if the Chikou Span crosses above the latest close price
	if latestChikou > latestClose && cccr.Chikou[index-1] <= cccr.Close[index-1] {
		return true
	}

	return false
}
