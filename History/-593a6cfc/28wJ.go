package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXIchimokuLongStrategy struct {
	*BaseStrategy
}

func NewADXIchimokuLongStrategy(adx, plusDI, minusDI []float64, adxThreshold float64) Strategy {
	plusDICrossAboveMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, true)
	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold, true)

	adxThresholdCrossoverRule := rules.NewCrossAboveThresholdRule(adx, adxThreshold, true)
	plusDIAboveMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, true)

	plusDICrossoverEntryRule := rules.And(plusDICrossAboveMinusDIRule, adxOverThresholdRule)
	adxThresholdCrossoverEntryRule := rules.And(adxThresholdCrossoverRule, plusDIAboveMinusDIRule)

	entryRule := rules.Or(plusDICrossoverEntryRule, adxThresholdCrossoverEntryRule)

	exitRule := rules.NewCrossDownRule(plusDI, minusDI, true)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &ADXIchimokuLongStrategy{baseStrategy}
}

func (s *ADXIchimokuLongStrategy) GetName() string {
	return "ADXBasicLongStrategy"
}

func (s *ADXIchimokuLongStrategy) GetDescription() string {
	return "ADX Basic Long Strategy"
}

func (s *ADXIchimokuLongStrategy) IsStrategyLive() bool {
	return true
}
