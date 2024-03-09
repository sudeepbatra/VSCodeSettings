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

func (ccdr ChikouCrossDownRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccdr.Value) {
		return false
	}

	valueAtTimeDelay := ccdr.Value[lastIndex-25]
	chikouValueAtTimeDelay := ccdr.Chikou[lastIndex-25]

	valueAtTimeDelayMinus1 := ccdr.Value[lastIndex-25-1]
	chikouValueAtMinus26Minus1 := ccdr.Chikou[lastIndex-25-1]

	return chikouValueAtTimeDelay < valueAtTimeDelay &&
		chikouValueAtMinus26Minus1 > valueAtTimeDelayMinus1
}
