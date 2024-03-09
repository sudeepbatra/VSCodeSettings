package rules

import (
	"github.com/sudeepbatra/alpha-hft/logger"
)

type UnderThresholdIndicatorRule struct {
	first     []float64
	threshold float64
	traceRule bool
}

func NewUnderThresholdIndicatorRule(first []float64, threshold float64, traceRule bool) Rule {
	return &UnderThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
		traceRule: traceRule,
	}
}

func (uir *UnderThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(uir.first) {
		return false
	}

	satisfied := uir.first[index] < uir.threshold

	if satisfied {
		if uir.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("uir.first[index]", uir.first[index]).
				Float64("threshold", uir.threshold).
				Msg("under threshold indicator rule satisfied")
		}
	}

	return satisfied
}
