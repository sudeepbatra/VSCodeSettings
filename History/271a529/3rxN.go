package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuEmaLongStrategy struct {
	*BaseStrategy
}

func NewIchimokuEmaLongStrategy(close, ema10, ema20, chikou, psar []float64) Strategy {
	chikouCrossUpRule := rules.NewChikouCrossUpRule(chikou, close, true)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, traceRule)

	entryCrossUpRule := rules.Or(chikouCrossUpRule, ema10CrossUpEma20Rule)

	chikouAboveRule := rules.NewChikouAboveRule(chikou, close, false)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, traceRule)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, traceRule)

	entryVerificaionRule := rules.And(chikouAboveRule, ema10OverEma20Rule, closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificaionRule)

	exitRule := rules.NewCrossDownRule(ema10, ema20, traceRule)

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