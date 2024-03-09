package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXBasicLongStrategy struct {
	*BaseStrategy
}

func NewADXBasicLongStrategy(adx, plusDI, minusDI []float64, adxThreshold float64) Strategy {
	plusDICrossAboveMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, true)
	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold, true)

	adxThresholdCrossoverRule := rules.NewCrossAboveRule(adx, adxThreshold)
	plusDIAboveMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI)

	plusDICrossoverEntryRule := rules.And(plusDICrossAboveMinusDIRule, adxOverThresholdRule)
	adxThresholdCrossoverEntryRule := rules.And(adxThresholdCrossoverRule, plusDIAboveMinusDIRule)

	entryRule := rules.Or(plusDICrossoverEntryRule, adxThresholdCrossoverEntryRule)

	exitRule := rules.NewCrossDownRule(plusDI, minusDI)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &ADXBasicLongStrategy{baseStrategy}
}

func (s *ADXBasicLongStrategy) GetName() string {
	return "ADXBasicLongStrategy"
}

func (s *ADXBasicLongStrategy) GetDescription() string {
	return "ADX Basic Long Strategy"
}

func IsStrategyLive() bool {
	return true
}
