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

	closeValueAtMinus26 := cccr.Value[lastIndex-26]
	chikouValueAtMinus26 := cccr.Chikou[lastIndex-26]

	closeValueAtMinus26Minus1 := cccr.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cccr.Chikou[lastIndex-26]

	// Check if the Chikou Span at index - 26 is above the Close at index - 26
	// and the Chikou Span at index - 26 - 1 is below the Close at index - 26 - 1
	if chikouValueAtMinus26 > closeValueAtMinus26 &&
		chikouValueAtMinus26Minus1 < closeValueAtMinus26Minus1 {
		return true
	}

	return false
}
