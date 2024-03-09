package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuEmaLongStrategy struct {
	*BaseStrategy
}

func NewADXEMIPSARLongStrategy(close, ema10, ema20, chikou []float64) Strategy {
	rules.NewChikouAboveRule(chikou, close, false)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, traceRule)

	entryCrossUpRule := rules.Or(plusDICrossUpMinusDIRule, ema10CrossUpEma20Rule)

	plusDIOverMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, traceRule)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, traceRule)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, traceRule)

	entryVerificaionRule := rules.And(plusDIOverMinusDIRule, ema10OverEma20Rule, closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificaionRule)

	closeUnderPsarRule := rules.NewUnderIndicatorRule(close, psar, traceRule)
	closeUnderEma20Rule := rules.NewUnderIndicatorRule(close, ema20, traceRule)

	exitRule := rules.Or(closeUnderPsarRule, closeUnderEma20Rule)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating new adx emi psar long strategy")
	}

	return &ADXEMIPSARLongStrategy{baseStrategy}
}

func (s *ADXEMIPSARLongStrategy) GetName() string {
	return "ADXEMIPSARLongStrategy"
}

func (s *ADXEMIPSARLongStrategy) GetDescription() string {
	return "ADX EMI PSAR Long Strategy"
}

func (s *ADXEMIPSARLongStrategy) IsStrategyLive() bool {
	return true
}
