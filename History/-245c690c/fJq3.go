package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXEMIRSILongStrategy struct {
	*BaseStrategy
}

func NewADXEMIRSILongStrategy(close, adx, plusDI, minusDI, sma20, ema10, ema20, ema60, rsi, psar []float64, adxThreshold float64, rsiThreshold float64) Strategy {
	plusDICrossUpMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, true)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, true)
	adxCrossAboveThresholdRule := rules.NewCrossAboveRule(adx, adxThreshold, true)
	rsiCrossAboveThresholdRule := rules.NewCrossAboveRule(rsi, rsiThreshold, true)
	closeCrossUpSMA20EntryRule := rules.NewCrossUpRule(close, sma20, true)
	closeCrossUpEMA60EntryRule := rules.NewCrossUpRule(close, ema60, true)
	closeCrossUpPsarEntryRule := rules.NewCrossUpRule(close, psar, true)

	entryCrossUpRule := rules.Or(plusDICrossUpMinusDIRule,
		ema10CrossUpEma20Rule,
		adxCrossAboveThresholdRule,
		rsiCrossAboveThresholdRule,
		closeCrossUpSMA20EntryRule,
		closeCrossUpEMA60EntryRule,
		closeCrossUpPsarEntryRule)

	plusDIOverMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, true)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, true)
	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold, true)
	rsiOverThresholdRule := rules.NewOverThresholdIndicatorRule(rsi, rsiThreshold, true)
	closeOverSMA20Rule := rules.NewOverIndicatorRule(close, sma20, true)
	closeOverEMA60Rule := rules.NewOverIndicatorRule(close, ema60, true)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, true)

	entryVerificaionRule := rules.And(plusDIOverMinusDIRule,
		ema10OverEma20Rule,
		adxOverThresholdRule,
		rsiOverThresholdRule,
		closeOverSMA20Rule,
		closeOverEMA60Rule,
		closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificaionRule)

	ema10UnderEma20Rule := rules.NewUnderIndicatorRule(ema10, ema20, true)
	closeUnderEma60Rule := rules.NewUnderIndicatorRule(close, ema60, true)
	closeUnderPsarRule := rules.NewUnderIndicatorRule(close, psar, true)

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
