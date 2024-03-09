package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type ChikouAboveRule struct {
	Chikou    []float64
	Value     []float64
	traceRule bool
}

func NewChikouAboveRule(chikou, value []float64, traceRule bool) Rule {
	return ChikouCrossUpRule{
		Chikou:    chikou,
		Value:     value,
		traceRule: traceRule,
	}
}

func (car ChikouAboveRule) IsSatisfied(lastIndex int) bool {
	if lastIndex < 26 || lastIndex >= len(car.Value) {
		return false
	}

	valueAtTimeDelay := car.Value[lastIndex-25]
	chikouAtTimeDelay := car.Chikou[lastIndex-25]

	valueAtTimeDelayMinus1 := car.Value[lastIndex-25-1]
	chikouAtTimeDelayMinus1 := car.Chikou[lastIndex-25-1]

	satisfied := chikouAtTimeDelay > valueAtTimeDelay &&
		chikouAtTimeDelayMinus1 < valueAtTimeDelayMinus1

	if satisfied && car.traceRule {
		logger.Log.Info().
			Int("lastIndex", lastIndex).
			Float64("valueAtTimeDelay", valueAtTimeDelay).
			Float64("chikouAtTimeDelay", chikouAtTimeDelay).
			Float64("valueAtTimeDelayMinus1", valueAtTimeDelayMinus1).
			Float64("chikouAtTimeDelayMinus1", chikouAtTimeDelayMinus1).
			Msg("ChikouCrossUpRule satisfied")
	}

	return satisfied
}
