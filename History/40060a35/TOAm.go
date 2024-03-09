package strategies

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyName        = "RSI Above 70"
	StrategyDescription = "Buy when RSI above 70"
	LongThreshold       = 70
	ExitThreshold       = 30
)

func NewRSIAboveStrategy(rsiSeries []float64, threshold float64) ta.Strategy {
	rsiAboveStrategy, err := ta.NewStrategy(StrategyName, StrategyDescription, ta.StrategyTypeLong, rules.NewCrossAboveRule(rsiSeries, LongThreshold), rules.NewCrossAboveRule(rsiSeries, ExitThreshold))

	if err != nil {
		logger.Log.Error().Err(err).Msg("Error scanning row")
	}

	return rsiAboveStrategy
}
