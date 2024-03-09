package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	traceRule = true
)

type ADXEMIRSILongStrategy struct {
	*BaseStrategy
}

func NewADXEMIRSILongStrategy(close, adx, plusDI, minusDI, sma20, ema10,
	ema20, ema60, rsi, psar []float64, adxThreshold float64, rsiThreshold float64,
) Strategy {
	plusDICrossUpMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, traceRule)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, traceRule)
	adxCrossAboveThresholdRule := rules.NewCrossAboveThresholdRule(adx, adxThreshold, traceRule)
	rsiCrossAboveThresholdRule := rules.NewCrossAboveThresholdRule(rsi, rsiThreshold, traceRule)
	closeCrossUpSma20EntryRule := rules.NewCrossUpRule(close, sma20, traceRule)
	closeCrossUpEma60EntryRule := rules.NewCrossUpRule(close, ema60, traceRule)
	closeCrossUpPsarEntryRule := rules.NewCrossUpRule(close, psar, traceRule)

	entryCrossUpRule := rules.Or(plusDICrossUpMinusDIRule,
		ema10CrossUpEma20Rule,
		adxCrossAboveThresholdRule,
		rsiCrossAboveThresholdRule,
		closeCrossUpSma20EntryRule,
		closeCrossUpEma60EntryRule,
		closeCrossUpPsarEntryRule)

	plusDIOverMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, traceRule)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, traceRule)
	adxOverThresholdRule := rules.NewOverThresholdIndicatorRule(adx, adxThreshold, traceRule)
	rsiOverThresholdRule := rules.NewOverThresholdIndicatorRule(rsi, rsiThreshold, traceRule)
	closeOverSma20Rule := rules.NewOverIndicatorRule(close, sma20, traceRule)
	closeOverEma60Rule := rules.NewOverIndicatorRule(close, ema60, traceRule)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, traceRule)

	entryVerificaionRule := rules.And(plusDIOverMinusDIRule,
		ema10OverEma20Rule,
		adxOverThresholdRule,
		rsiOverThresholdRule,
		closeOverSma20Rule,
		closeOverEma60Rule,
		closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificaionRule)

	ema10UnderEma20Rule := rules.NewUnderIndicatorRule(ema10, ema20, traceRule)
	closeUnderEma60Rule := rules.NewUnderIndicatorRule(close, ema60, traceRule)
	closeUnderPsarRule := rules.NewUnderIndicatorRule(close, psar, traceRule)
	closeUnderEma10Rule := rules.NewUnderIndicatorRule(close, ema10, traceRule)
	closeUnderEma20Rule := rules.NewUnderIndicatorRule(close, ema20, traceRule)

	exitRule := rules.Or(ema10UnderEma20Rule,
		closeUnderEma60Rule,
		closeUnderPsarRule,
		closeUnderEma10Rule,
		closeUnderEma20Rule)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating new adx emi rsi long strategy")
	}

	return &ADXEMIRSILongStrategy{baseStrategy}
}

func (s *ADXEMIRSILongStrategy) GetName() string {
	return "ADXEMIRSILongStrategy"
}

func (s *ADXEMIRSILongStrategy) GetDescription() string {
	return "ADX EMI RSI Long Strategy"
}

func (s *ADXEMIRSILongStrategy) IsStrategyLive() bool {
	return true
}
