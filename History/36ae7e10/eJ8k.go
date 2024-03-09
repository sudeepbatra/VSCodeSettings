package rules

type ChikouCrossDownLowHighNBarsRule struct {
	Chikou []float64
	High   []float64
	Low    []float64
	nBars  int
}

func NewChikouCrossDownLowHighNBarsRule(chikou, high, low []float64, nBars int) Rule {
	return ChikouCrossDownLowHighNBarsRule{
		Chikou: chikou,
		High:   high,
		Low:    low,
		nBars:  nBars,
	}
}

func (ccdr ChikouCrossDownLowHighNBarsRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccdr.High) {
		return false
	}

	lowAtTimeDelay := ccdr.Low[lastIndex-25]
	chikouAtTimeDelay := ccdr.Chikou[lastIndex-25]

	if chikouAtTimeDelay >= lowAtTimeDelay {
		return false
	}

	for i := 1; i <= ccdr.nBars; i++ {
		index := lastIndex - 25 - i
		if index < 0 {
			break
		}

		chikouAtTimeDelayMinusI := ccdr.Chikou[index]
		highAtTimeDelayMinusI := ccdr.High[index]

		if chikouAtTimeDelayMinusI <= highAtTimeDelayMinusI {
			return false
		}
	}

	return true
}
