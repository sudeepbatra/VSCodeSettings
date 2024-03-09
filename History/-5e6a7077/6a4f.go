package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type crossDownRule struct {
	upper     []float64
	lower     []float64
	traceRule bool
}

func NewCrossDownRule(upper, lower []float64, traceRule bool) Rule {
	return crossDownRule{
		upper:     upper,
		lower:     lower,
		traceRule: traceRule,
	}
}

func (cr crossDownRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		if cr.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("cr.lower[index-1]", cr.lower[index-1]).
				Float64("cr.upper[index-1]", cr.upper[index-1]).
				Float64("cr.lower[index]", cr.lower[index]).
				Float64("cr.upper[index]", cr.upper[index]).
				Msg("cross down rule satisfied")
		}

		return true
	}

	return false
}
