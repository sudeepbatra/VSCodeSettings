package rules

type ChikouCrossUpLowHighNBarsRule struct {
	Chikou []float64
	High   []float64
	Low    []float64
	nBars  int
}

func NewChikouCrossUpLowHighNBarsRule(chikou, high, low []float64, nBars int) Rule {
	return ChikouCrossUpLowHighNBarsRule{
		Chikou: chikou,
		High:   high,
		Low:    low,
		nBars:  nBars,
	}
}

func (ccur ChikouCrossUpLowHighNBarsRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccur.High) {
		return false
	}

	highAtMinus26 := ccur.High[lastIndex-26]
	chikouValueAtMinus26 := ccur.Chikou[lastIndex-26]

	lowAtMinus26Minus1 := ccur.Low[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := ccur.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > highAtMinus26 &&
		chikouValueAtMinus26Minus1 < lowAtMinus26Minus1
}
