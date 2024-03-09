package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type UnderThresholdIndicatorRule struct {
	first     []float64
	threshold float64
}

func NewUnderThresholdIndicatorRule(first []float64, threshold float64) Rule {
	return &UnderThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
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
				Float64("oir.first[index]", oir.first[index]).
				Float64("threshold", oir.threshold).
				Msg("over threshold indicator rule satisfied")
		}
	}

	return satisfied
}