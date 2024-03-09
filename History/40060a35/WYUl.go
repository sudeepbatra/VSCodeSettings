package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyName        = "RSI Above 70"
	StrategyDescription = "Buy when RSI above 70"
	LongThreshold       = 50
	ExitThreshold       = 30
)

type RSIAboveStrategy struct {
	*BaseStrategy
}

func NewRSIAboveStrategy(rsiSeries []float64) Strategy {
	entryRule := rules.NewCrossAboveRule(rsiSeries, LongThreshold)
	exitRule := rules.NewCrossBelowRule(rsiSeries, ExitThreshold)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &RSIAboveStrategy{baseStrategy}
}

func (s *RSIAboveStrategy) GetName() string {
	return StrategyName
}

func (s *RSIAboveStrategy) GetDescription() string {
	return StrategyDescription
}
