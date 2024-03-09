package rules

type ChikouCrossDownRule struct {
	Chikou []float64
	Value  []float64
}

func NewChikouCrossDownRule(chikou, value []float64) Rule {

	return ChikouCrossDownRule{
		Chikou: chikou,
		Value:  value,
	}
}

// IsSatisfied checks if the Chikou Span crosses below the Value at the given lastIndex.
func (ccdr ChikouCrossDownRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccdr.Value) {
		return false
	}

	valueAtMinus26 := ccdr.Value[lastIndex-26]
	chikouValueAtMinus26 := ccdr.Chikou[lastIndex-26]

	valueAtMinus26Minus1 := ccdr.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := ccdr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 < valueAtMinus26 &&
		chikouValueAtMinus26Minus1 > valueAtMinus26Minus1, nil
}
