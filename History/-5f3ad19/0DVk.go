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

	satisfied := chikouValueAtTimeDelay < valueAtTimeDelay
	if satisfied && cbr.traceRule {
		logger.Log.Info().
			Int("lastIndex", lastIndex).
			Float64("valueAtTimeDelay", valueAtTimeDelay).
			Float64("chikouValueAtTimeDelay", chikouValueAtTimeDelay).
			Msg("ChikouBelowRule satisfied")
	}

	return satisfied
}
