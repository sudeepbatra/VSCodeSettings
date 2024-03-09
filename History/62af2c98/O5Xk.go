package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type crossAboveThresholdRule struct {
	series    []float64
	threshold float64
	traceRule bool
}

func NewCrossAboveRule(series []float64, threshold float64, traceRule bool) Rule {
	return crossAboveThresholdRule{
		series:    series,
		threshold: threshold,
		traceRule: traceRule,
	}
}

func (cr crossAboveThresholdRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.series) {
		return false
	}

	if cr.series[index-1] < cr.threshold && cr.series[index] >= cr.threshold {
		if cr.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("cr.series[index-1]", cr.series[index-1]).
				Float64("cr.threshold", cr.threshold).
				Float64("cr.series[index]", cr.series[index]).
				Float64("cr.threshold", cr.threshold).
				Msg("cross above threshold rule satisfied")
		}

		return true
	}

	return false
}
