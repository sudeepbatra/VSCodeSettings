package common

type Rule interface {
	IsSatisfied(index int) bool
}

type andRule struct {
	r1 Rule
	r2 Rule
}

func And(r1, r2 Rule) Rule {
	return andRule{r1, r2}
}

func (ar *andRule) IsSatisfied(index int) (bool, error) {
	s1, err := ar.r1.IsSatisfied(index)
	if err != nil {
		return false, err
	}

	s2, err := ar.r2.IsSatisfied(index)
	if err != nil {
		return false, err
	}

	return s1 && s2, nil
}
