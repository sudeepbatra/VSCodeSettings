package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type OverThresholdIndicatorRule struct {
	first     []float64
	threshold float64
	traceRule bool
}

func NewOverThresholdIndicatorRule(first []float64, threshold float64, traceRule bool) Rule {
	return &OverThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
		traceRule: traceRule,
	}
}

func (oir *OverThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) {
		return false
	}

	satisfied := oir.first[index] > oir.threshold

	if satisfied {
		if oir.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("oir.first[index]", oir.first[index]).
				Float64("threshold", oir.threshold).
				Msg("over threshold indicator rule satisfied")
		}
	}

	return satisfied
}
