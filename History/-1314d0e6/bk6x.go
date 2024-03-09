package rules

import (
	"errors"
)

type ChikouCrossoverRule struct {
	Chikou []float64
	Value  []float64
}

// NewChikouCrossoverRule creates a new ChikouCrossoverRule with Chikou and Value slices.
func NewChikouCrossoverRule(chikou, value []float64) (Rule, error) {
	if len(chikou) != len(value) {
		return nil, errors.New("Chikou and Value slices must have the same length")
	}

	return ChikouCrossoverRule{
		Chikou: chikou,
		Value:  value,
	}, nil
}

// IsSatisfied checks if the Chikou Span crosses above the Value at the given lastIndex.
func (cccr ChikouCrossoverRule) IsSatisfied(lastIndex int) (bool, error) {
	if lastIndex < 26 || lastIndex >= len(cccr.Value) {
		return false, errors.New("Index out of range")
	}

	closeValueAtMinus26 := cccr.Value[lastIndex-26]
	chikouValueAtMinus26 := cccr.Chikou[lastIndex-26]

	closeValueAtMinus26Minus1 := cccr.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := cccr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 > closeValueAtMinus26 && chikouValueAtMinus26Minus1 < closeValueAtMinus26Minus1, nil
}
