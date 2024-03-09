package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type ChikouBelowRule struct {
	Chikou    []float64
	Value     []float64
	traceRule bool
}

func NewChikouBelowRule(chikou, value []float64, traceRule bool) Rule {
	return ChikouCrossDownRule{
		Chikou:    chikou,
		Value:     value,
		traceRule: traceRule,
	}
}

func (cbr ChikouBelowRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(cbr.Value) {
		return false
	}

	valueAtTimeDelay := cbr.Value[lastIndex-25]
	chikouValueAtTimeDelay := cbr.Chikou[lastIndex-25]

	valueAtTimeDelayMinus1 := cbr.Value[lastIndex-25-1]
	chikouValueAtTimeDelayMinus1 := cbr.Chikou[lastIndex-25-1]

	satisfied := chikouValueAtTimeDelay < valueAtTimeDelay &&
		chikouValueAtTimeDelayMinus1 > valueAtTimeDelayMinus1

	if satisfied && cbr.traceRule {
		logger.Log.Info().
			Int("lastIndex", lastIndex).
			Float64("valueAtTimeDelay", valueAtTimeDelay).
			Float64("chikouValueAtTimeDelay", chikouValueAtTimeDelay).
			Float64("valueAtTimeDelayMinus1", valueAtTimeDelayMinus1).
			Float64("chikouValueAtTimeDelayMinus1", chikouValueAtTimeDelayMinus1).
			Msg("ChikouCrossDownRule satisfied")
	}

	return satisfied
}
