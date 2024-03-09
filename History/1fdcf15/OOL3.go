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
