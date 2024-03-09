package ta

import (
	"errors"

	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyTypeLong  = "LONG"
	StrategyTypeShort = "SHORT"
)

var (
	ErrInvalidStrategyType = errors.New("invalid strategy type. Must be 'LONG' or 'SHORT'")
)

type Strategy interface {
	GetName() string
	GetDescription() string
	GetStrategyType() string
	ShouldEnter(index int) bool
	ShouldExit(index int) bool
	LoggingAndReport() string
	SetParameters(params map[string]interface{}) error
}

type BaseStrategy struct {
	entryRule rules.Rule
	exitRule  rules.Rule
}

func NewBaseStrategy(entryRule, exitRule rules.Rule) *BaseStrategy {
	// if strategyType != StrategyTypeLong && strategyType != StrategyTypeShort {
	// 	return nil, ErrInvalidStrategyType
	// }

	return &BaseStrategy{
		entryRule: entryRule,
		exitRule:  exitRule,
	}
}

func (rs *BaseStrategy) ShouldEnter(index int) bool {
	return rs.entryRule != nil && rs.entryRule.IsSatisfied(index)
}

func (rs *BaseStrategy) ShouldExit(index int) bool {
	return rs.exitRule != nil && rs.exitRule.IsSatisfied(index)
}

func (rs *BaseStrategy) LoggingAndReport() string {
	return "Strategy report"
}
