package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuChikouHighLowNBarsCrossoverStrategy struct {
	*BaseStrategy
}

func NewIchimokuChikouHighLowNBarsCrossoverStrategy(chikou, high, low []float64) Strategy {
	entryRule := rules.NewChikouCrossUpLowHighNBarsRule(chikou, high, low, 10)
	exitRule := rules.NewChikouCrossDownLowHighNBarsRule(chikou, high, low, 10)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &IchimokuChikouHighLowNBarsCrossoverStrategy{baseStrategy}
}

func (s *IchimokuChikouHighLowNBarsCrossoverStrategy) GetName() string {
	return "IchimokuChikouHighLowNBarsCrossoverStrategy"
}

func (s *IchimokuChikouHighLowNBarsCrossoverStrategy) GetDescription() string {
	return "Ichimoku Chikou High Low NBars Crossover Strategy"
}
