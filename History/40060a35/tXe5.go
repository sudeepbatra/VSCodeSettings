package strategies

import (
	"github.com/sudeepbatra/alpha-hft/ta"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyName        = "RSI Above 70"
	StrategyDescription = "Buy when RSI above 70"
)

func NewRSIAboveStrategy(rsiSeries []float64, threshold float64) ta.Strategy {

	ta.NewStrategy(StrategyName, StrategyDescription, ta.StrategyTypeLong, rules.NewCrossAboveRule(rsiSeries, threshold), rules.NewCrossAboveRule(rsiSeries, threshold))
}
