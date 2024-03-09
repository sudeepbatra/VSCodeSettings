package rules

type OverThresholdIndicatorRule struct {
	first  []float64
	second []float64
}

func NewOverThresholdIndicatorRule(first, second []float64) Rule {
	return &OverIndicatorRule{
		first:  first,
		second: second,
	}
}

func (oir *OverThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) || index >= len(oir.second) {
		return false
	}

	satisfied := oir.first[index] > oir.second[index]

	return satisfied
}
