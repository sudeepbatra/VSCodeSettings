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

	highAtTimeDelay := cculhr.High[lastIndex-25]
	chikouValueAtTimeDelay := cculhr.Chikou[lastIndex-25]

	lowAtTimeDelayMinus1 := cculhr.Low[lastIndex-25-1]
	chikouValueAtTimeDelayMinus1 := cculhr.Chikou[lastIndex-25-1]

	return chikouValueAtTimeDelay > highAtTimeDelay &&
		chikouValueAtTimeDelayMinus1 < lowAtTimeDelayMinus1
}
