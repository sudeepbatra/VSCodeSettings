package rules

type ChikouCrossUpRule struct {
	Chikou    []float64
	Value     []float64
	traceRule bool
}

func NewChikouCrossUpRule(chikou, value []float64, traceRule bool) Rule {
	return ChikouCrossUpRule{
		Chikou:    chikou,
		Value:     value,
		traceRule: traceRule,
	}
}

func (ccur ChikouCrossUpRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(ccur.Value) {
		return false
	}

	valueAtTimeDelay := ccur.Value[lastIndex-25]
	chikouAtTimeDelay := ccur.Chikou[lastIndex-25]

	valueAtTimeDelayMinus1 := ccur.Value[lastIndex-25-1]
	chikouAtTimeDelayMinus1 := ccur.Chikou[lastIndex-25-1]

	satisfied := chikouAtTimeDelay > valueAtTimeDelay &&
		chikouAtTimeDelayMinus1 < valueAtTimeDelayMinus1
}
