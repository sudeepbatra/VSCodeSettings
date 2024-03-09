package rules

type UnderIndicatorRule struct {
	a []float64
	b []float64
}

func NewUnderIndicatorRule(a, b []float64) Rule {
	return &UnderIndicatorRule{
		a: a,
		b: b,
	}
}

func (uir *UnderIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(uir.a) || index >= len(uir.b) {
		return false
	}

	satisfied := uir.a[index] < uir.b[index]
	return uir.a[index] < uir.b[index]
}
