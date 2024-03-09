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

	entryVerificationRule := rules.And(chikouAboveRule, ema10OverEma20Rule, closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificationRule)

	chikouCrossUpRule := rules.NewChikouCrossUpRule(chikou, close, true)

	exitRule := rules.NewCrossDownRule(ema10, ema20, traceRule)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating ichimoku ema long strategy")
	}

	return &IchimokuEmaLongStrategy{baseStrategy}
}

func (s *IchimokuEmaLongStrategy) GetName() string {
	return "IchimokuEmaLongStrategy"
}

func (s *IchimokuEmaLongStrategy) GetDescription() string {
	return "Ichimoku Ema Long Strategy"
}

func (s *IchimokuEmaLongStrategy) IsStrategyLive() bool {
	return true
}
