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
	SetParameters(params map[string]interface{})
}

type BaseStrategy struct {
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
}

func (rs *BaseStrategy) ShouldExit(index int) bool {
	return rs.exitRule != nil && rs.exitRule.IsSatisfied(index)
}

func (rs *BaseStrategy) LoggingAndReport() string {
	return "Strategy report"
}

func (rs *BaseStrategy) SetParameters(params map[string]interface{}) {
	rs.Parameters = params
}
