package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuChikouCrossoverStrategy struct {
	*BaseStrategy
}

func NewIchimokuChikouCrossoverStrategy(chikou, value []float64, rsiLongThreshold, rsiExitThreshold float64) Strategy {
	rules.NewChikouCrossUpRule(chikou, value)

	entryRule := rules.NewCrossAboveRule(rsiSeries, rsiLongThreshold)
	exitRule := rules.NewCrossBelowRule(rsiSeries, rsiExitThreshold)

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
