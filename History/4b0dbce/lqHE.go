package ta

const (
	StrategyTypeLong  = "LONG"
	StrategyTypeShort = "SHORT"
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
