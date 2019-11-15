package matrix

import (
	"sync"

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

func (l Matrix) Multiple(r Matrix) Matrix {
	result := New(len(l), len(r[0]))

	var wg sync.WaitGroup
	for rw := range l {
		wg.Add(1)
		go func(row int) {
			for i := range l[row] {
				for j := range r {
					result[row][i] = result[row][i] + l[row][j]*r[j][i]
				}
			}
			wg.Done()
		}(rw)
	}
	wg.Wait()

	return result
}

func (l Matrix) Multiple2(r Matrix) Matrix {
	result := New(len(l), len(r[0]))

	for i := range l {
		for j := range l[i] {

			result[i][j] =
				l[i][0]*r[0][j] +
					l[i][1]*r[1][j] +
					l[i][2]*r[2][j] +
					l[i][3]*r[3][j]
		}
	}
	return result
}
