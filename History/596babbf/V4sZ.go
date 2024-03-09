package common

type xorRule struct {
	rules []Rule
}

func Xor(rules ...Rule) Rule {
	return &xorRule{rules}
}

func (xr *xorRule) IsSatisfied(index int) (bool, error) {
	satisfiedCount := 0

	for _, rule := range xr.rules {
		satisfied, err := rule.IsSatisfied(index)
		if err != nil {
			return false, err
		}

		if satisfied {
			satisfiedCount++
		}
	}

	// XOR is satisfied when an odd number of rules are satisfied
	return satisfiedCount%2 == 1, nil
}
