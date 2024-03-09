package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuChikouHighLowCrossoverStrategy struct {
	*BaseStrategy
}

func NewIchimokuChikouHighLowCrossoverStrategy(chikou, high, low []float64) Strategy {
	entryRule := rules.NewChikouCrossUpLowHighRule(chikou, high, low)
	exitRule := rules.NewChikouCrossDownLowHighRule(chikou, high, low)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &IchimokuChikouHighLowCrossoverStrategy{baseStrategy}
}

func (s *IchimokuChikouHighLowCrossoverStrategy) GetName() string {
	return "IchimokuChikouCrossoverStrategy"
}

func (s *IchimokuChikouHighLowCrossoverStrategy) GetDescription() string {
	return "Ichimoku Chikou Crossover Strategy"
}
