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

	lowAtMinus26 := cculhr.Low[lastIndex-26]
	chikouValueAtMinus26 := cculhr.Chikou[lastIndex-26]

	highAtMinus26Minus1 := cculhr.High[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cculhr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > lowAtMinus26 &&
		chikouValueAtMinus26Minus1 < highAtMinus26Minus1
}
