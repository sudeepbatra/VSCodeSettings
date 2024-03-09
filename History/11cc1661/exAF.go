package rules

type ChikouCrossDownLowHighRule struct {
	Chikou    []float64
	High      []float64
	Low       []float64
	traceRule bool
}

func NewChikouCrossDownLowHighRule(chikou, high, low []float64, traceRule bool) Rule {
	return ChikouCrossDownLowHighRule{
		Chikou:    chikou,
		High:      high,
		Low:       low,
		traceRule: traceRule,
	}
}

func (ccdlhr ChikouCrossDownLowHighRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccdlhr.High) {
		return false
	}

	lowAtTimeDelay := ccdlhr.Low[lastIndex-25]
	chikouValueAtTimeDelay := ccdlhr.Chikou[lastIndex-25]

	highAtTimeDelayMinus1 := ccdlhr.High[lastIndex-25-1]
	chikouValueAtTimeDelayMinus1 := ccdlhr.Chikou[lastIndex-25-1]

	return chikouValueAtTimeDelay < lowAtTimeDelay &&
		chikouValueAtTimeDelayMinus1 > highAtTimeDelayMinus1
}
