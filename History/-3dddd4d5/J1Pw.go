package rules

type orRule struct {
	rules []Rule
}

func Or(rules ...Rule) Rule {
	return &orRule{rules}
}

func (or *orRule) IsSatisfied(index int) bool {
	for _, rule := range or.rules {
		satisfied, err := rule.IsSatisfied(index)
		if err != nil {
			return false, err
		}

		if satisfied {
			return true
		}
	}

	return false
}
