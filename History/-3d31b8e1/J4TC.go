package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type IchimokuLowHighEmaLongStrategy struct {
	*BaseStrategy
}

func NewIchimokuLowHighEmaLongStrategy(low, high, close, ema10, ema20, chikou, psar []float64) Strategy {
	chikouCrossUpLowHighRule := rules.NewChikouCrossUpLowHighRule(chikou, high, low, traceRule)
	ema10CrossUpEma20Rule := rules.NewCrossUpRule(ema10, ema20, traceRule)

	entryCrossUpRule := rules.Or(chikouCrossUpLowHighRule, ema10CrossUpEma20Rule)

	chikouAboveHighRule := rules.NewChikouAboveRule(chikou, high, traceRule)
	ema10OverEma20Rule := rules.NewOverIndicatorRule(ema10, ema20, traceRule)
	closeOverPsarRule := rules.NewOverIndicatorRule(close, psar, traceRule)

	entryVerificationRule := rules.And(chikouAboveHighRule, ema10OverEma20Rule, closeOverPsarRule)

	entryRule := rules.And(entryCrossUpRule, entryVerificationRule)

	chikouCrossDownLowRule := rules.NewChikouCrossDownRule(chikou, low, true)
	ema10CrossDownEma20Rule := rules.NewCrossDownRule(ema10, ema20, traceRule)
	exitRule := rules.Or(ema10CrossDownEma20Rule, chikouCrossDownLowRule)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating ichimoku ema long strategy")
	}

	return &IchimokuLowHighEmaLongStrategy{baseStrategy}
}

func (s *IchimokuLowHighEmaLongStrategy) GetName() string {
	return "IchimokuLowHighEmaLongStrategy"
}

func (s *IchimokuLowHighEmaLongStrategy) GetDescription() string {
	return "Ichimoku Low High Ema Long Strategy"
}

func (s *IchimokuLowHighEmaLongStrategy) IsStrategyLive() bool {
	return true
}
