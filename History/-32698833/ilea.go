package rules

type Rule interface {
	IsSatisfied(index int) bool
}
