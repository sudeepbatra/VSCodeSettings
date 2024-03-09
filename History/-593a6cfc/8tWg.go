package ta

import (
	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

type ADXIchimokuLongStrategy struct {
	*BaseStrategy
}

func NewADXIchimokuLongStrategy(close, chikou, plusDI, minusDI []float64) Strategy {
	plusDICrossAboveMinusDIRule := rules.NewCrossUpRule(plusDI, minusDI, true)
	chikouCrossAboveCloseRule := rules.NewChikouCrossUpRule(chikou, close)

	entryCrossUpRule := rules.Or(plusDICrossAboveMinusDIRule, chikouCrossAboveCloseRule)

	plusDIAboveMinusDIRule := rules.NewOverIndicatorRule(plusDI, minusDI, true)
	exitRule := rules.NewChikouCrossDownRule(chikou, value)

	entryRule := rules.Or(plusDICrossoverEntryRule, adxThresholdCrossoverEntryRule)

	exitRule := rules.NewCrossDownRule(plusDI, minusDI, true)

	baseStrategy, err := NewBaseStrategy(StrategyTypeLong, entryRule, exitRule)
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return &ADXIchimokuLongStrategy{baseStrategy}
}

func (s *ADXIchimokuLongStrategy) GetName() string {
	return "ADXIchimokuLongStrategy"
}

func (s *ADXIchimokuLongStrategy) GetDescription() string {
	return "ADX Ichimoku Long Strategy"
}

func (s *ADXIchimokuLongStrategy) IsStrategyLive() bool {
	return true
}
