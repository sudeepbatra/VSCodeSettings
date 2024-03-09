package common

type notRule struct {
	rule Rule
}

func Not(rule Rule) Rule {
	return &notRule{rule}
}
