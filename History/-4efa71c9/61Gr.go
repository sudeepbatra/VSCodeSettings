package common

type notRule struct {
	rule Rule
}

func Not(rule Rule) Rule {
	return &notRule{rule}
}

func (nr *notRule) IsSatisfied(index int) (bool, error) {
	satisfied, err := nr.rule.IsSatisfied(index)
	if err != nil {
		return false, err
	}

	return !satisfied, nil
}
