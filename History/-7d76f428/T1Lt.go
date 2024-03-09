package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type OverThresholdIndicatorRule struct {
	first     []float64
	threshold float64
	traceRule bool
}

func NewOverThresholdIndicatorRule(first []float64, threshold float64) Rule {
	return &OverThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
	}
}

func (oir *OverThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) {
		return false
	}

	satisfied := oir.first[index] > oir.threshold

	if satisfied {
		if cr.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("cr.lower[index-1]", cr.lower[index-1]).
				Float64("cr.upper[index-1]", cr.upper[index-1]).
				Float64("cr.lower[index]", cr.lower[index]).
				Float64("cr.upper[index]", cr.upper[index]).
				Msg("cross up rule satisfied")
		}
	}

	return satisfied
}
