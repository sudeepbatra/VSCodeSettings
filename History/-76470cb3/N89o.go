package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	PSARStrategyName        = "PSAR"
	PSARStrategyDescription = "PSAR"
)

func NewPSARStrategy(rsiSeries []float64) Strategy {
	rsiAboveStrategy, err := rules.NewCrossAboveRule(
		PSARStrategyName,
		PSARStrategyDescription,
		StrategyTypeLong,
		rules.NewCrossAboveRule(rsiSeries, LongThreshold),
		rules.NewCrossBelowRule(rsiSeries, ExitThreshold))
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return rsiAboveStrategy
}
