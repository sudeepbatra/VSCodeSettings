package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXEMIRSILongStrategy struct {
	*BaseStrategy
}

func NewADXEMIRSILongStrategy(close, adx, plusDI, minusDI, sma20, ema10, ema20, ema60, rsi, psar []float64, adxThreshold float64, rsiThreshold float64) Strategy {
	plusDICrossAboveMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, true)
	ema10OverEma20Rule := rules.NewCrossUpRule(ema10, ema20, true)
	adxCrossAboveThresholdRule := rules.NewCrossAboveRule(adx, adxThreshold, true)
	rsiCrossAboveThresholdRule := rules.NewCrossAboveRule(rsi, rsiThreshold, true)
	closeCrossUpSMA20EntryRule := rules.NewCrossUpRule(close, sma20, true)
	closeCrossUpEMA60EntryRule := rules.NewCrossUpRule(close, ema60, true)
	closeCrossUpPsarEntryRule := rules.NewCrossUpRule(close, psar, true)
	// closeCrossUpPsarEntryRule := rules.NewCrossDownRule(close, psar, false)

	// Rule entryCrossUpRule := rules.Or(plusDICrossAboveMinusDIRule, ema10OverEma20Rule, adxCrossAboveThresholdRule, rsiCrossAboveThresholdRule)

	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold, true)

	plusDIAboveMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, true)

	plusDICrossoverEntryRule := rules.And(plusDICrossAboveMinusDIRule, adxOverThresholdRule)
	adxThresholdCrossoverEntryRule := rules.And(adxCrossAboveThresholdRule, plusDIAboveMinusDIRule)

	entryRule := rules.Or(plusDICrossoverEntryRule, adxThresholdCrossoverEntryRule, ema10OverEma20Rule)

	exitRule := rules.NewCrossDownRule(plusDI, minusDI, true)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &ADXEMIRSILongStrategy{baseStrategy}
}

func (s *ADXEMIRSILongStrategy) GetName() string {
	return "ADXBasicLongStrategy"
}

func (s *ADXEMIRSILongStrategy) GetDescription() string {
	return "ADX Basic Long Strategy"
}

func (s *ADXEMIRSILongStrategy) IsStrategyLive() bool {
	return true
}
