package common

import "reflect"

func CheckSameSizeFloat64(values ...[]float64) {
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

func CheckSameSize(values ...interface{}) bool {
	if len(values) < 2 {
		return true
	}

	n := reflect.ValueOf(values[0]).Len()

	for i := 1; i < len(values); i++ {
		if reflect.ValueOf(values[i]).Len() != n {
			return false
		}
	}

	return true
}
