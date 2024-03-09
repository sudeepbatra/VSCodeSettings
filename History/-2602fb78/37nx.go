package rules

type ChikouCrossUpLowHighRule struct {
	Chikou []float64
	Low    []float64
	High   []float64
}

func ChikouCrossUpLowHighRule(chikou, high, low []float64) Rule {
	return ChikouCrossUpLowHighRule{
		Chikou: chikou,
		High:   high,
		Low:    low,
	}
}

func (cculhr ChikouCrossUpLowHighRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(cculhr.High) {
		return false
	}

	valueAtMinus26 := cculhr.Value[lastIndex-26]
	chikouValueAtMinus26 := cculhr.Chikou[lastIndex-26]

	valueAtMinus26Minus1 := cculhr.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cculhr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > valueAtMinus26 &&
		chikouValueAtMinus26Minus1 < valueAtMinus26Minus1
}
