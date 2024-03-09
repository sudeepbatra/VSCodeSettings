package ta

import (
	"errors"

	"github.com/sudeepbatra/alpha-hft/logger"
	"github.com/sudeepbatra/alpha-hft/ta/rules"
)

const (
	StrategyTypeLong  = "LONG"
	StrategyTypeShort = "SHORT"
)

var (
	ErrEntryRuleIsNil      = errors.New("entry rule is nil")
	ErrExitRuleIsNil       = errors.New("exit rule is nil")
	ErrInvalidStrategyType = errors.New("invalid strategy type")
)

type Strategy interface {
	GetName() string
	GetDescription() string
	GetStrategyType() string
	ShouldEnter(index int) (bool, error)
	ShouldExit(index int) (bool, error)
	LoggingAndReport() string
	SetParameters(params map[string]interface{}) error
}

type RuleStrategy struct {
	Name         string
	Description  string
	StrategyType string
	EntryRule    rules.Rule
	ExitRule     rules.Rule
	Parameters   map[string]interface{}
}

func NewStrategy(name, description, strategyType string, entryRule, exitRule rules.Rule) (*RuleStrategy, error) {
	if strategyType != StrategyTypeLong && strategyType != StrategyTypeShort {
		return nil, ErrInvalidStrategyType
	}

	return &RuleStrategy{
		Name:         name,
		Description:  description,
		StrategyType: strategyType,
		EntryRule:    entryRule,
		ExitRule:     exitRule,
		Parameters:   make(map[string]interface{}),
	}, nil
}

func (rs *RuleStrategy) GetName() string {
	return rs.Name
}

func (rs *RuleStrategy) GetDescription() string {
	return rs.Description
}

func (rs *RuleStrategy) GetStrategyType() string {
	return rs.StrategyType
}

func (rs *RuleStrategy) ShouldEnter(index int) (bool, error) {
	if rs.EntryRule == nil {
		return false, ErrEntryRuleIsNil
	}

	return rs.EntryRule.IsSatisfied(index)
}

func (rs *RuleStrategy) ShouldExit(index int) (bool, error) {
	if rs.ExitRule == nil {
		return false, ErrExitRuleIsNil
	}

	return rs.ExitRule.IsSatisfied(index)
}

func (rs *RuleStrategy) LoggingAndReport() string {
	return "Strategy report"
}

func (rs *RuleStrategy) SetParameters(params map[string]interface{}) error {
	rs.Parameters = params
	return nil
}

const (
	StrategyName        = "RSI Above 70"
	StrategyDescription = "Buy when RSI above 70"
	LongThreshold       = 70
	ExitThreshold       = 30
)

func NewRSIAboveStrategy(rsiSeries []float64) Strategy {
	rsiAboveStrategy, err := NewStrategy(
		StrategyName,
		StrategyDescription,
		StrategyTypeLong,
		rules.NewCrossAboveRule(rsiSeries, LongThreshold),
		rules.NewCrossBelowRule(rsiSeries, ExitThreshold))
	if err != nil {
		logger.Log.Error().Err(err).Msg("error in creating strategy")
	}

	return rsiAboveStrategy
}
