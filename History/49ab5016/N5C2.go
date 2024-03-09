package rules

type InPipeRule struct {
	ref   []float64
	upper []float64
	lower []float64
}

func NewInPipeRule(ref, upper, lower []float64) Rule {
	return &InPipeRule{
		ref:   ref,
		upper: upper,
		lower: lower,
	}
}

func (ipr *InPipeRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(ipr.ref) {
		return false
	}

	refValue := ipr.ref[index]
	upperValue := ipr.upper[index]
	lowerValue := ipr.lower[index]

	satisfied := refValue >= lowerValue && refValue <= upperValue

	return satisfied
}
