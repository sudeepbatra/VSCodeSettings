package rules

type orRule struct {
	rules []Rule
}

func Or(rules ...Rule) Rule {
	return &orRule{rules}
}

func (or *orRule) IsSatisfied(index int) bool {
	satisfied := false
	for _, rule := range or.rules {
		satisfied := rule.IsSatisfied(index)

		if satisfied {
			return satisfied
		}
	}

	return satisfied
}
