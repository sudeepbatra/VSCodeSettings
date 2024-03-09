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

func (cccr ChikouCrossoverRule) IsSatisfied(index int) bool {
	if index < 26 || index >= len(cccr.Close) {
		return false
	}

	latestChikou := cccr.Chikou[index]
	latestClose := cccr.Close[index]

	// Check if the Chikou Span crosses above the latest close price
	if latestChikou > latestClose && cccr.Chikou[index-1] <= cccr.Close[index-1] {
		return true
	}

	return false
}
