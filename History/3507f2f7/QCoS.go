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

	highAtTimeDelay := ccur.High[lastIndex-25]
	chikouAtTimeDelay := ccur.Chikou[lastIndex-25]

	if chikouAtTimeDelay <= highAtTimeDelay {
		return false
	}

	for i := 1; i <= ccur.nBars; i++ {
		index := lastIndex - 25 - i
		if index < 0 {
			break
		}

		chikouAtTimeDelayMinusI := ccur.Chikou[index]
		lowAtTimeDelayMinusI := ccur.Low[index]

		if chikouAtTimeDelayMinusI >= lowAtTimeDelayMinusI {
			return false
		}
	}

	return true
}
