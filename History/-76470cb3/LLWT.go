package ta

const (
	PSARStrategyName        = "PSAR"
	PSARStrategyDescription = "PSAR"
)

type PSARStrategy struct {
	Name         string
	Description  string
	StrategyType string
	close        []float64
	psar         []float64
}

func NewPSARStrategy() *PSARStrategy {
	return &PSARStrategy{
		Name:         PSARStrategyName,
		Description:  PSARStrategyDescription,
		StrategyType: StrategyTypeLong,
	}
}

// func NewPSARStrategy(close, psar []float64) Strategy {
// 	rsiAboveStrategy, err := NewStrategy(
// 		PSARStrategyName,
// 		PSARStrategyDescription,
// 		StrategyTypeLong,
// 		rules.NewCrossUpRule(close, psar),
// 		rules.NewCrossDownRule(close, psar))
// 	if err != nil {
// 		logger.Log.Error().Err(err).Msg("error in creating strategy")
// 	}

// 	return rsiAboveStrategy
// }
