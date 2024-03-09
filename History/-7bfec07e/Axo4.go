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
		cmp:   1,
	}
}

func NewCrossDownRule(upper, lower []float64) Rule {
	return crossRule{
		upper: upper,
		lower: lower,
		cmp:   -1, // Negative value for cross-down
	}
}

type Rule interface {
	IsSatisfied(index int) bool
}

type crossRule struct {
	upper []float64
	lower []float64
	cmp   int
}
