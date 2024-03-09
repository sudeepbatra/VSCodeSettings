package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXBasicLongStrategy struct {
	*BaseStrategy
}

func NewADXBasicLongStrategy(adx, plusDI, minusDI []float64, adxThreshold float64) Strategy {
	plusDICrossAboveMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI)
	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold)

	adxThresholdCrossoverRule := rules.NewCrossAboveRule(adx, adxThreshold)
	plusDIAboveMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI)

	plusDICrossoverEntryRule := rules.And(plusDICrossAboveMinusDIRule, adxOverThresholdRule)
	adxThresholdCrossoverEntryRule := rules.And(adxThresholdCrossoverRule, plusDIAboveMinusDIRule)

	entryRule := rules.Or(plusDICrossoverEntryRule, adxThresholdCrossoverEntryRule)

	plusDICrossDownMinusDIRule := rules.NewCrossDownRule(plusDI, minusDI)

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
