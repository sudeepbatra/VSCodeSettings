package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type UnderIndicatorRule struct {
	a         []float64
	b         []float64
	traceRule bool
}

func NewUnderIndicatorRule(a, b []float64) Rule {
	return &UnderIndicatorRule{
		a: a,
		b: b,
	}
}

func (uir *UnderIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(uir.a) || index >= len(uir.b) {
		return false
	}

	satisfied := uir.a[index] < uir.b[index]

	if satisfied {
		logger.Log.Info().
			Int("index", index).
			Float64("uir.a[index]", uir.a[index]).
			Float64("uir.b[index]", uir.b[index]).
			Msg("under indicator rule satisfied")
	}

	return satisfied
}
