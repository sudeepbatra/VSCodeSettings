package common

type Rule interface {
	IsSatisfied(index int) (bool, error)
}

type andRule struct {
	r1 Rule
	r2 Rule
}

func And(r1, r2 Rule) Rule {
	return &andRule{r1, r2}
}
