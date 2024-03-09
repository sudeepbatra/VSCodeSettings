package rules

type ChikouCrossUpLowHighRule struct {
	Chikou    []float64
	High      []float64
	Low       []float64
	traceRule bool
}

func NewChikouCrossUpLowHighRule(chikou, high, low []float64, traceRule bool) Rule {
	return ChikouCrossUpLowHighRule{
		Chikou:    chikou,
		High:      high,
		Low:       low,
		traceRule: traceRule,
	}
}

func (cculhr ChikouCrossUpLowHighRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(cculhr.High) {
		return false
	}

	highAtTimeDelay := cculhr.High[lastIndex-25]
	chikouAtTimeDelay := cculhr.Chikou[lastIndex-25]

	lowAtTimeDelayMinus1 := cculhr.Low[lastIndex-25-1]
	chikouAtTimeDelayMinus1 := cculhr.Chikou[lastIndex-25-1]

	satisfied := chikouAtTimeDelay > highAtTimeDelay &&
		chikouAtTimeDelayMinus1 < lowAtTimeDelayMinus1
	return
}
