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

func (cr crossUpRule) IsSatisfied(index int) (bool, error) {
	if index <= 0 || index >= len(cr.upper) {
		return false, nil
	}

	if cr.lower[index-1] < cr.upper[index-1] && cr.lower[index] >= cr.upper[index] {
		return true, nil
	}

	return false, nil
}

func (cr crossDownRule) IsSatisfied(index int) (bool, error) {
	if index <= 0 || index >= len(cr.upper) {
		return false, nil
	}

	if cr.upper[index-1] < cr.lower[index-1] && cr.upper[index] >= cr.lower[index] {
		return true, nil
	}

	return false, nil
}
