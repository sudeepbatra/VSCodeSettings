package common

type OverIndicatorRule struct {
	first  []float64
	second []float64
}

func NewOverIndicatorRule(first, second []float64) Rule {
	return &OverIndicatorRule{
		first:  first,
		second: second,
	}
}

func (oir *OverIndicatorRule) IsSatisfied(index int) (bool, error) {
	if index < 0 || index >= len(oir.first) || index >= len(oir.second) {
		return false, nil
	}

	satisfied := oir.first[index] > oir.second[index]

	return satisfied, nil
}
