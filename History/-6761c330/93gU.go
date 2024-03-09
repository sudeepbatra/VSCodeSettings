package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXBasicLongStrategy struct {
	*BaseStrategy
}

func NewADXBasicLongStrategy(adx, plusDI, minusDI []float64, adxThreshold float64) Strategy {
	plusDICrossoverRule := rules.NewCrossUpRule(plusDI, minusDI)

	minusDICrossoverRule := rules.NewCrossDownRule(plusDI, minusDI)

	entryRule := rules.NewChikouCrossUpRule(chikou, value)
	exitRule := rules.NewChikouCrossDownRule(chikou, value)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &IchimokuChikouCrossoverStrategy{baseStrategy}
}

func (s *ADXBasicLongStrategy) GetName() string {
	return "ADXBasicLongStrategy"
}

func (s *ADXBasicLongStrategy) GetDescription() string {
	return "ADX Basic Long Strategy"
}
