package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuLowHighEmaLongStrategy struct {
	*BaseStrategy
}

func NewIchimokuLowHighEmaLongStrategy(low, high, close, ema10, ema20, chikou, psar []float64) Strategy {
	chikouCrossUpRule := rules.NewChikouCrossUpLowHighRule(chikou, high, low, traceRule)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, traceRule)

	entryCrossUpRule := rules.Or(chikouCrossUpRule, ema10CrossUpEma20Rule)

	chikouAboveRule := rules.NewChikouAboveRule(chikou, high, traceRule)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, traceRule)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, traceRule)

	entryVerificationRule := rules.And(chikouAboveRule, ema10OverEma20Rule, closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificationRule)

	chikouCrossDownRule := rules.NewChikouCrossDownRule(chikou, close, true)
	ema10CrossDownEma20Rule := rules.NewCrossDownRule(ema10, ema20, traceRule)
	exitRule := rules.Or(ema10CrossDownEma20Rule, chikouCrossDownRule)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating ichimoku ema long strategy")
	}

	return &IchimokuLowHighEmaLongStrategy{baseStrategy}
}

func (s *IchimokuLowHighEmaLongStrategy) GetName() string {
	return "IchimokuEmaLongStrategy"
}

func (s *IchimokuLowHighEmaLongStrategy) GetDescription() string {
	return "Ichimoku Ema Long Strategy"
}

func (s *IchimokuLowHighEmaLongStrategy) IsStrategyLive() bool {
	return true
}
