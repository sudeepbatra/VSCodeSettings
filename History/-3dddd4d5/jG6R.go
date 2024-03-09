package rules

type orRule struct {
	rules []Rule
}

func Or(rules ...Rule) Rule {
	return &orRule{rules}
}

func (or *orRule) IsSatisfied(index int) bool {
	for _, rule := range or.rules {
		satisfied := rule.IsSatisfied(index)

		if satisfied {
			return true
		}
	}

	return false
}
