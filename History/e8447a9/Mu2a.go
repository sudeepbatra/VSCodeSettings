package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuChikouCrossoverStrategy struct {
	*BaseStrategy
}

func NewIchimokuChikouCrossoverStrategy(chikou, value []float64) Strategy {
	entryRule := rules.NewChikouCrossUpRule(chikou, value)
	exitRule := rules.NewChikouCrossDownRule(chikou, value)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &IchimokuChikouCrossoverStrategy{baseStrategy}
}

func (s *IchimokuChikouCrossoverStrategy) GetName() string {
	return "IchimokuChikouCrossoverStrategy"
}

func (s *IchimokuChikouCrossoverStrategy) GetDescription() string {
	return "Ichimoku Chikou Crossover Strategy"
}
