package ta

import (
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyTypeLong  = "LONG"
	StrategyTypeShort = "SHORT"
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
	Name         string
	Description  string
	StrategyType string
	entryRule    rules.Rule
	exitRule     rules.Rule
	Parameters   map[string]interface{}
}

func NewBaseStrategy(name, description, strategyType string, entryRule, exitRule rules.Rule) (*BaseStrategy, error) {
	if strategyType != StrategyTypeLong && strategyType != StrategyTypeShort {
		return nil, ErrInvalidStrategyType
	}

	return &BaseStrategy{
		Name:         name,
		Description:  description,
		StrategyType: strategyType,
		entryRule:    entryRule,
		exitRule:     exitRule,
		Parameters:   make(map[string]interface{}),
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

	if rs.EntryRule == nil {
		return false, ErrEntryRuleIsNil
	}

	return rs.EntryRule.IsSatisfied(index)
}

func (rs *BaseStrategy) ShouldExit(index int) (bool, error) {
	if rs.ExitRule == nil {
		return false, ErrExitRuleIsNil
	}

	return rs.ExitRule.IsSatisfied(index)
}

func (rs *BaseStrategy) LoggingAndReport() string {
	return "Strategy report"
}

func (rs *BaseStrategy) SetParameters(params map[string]interface{}) error {
	rs.Parameters = params
	return nil
}
