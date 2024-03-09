package common

func CheckSameSize(values ...[]float64) {
	if len(values) < 2 {
		return
	}

	n := len(values[0])

	for i := 1; i < len(values); i++ {
		if len(values[i]) != n {
			panic("not all same size")
		}
	}
}
