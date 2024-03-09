package rules

type ChikouCrossUpRule struct {
	Chikou []float64
	Value  []float64
}

func NewChikouCrossUpRule(chikou, value []float64) Rule {
	return ChikouCrossUpRule{
		Chikou: chikou,
		Value:  value,
	}
}

func (ccur ChikouCrossUpRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccur.Value) {
		return false
	}

	valueAtTimeDelay := ccur.Value[lastIndex-25]
	chikouAtTimeDelay := ccur.Chikou[lastIndex-25]

	valueAtMinus26Minus1 := ccur.Value[lastIndex-25-1]
	chikouValueAtMinus26Minus1 := ccur.Chikou[lastIndex-25-1]

	return chikouAtTimeDelay > valueAtTimeDelay &&
		chikouValueAtMinus26Minus1 < valueAtMinus26Minus1
}
