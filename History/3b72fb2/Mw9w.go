package rules

import "math"

type PercentChangeRule struct {
	Values  []float64
	Percent float64
}

func NewPercentChangeRule(values []float64, percent float64) Rule {
	return &PercentChangeRule{
		Values:  values,
		Percent: percent,
	}
}

func (pcr *PercentChangeRule) IsSatisfied(index int) bool {
	if index == 0 || index >= len(pcr.Values) {
		return false
	}

	cp := pcr.Values[index]
	cplast := pcr.Values[index-1]
	percentChange := ((cp - cplast) / cplast) * 100.0

	absPercentChange := math.Abs(percentChange)

	satisfied := absPercentChange > pcr.Percent

	return satisfied
}
