package rules

type ChikouCrossDownLowHighRule struct {
	Chikou []float64
	High   []float64
	Low    []float64
}

func NewChikouCrossDownLowHighRule(chikou, high, low []float64) Rule {
	return ChikouCrossDownLowHighRule{
		Chikou: chikou,
		High:   high,
		Low:    low,
	}
}

func (ccdlhr ChikouCrossDownLowHighRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccdlhr.High) {
		return false
	}

	highAtMinus26 := ccdlhr.High[lastIndex-26]
	chikouValueAtMinus26 := ccdlhr.Chikou[lastIndex-26]

	lowAtMinus26Minus1 := ccdlhr.Low[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := ccdlhr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > highAtMinus26 &&
		chikouValueAtMinus26Minus1 < lowAtMinus26Minus1
}
