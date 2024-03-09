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

	lowAtMinus26 := ccdr.Low[lastIndex-25]
	chikouValueAtMinus26 := ccdr.Chikou[lastIndex-25]

	if chikouValueAtMinus26 >= lowAtMinus26 {
		return false
	}

	for i := 1; i <= ccdr.nBars; i++ {
		index := lastIndex - 25 - i
		if index < 0 {
			break
		}

		chikouValueAtIndex := ccdr.Chikou[index]
		highAtMinus26MinusI := ccdr.High[index-i]

		if chikouValueAtIndex <= highAtMinus26MinusI {
			return false
		}
	}

	return true
}
