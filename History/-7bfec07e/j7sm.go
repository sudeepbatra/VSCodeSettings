package common

type crossRule struct {
	upper []float64
	lower []float64
	cmp   int
}

func NewCrossUpRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
		cmp:   talib.CmpGreater,
	}
}
