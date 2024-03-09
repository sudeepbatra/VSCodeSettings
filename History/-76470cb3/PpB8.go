package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	PSARStrategyName        = "PSAR"
	PSARStrategyDescription = "PSAR"
)

func NewPSARStrategy(close, psar []float64) Strategy {
	rsiAboveStrategy, err := NewBaseStrategy(
		PSARStrategyName,
		PSARStrategyDescription,
		StrategyTypeLong,
		rules.NewCrossUpRule(close, psar),
		rules.NewCrossDownRule(close, psar))
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return rsiAboveStrategy
}
