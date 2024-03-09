package rules

type ChikouCrossUpLowHighRule struct {
	Chikou []float64
	High   []float64
	Low    []float64
}

func NewChikouCrossUpLowHighRule(chikou, high, low []float64) Rule {
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

	highAtTimeDelay := cculhr.High[lastIndex-26]
	chikouValueAtTimeDelay := cculhr.Chikou[lastIndex-26]

	lowAtMinus26Minus1 := cculhr.Low[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cculhr.Chikou[lastIndex-26-1]

	return chikouValueAtTimeDelay > highAtTimeDelay &&
		chikouValueAtMinus26Minus1 < lowAtMinus26Minus1
}
