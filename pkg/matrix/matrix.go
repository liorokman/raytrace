package matrix

import (
	"fmt"
	"math"

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

func (l Matrix) Transpose() Matrix {
	result := New(len(l), len(l[0]))
	for i := range l {
		for j := range l[0] {
			result[i][j] = l[j][i]
		}
	}
	return result
}

func (l Matrix) Submatrix(row, col int) Matrix {

	result := New(len(l)-1, len(l[0])-1)

	for index := range l {
		realI := index
		switch {
		case index == row:
			continue
		case index > row:
			realI = realI - 1
		}
		copy(result[realI][0:col], l[index][0:col])
		copy(result[realI][col:], l[index][col+1:])
	}

	return result
}

func (m Matrix) Determinant() float64 {
	if len(m) == 2 && len(m[0]) == 2 {
		return m[0][0]*m[1][1] - m[0][1]*m[1][0]
	}
	det := 0.0
	for col := range m[0] {
		det = det + m[0][col]*m.Cofactor(0, col)
	}
	return det
}

func (m Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

func (m Matrix) Cofactor(row, col int) float64 {
	return math.Pow(-1, float64(row+col)) * m.Minor(row, col)
}

func (m Matrix) Inverse() (Matrix, error) {
	d := m.Determinant()
	if utils.FloatEqual(d, 0.0) {
		return m, fmt.Errorf("Matrix is not invertable")
	}
	result := New(len(m), len(m[0]))
	for row := range m {
		for col := range m[0] {
			c := m.Cofactor(row, col)
			result[col][row] = c / d
		}

	}
	return result, nil
}
