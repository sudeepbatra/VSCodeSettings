package ta

import (
	"errors"
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
