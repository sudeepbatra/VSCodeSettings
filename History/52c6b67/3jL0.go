package rules

type IsEqualRule struct {
	first  []float64
	second []float64
}

func NewIsEqualRule(first, second []float64) Rule {
	return &IsEqualRule{
		first:  first,
		second: second,
	}
}

func (ier *IsEqualRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(ier.first) || index >= len(ier.second) {
		return false
	}

	return ier.first[index] == ier.second[index]
}
