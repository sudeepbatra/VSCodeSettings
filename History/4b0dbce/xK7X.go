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
	ErrInvalidStrategyType = errors.New("invalid strategy type")
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
	StrategyType string
	entryRule    rules.Rule
	exitRule     rules.Rule
}

func NewBaseStrategy(name, description, strategyType string, entryRule, exitRule rules.Rule) (*BaseStrategy, error) {
	if strategyType != StrategyTypeLong && strategyType != StrategyTypeShort {
		return nil, ErrInvalidStrategyType
	}

	return &BaseStrategy{
		StrategyType: strategyType,
		entryRule:    entryRule,
		exitRule:     exitRule,
	}, nil
}

func (rs *BaseStrategy) GetName() string {
	return rs.Name
}

func (rs *BaseStrategy) GetDescription() string {
	return rs.Description
}

func (rs *BaseStrategy) GetStrategyType() string {
	return rs.StrategyType
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

func (rs *BaseStrategy) SetParameters(params map[string]interface{}) error {
	rs.Parameters = params
	return nil
}
