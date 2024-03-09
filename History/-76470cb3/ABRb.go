package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type PSARStrategy struct {
	*BaseStrategy
}

func NewPSARStrategy(close, psar []float64) Strategy {
	entryRule := rules.NewCrossUpRule(close, psar)
	exitRule := rules.NewCrossDownRule(close, psar)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &PSARStrategy{baseStrategy}
}

func (s *PSARStrategy) GetName() string {
	return "PSAR Strategy"
}

func (s *PSARStrategy) GetDescription() string {
	return "PSAR Strategy"
}
