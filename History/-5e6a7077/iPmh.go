package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type crossDownRule struct {
	a         []float64
	b         []float64
	traceRule bool
}

func NewCrossDownRule(upper, lower []float64, traceRule bool) Rule {
	return crossDownRule{
		a:         upper,
		b:         lower,
		traceRule: traceRule,
	}
}

func (cr crossDownRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.a) {
		return false
	}

	if cr.a[index-1] > cr.b[index-1] && cr.a[index] < cr.b[index] {
		if cr.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("cr.b[index-1]", cr.b[index-1]).
				Float64("cr.a[index-1]", cr.a[index-1]).
				Float64("cr.lower[index]", cr.b[index]).
				Float64("cr.upper[index]", cr.a[index]).
				Msg("cross down rule satisfied")
		}

		return true
	}

	return false
}
