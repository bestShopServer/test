package util

import "math"

const MIN = 0.000001

// Func: Compare Two Float64 Num
func Float64Equal(f1, f2 float64) bool {
	if math.Dim(f1, f2) < MIN {
		return true
	}
	return false
}

// Func: Min
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
