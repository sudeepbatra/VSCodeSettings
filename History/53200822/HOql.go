package rules

type UnderThresholdIndicatorRule struct {
	first     []float64
	threshold float64
}

func NewUnderThresholdIndicatorRule(first []float64, threshold float64) Rule {
	return &UnderThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
	}
}

func (uir *UnderThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(uir.first) {
		return false
	}

	return uir.first[index] < uir.threshold
}
