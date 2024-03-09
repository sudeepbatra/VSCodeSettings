package rules

type UnderIndicatorRule struct {
	first  []float64
	second []float64
}

func NewUnderIndicatorRule(first, second []float64) Rule {
	return &UnderIndicatorRule{
		first:  first,
		second: second,
	}
}

func (uir *UnderIndicatorRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(uir.first) || index >= len(uir.second) {
		return false, nil
	}

	satisfied := uir.first[index] < uir.second[index]

	return satisfied, nil
}
