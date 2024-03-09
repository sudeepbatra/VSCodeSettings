package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type crossUpRule struct {
	upper     []float64
	lower     []float64
	traceRule bool
}

func NewCrossUpRule(upper, lower []float64, traceRule bool) Rule {
	return crossUpRule{
		upper:     upper,
		lower:     lower,
		traceRule: traceRule,
	}
}

func (cr crossUpRule) IsSatisfied(index int) bool {
	if index <= 0 || index >= len(cr.upper) {
		return false
	}

	if cr.lower[index-1] < cr.upper[index-1] && cr.lower[index] >= cr.upper[index] {
		if cr.traceRule {
			logger.Log.Info()
		}

		return true
	}

	return false
}
