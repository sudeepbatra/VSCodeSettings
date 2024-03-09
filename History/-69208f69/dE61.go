package common

import "math"

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func isNaN(val float64) bool {
	return math.IsNaN(val)
}
