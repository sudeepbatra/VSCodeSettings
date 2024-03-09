package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type OverIndicatorRule struct {
	a         []float64
	b         []float64
	traceRule bool
}

func NewOverIndicatorRule(first, second []float64, traceRule bool) Rule {
	return &OverIndicatorRule{
		a:         first,
		b:         second,
		traceRule: traceRule,
	}
}

func (oir *OverIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.a) || index >= len(oir.b) {
		return false
	}

	satisfied := oir.a[index] > oir.b[index]
	if satisfied {
		logger.Log.Info().
			Int("index", index).
			Float64("oir.first[index]", oir.a[index]).
			Float64("oir.second[index]", oir.b[index]).
			Msg("over indicator rule satisfied")
	}

	return satisfied
}
