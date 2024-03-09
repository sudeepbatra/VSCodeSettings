package common

type Rule interface {
	IsSatisfied(index int) (bool, error)
}

type andRule struct {
	rules []Rule
}

func And(rules ...Rule) Rule {
	return &andRule{rules}
}

func (ar *andRule) IsSatisfied(index int) (bool, error) {
	for _, rule := range ar.rules {
		satisfied, err := rule.IsSatisfied(index)
		if err != nil {
			return false, err
		}

		if !satisfied {
			return false, nil
		}
	}

	return true, nil
}

type orRule struct {
	rules []Rule
}

func Or(rules ...Rule) Rule {
	return &orRule{rules}
}

func (or *orRule) IsSatisfied(index int) (bool, error) {
	for _, rule := range or.rules {
		satisfied, err := rule.IsSatisfied(index)
		if err != nil {
			return false, err
		}

		if satisfied {
			return true, nil
		}
	}
	return false, nil
}
