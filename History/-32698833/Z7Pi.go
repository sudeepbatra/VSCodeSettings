package ta

type Rule interface {
	IsSatisfied(index int) (bool, error)
}
