package rules

type crossUpRule struct {
	upper []float64
	lower []float64
}

type crossDownRule struct {
	upper []float64
	lower []float64
}

func NewCrossUpRule(upper, lower []float64) Rule {
	return crossUpRule{
		upper: upper,
		lower: lower,
	}
}

func NewCrossDownRule(upper, lower []float64) Rule {
	return crossDownRule{
		upper: upper,
		lower: lower,
	}
}
