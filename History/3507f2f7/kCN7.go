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

	if chikouValueAtMinus26 <= highAtMinus26 {
		return false
	}

	// Check the condition for the previous bars
	for i := 1; i <= ccur.nBars; i++ {
		index := lastIndex - 26 - i
		if index < 0 {
			// If we reach the beginning of the data, exit the loop.
			break
		}

		// Check the condition for the previous bars
		chikouValueAtIndex := ccur.Chikou[index]
		closeAtMinus26MinusI := ccur.Low[index-i]

		if chikouValueAtIndex >= closeAtMinus26MinusI {
			// If the condition is not met for any of the previous bars, return false.
			return false
		}
	}

	// If the condition is met for all previous bars, return true.
	return true
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