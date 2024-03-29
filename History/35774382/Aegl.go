package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type crossUpRule struct {
	a         []float64
	b         []float64
	traceRule bool
}

func NewCrossUpRule(a, b []float64, traceRule bool) Rule {
	return crossUpRule{
		a:         a,
		b:         b,
		traceRule: traceRule,
	}
}

func (cr crossUpRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.a) {
		return false
	}

	if cr.b[index-1] < cr.a[index-1] && cr.b[index] >= cr.a[index] {
		if cr.traceRule {
			logger.Log.Info().
				Int("index", index).
				Float64("cr.lower[index-1]", cr.b[index-1]).
				Float64("cr.upper[index-1]", cr.a[index-1]).
				Float64("cr.lower[index]", cr.b[index]).
				Float64("cr.upper[index]", cr.a[index]).
				Msg("cross up rule satisfied")
		}

		return true
	}

	return false
}
