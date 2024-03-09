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
	if lastIndex < 26 || lastIndex > len(ccur.High) {
		return false
	}

	highAtMinus26 := ccur.High[lastIndex-26+1]
	chikouValueAtMinus26 := ccur.Chikou[lastIndex-26]

	if chikouValueAtMinus26 <= highAtMinus26 {
		return false
	}

	for i := 1; i <= ccur.nBars; i++ {
		index := lastIndex - 26 - i
		if index < 0 {
			break
		}

		chikouValueAtIndex := ccur.Chikou[index]
		lowAtMinus26MinusI := ccur.Low[index-i]

		if chikouValueAtIndex >= lowAtMinus26MinusI {
			return false
		}
	}

	return true
}
