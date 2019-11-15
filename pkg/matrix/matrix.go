package matrix

import (
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type Matrix [][]float64

func New(row, col int) Matrix {
	m := make([][]float64, row)
	for i, _ := range m {
		m[i] = make([]float64, col)
	}
	return m
}

func (m Matrix) Equals(other Matrix) bool {
	if len(m) != len(other) || len(m[0]) != len(other[0]) {
		return false
	}
	for i := range m {
		for j := range m[i] {
			if !utils.FloatEqual(m[i][j], other[i][j]) {
				return false
			}
		}
	}
	return true
}

func (l Matrix) Multiply(r Matrix) Matrix {
	result := New(len(l), len(r[0]))

	for i := range l {
		for j := range r[i] {
			for t := range l[j] {
				result[i][j] = result[i][j] + l[i][t]*r[t][j]
			}
		}
	}
	return result
}

func (l Matrix) MultiplyTuple(r tuple.Tuple) tuple.Tuple {

	right := Matrix{
		[]float64{r.X()},
		[]float64{r.Y()},
		[]float64{r.Z()},
		[]float64{r.W()},
	}
	result := l.Multiply(right)
	return tuple.Tuple{result[0][0], result[1][0], result[2][0], result[3][0]}
}
