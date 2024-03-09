package rules

import "errors"

type ChikouCrossDownRule struct {
	Chikou []float64
	Value  []float64
}

// NewChikouCrossDownRule creates a new ChikouCrossDownRule with Chikou and Value slices.
func NewChikouCrossDownRule(chikou, value []float64) (Rule, error) {
	if len(chikou) != len(value) {
		return nil, errors.New("Chikou and Value slices must have the same length")
	}

	return ChikouCrossDownRule{
		Chikou: chikou,
		Value:  value,
	}, nil
}

// IsSatisfied checks if the Chikou Span crosses below the Value at the given lastIndex.
func (ccdr ChikouCrossDownRule) IsSatisfied(lastIndex int) (bool, error) {
	if lastIndex < 26 || lastIndex >= len(ccdr.Value) {
		return false, errors.New("Index out of range")
	}

	valueAtMinus26 := ccdr.Value[lastIndex-26]
	chikouValueAtMinus26 := ccdr.Chikou[lastIndex-26]

	valueAtMinus26Minus1 := ccdr.Value[lastIndex-26-1]
	chikouValueAtMinus26Minus1 := ccdr.Chikou[lastIndex-26-1]

	return chikouValueAtMinus26 < valueAtMinus26 &&
		chikouValueAtMinus26Minus1 > valueAtMinus26Minus1, nil
}
