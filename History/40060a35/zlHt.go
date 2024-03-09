package strategies

import "github.com/sudeepbatra/alpha-hft/ta"

const (
	StrategyName = "RSI Above 70"
)

func NewRSIAboveStrategy(rsiSeries []float64, threshold float64) ta.Strategy {

	ta.NewStrategy(StrategyName)
}
