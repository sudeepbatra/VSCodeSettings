package rules

type ChikouCrossUpLowHighRule struct {
	Chikou []float64
	High   []float64
	Low    []float64
}

func NewChikouCrossUpLowHighRule(chikou, high, low []float64, nBars int) Rule {
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

	highAtMinus26 := cculhr.High[lastIndex-26]
	chikouValueAtMinus26 := cculhr.Chikou[lastIndex-26]

	lowAtMinus26Minus1 := cculhr.Low[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cculhr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > highAtMinus26 &&
		chikouValueAtMinus26Minus1 < lowAtMinus26Minus1
}
