package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyName        = "PSAR Strategy"
	StrategyDescription = "PSAR Strategy"
	LongThreshold       = 50
	ExitThreshold       = 30
)

func NewRSIAboveStrategy(rsiSeries []float64) Strategy {
	rsiAboveStrategy, err := NewStrategy(
		StrategyName,
		StrategyDescription,
		StrategyTypeLong,
		rules.NewCrossAboveRule(rsiSeries, LongThreshold),
		rules.NewCrossBelowRule(rsiSeries, ExitThreshold))
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return rsiAboveStrategy
}
