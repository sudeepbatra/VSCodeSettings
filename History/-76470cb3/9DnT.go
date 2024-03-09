package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type PSARStrategy struct {
	*BaseStrategy
}

const (
	PSARStrategyName        = "PSAR"
	PSARStrategyDescription = "PSAR"
)

func NewPSARStrategy(close, psar []float64) Strategy {
	entryRule := rules.NewCrossUpRule(close, psar)
	exitRule := rules.NewCrossDownRule(close, psar)

	psarStrategy, err := NewBaseStrategy(
		StrategyTypeLong,
		rules.NewCrossUpRule(close, psar),
		rules.NewCrossDownRule(close, psar))
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return rsiAboveStrategy
}
