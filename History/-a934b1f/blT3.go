package rules

type notRule struct {
	rule Rule
}

func Not(rule Rule) Rule {
	return &notRule{rule}
}

func (nr *notRule) IsSatisfied(index int) bool {
	satisfied := nr.rule.IsSatisfied(index)

	return !satisfied
}
