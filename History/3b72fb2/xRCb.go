package common

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

func (pcr *PercentChangeRule) IsSatisfied(index int) (bool, error) {
	if index == 0 || index >= len(pcr.Values) {
		return false, nil
	}

	cp := pcr.Values[index]
	cplast := pcr.Values[index-1]
	percentChange := ((cp - cplast) / cplast) * 100.0

	// Calculate the absolute value of the percent change
	absPercentChange := math.Abs(percentChange)

	satisfied := absPercentChange > pcr.Percent

	return satisfied, nil
}

type PercentChangeIndicator struct {
	Values []float64
}

func NewPercentChangeIndicator(values []float64) Indicator {
	return &PercentChangeIndicator{
		Values: values,
	}
}

func (pci *PercentChangeIndicator) Calculate(index int) float64 {
	if index == 0 || index >= len(pci.Values) {
		return 0.0
	}

	cp := pci.Values[index]
	cplast := pci.Values[index-1]
	return ((cp - cplast) / cplast) * 100.0
}
