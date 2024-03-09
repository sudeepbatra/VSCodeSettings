package rules

type ChikouCrossUpRule struct {
	Chikou []float64
	Value  []float64
}

func NewChikouCrossoverRule(chikou, value []float64) Rule {
	return ChikouCrossUpRule{
		Chikou: chikou,
		Value:  value,
	}
}

func (ccur ChikouCrossUpRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccur.Value) {
		return false
	}

	closeValueAtMinus26 := ccur.Value[lastIndex-26]
	chikouValueAtMinus26 := ccur.Chikou[lastIndex-26]

	closeValueAtMinus26Minus1 := ccur.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := ccur.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > closeValueAtMinus26 &&
		chikouValueAtMinus26Minus1 < closeValueAtMinus26Minus1
}
