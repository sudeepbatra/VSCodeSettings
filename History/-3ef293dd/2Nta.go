package common

type Rule interface {
	IsSatisfied(index int) bool
}

func And(r1, r2 Rule) Rule {
	return andRule{r1, r2}
}
