package rules

type UnderThresholdIndicatorRule struct {
	first     []float64
	threshold float64
}

func NewUnderThresholdIndicatorRule(first, second []float64) Rule {
	return &UnderThresholdIndicatorRule{
		first:  first,
		second: second,
	}
}

func (uir *UnderIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(uir.first) || index >= len(uir.second) {
		return false
	}

	return uir.first[index] < uir.second[index]
}
