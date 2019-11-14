package utils

import (
	"math"
)

const EPSILON float64 = 0.00001

func FloatEqual(l, r float64) bool {
	return math.Abs(l-r) < EPSILON
}
