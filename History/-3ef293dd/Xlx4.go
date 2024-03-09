package common

type Rule interface {
	IsSatisfied(index int) bool
}
