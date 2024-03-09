package rules

type andRule struct {
	rules []Rule
}

func And(rules ...Rule) Rule {
	return &andRule{rules}
}

func (ar *andRule) IsSatisfied(index int) bool {
	for _, rule := range ar.rules {
		satisfied := rule.IsSatisfied(index)

		if !satisfied {
			return false
		}
	}

	return true
}
