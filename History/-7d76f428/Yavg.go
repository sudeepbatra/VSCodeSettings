package rules

type OverThresholdIndicatorRule struct {
	first     []float64
	threshold float64
}

func NewOverThresholdIndicatorRule(first []float64, threshold float64) Rule {
	return &OverThresholdIndicatorRule{
		first:     first,
		threshold: threshold,
	}
}

func (oir *OverThresholdIndicatorRule) IsSatisfied(index int) bool {
	if index < 0 || index >= len(oir.first) {
		return false
	}

	satisfied := oir.first[index] > oir.second[index]

	return satisfied
}
