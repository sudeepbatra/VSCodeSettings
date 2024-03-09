package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	LongThreshold = 70
	ExitThreshold = 30
)

type RSIAboveStrategy struct {
	*BaseStrategy
}

func NewRSIAboveStrategy(rsiSeries []float64) Strategy {
	entryRule := rules.NewCrossAboveRule(rsiSeries, 70)
	exitRule := rules.NewCrossBelowRule(rsiSeries, ExitThreshold)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &RSIAboveStrategy{baseStrategy}
}

func (s *RSIAboveStrategy) GetName() string {
	return "RSI Above 70"
}

func (s *RSIAboveStrategy) GetDescription() string {
	return "Buy when RSI above 70"
}
