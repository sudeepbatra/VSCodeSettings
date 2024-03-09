package ta

import (
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	PSARStrategyName        = "PSAR"
	PSARStrategyDescription = "PSAR"
)

func NewPSARStrategy(close, psar []float64) Strategy {
	rsiAboveStrategy := NewBaseStrategy(
		rules.NewCrossUpRule(close, psar),
		rules.NewCrossDownRule(close, psar))

	return rsiAboveStrategy
}
