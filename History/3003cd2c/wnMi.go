package rules

import "github.com/sudeepbatra/alpha-hft/logger"

type OverIndicatorRule struct {
	first     []float64
	second    []float64
	traceRule bool
}

func NewOverIndicatorRule(first, second []float64, traceRule bool) Rule {
	return &OverIndicatorRule{
		first:     first,
		second:    second,
		traceRule: traceRule,
	}
}

func (oir *OverIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) || index >= len(oir.second) {
		return false
	}

	satisfied := oir.first[index] > oir.second[index]
	if satisfied {
		logger.Log.Info().
			Int("index", index).
			Float64("oir.first[index]", oir.first[index]).
			Float64("oir.second[index]", oir.second[index]).
			Msg("over indicator rule satisfied")
	}

	return satisfied
}
