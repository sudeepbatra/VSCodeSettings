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
