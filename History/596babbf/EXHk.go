package common

type xorRule struct {
	r1 Rule
	r2 Rule
}

func Xor(r1, r2 Rule) Rule {
	return &xorRule{r1, r2}
}

func (xr *xorRule) IsSatisfied(index int) (bool, error) {
	satisfied1, err1 := xr.r1.IsSatisfied(index)
	satisfied2, err2 := xr.r2.IsSatisfied(index)

	if err1 != nil {
		return false, err1
	}

	if err2 != nil {
		return false, err2
	}

	return (satisfied1 || satisfied2) && !(satisfied1 && satisfied2), nil
}
