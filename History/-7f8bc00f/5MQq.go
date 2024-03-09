package common

import "reflect"

func CheckSameSizeFloat64(values ...[]float64) bool {
	if len(values) < 2 {
		return true
	}

	n := len(values[0])

	for i := 1; i < len(values); i++ {
		if len(values[i]) != n {
			return false
		}
	}

	return true
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

func checkMinimumCandlesticks(values ...[]float64) bool {
	if len(values) < 2 {
		return false
	}

}
