package rules

type xorRule struct {
	r1 Rule
	r2 Rule
}

func Xor(r1, r2 Rule) Rule {
	return &xorRule{r1, r2}
}

func (xr *xorRule) IsSatisfied(index int) bool {
	satisfied1 := xr.r1.IsSatisfied(index)
	satisfied2 := xr.r2.IsSatisfied(index)

	return (satisfied1 || satisfied2) && !(satisfied1 && satisfied2)
}
