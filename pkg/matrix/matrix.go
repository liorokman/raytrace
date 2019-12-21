package matrix

import (
	"fmt"

	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type Matrix struct {
	rows, cols int
	data       []float64
}

func (m Matrix) Fill(d []float64) {
	copy(m.data[:], d[:])
}

func New(row, col int) Matrix {
	return Matrix{
		rows: row,
		cols: col,
		data: make([]float64, row*col),
	}
}

func (m Matrix) Equals(other Matrix) bool {
	if m.rows != other.rows || m.cols != other.cols {
		return false
	}
	for i := 0; i < m.rows*m.cols; i++ {
		if !utils.FloatEqual(m.data[i], other.data[i]) {
			return false
		}
	}
	return true
}

func (l Matrix) Multiply(r Matrix) Matrix {
	result := New(l.rows, r.cols)

	for i := 0; i < l.rows; i++ {
		for j := 0; j < r.cols; j++ {
			for t := 0; t < l.cols; t++ {
				result.data[i*r.cols+j] = result.data[i*r.cols+j] + l.data[i*l.cols+t]*
					r.data[t*r.cols+j]
			}
		}
	}
	return result
}

func (l Matrix) MultiplyTuple(r tuple.Tuple) tuple.Tuple {

	right := Matrix{
		rows: 4,
		cols: 1,
		data: []float64{r.X(), r.Y(), r.Z(), r.W()},
	}
	result := l.Multiply(right)
	return tuple.Tuple{result.data[0], result.data[1], result.data[2], result.data[3]}
}

func (l Matrix) Transpose() Matrix {
	result := New(l.rows, l.cols)
	for i := 0; i < l.rows; i++ {
		for j := 0; j < l.cols; j++ {
			result.data[i*result.rows+j] = l.data[j*l.rows+i]
		}
	}
	return result
}

func (l Matrix) Submatrix(row, col int) Matrix {

	result := New(l.rows-1, l.cols-1)

	for i := 0; i < l.rows; i++ {
		if i == row {
			continue
		}
		targetRow := i
		if i > row {
			targetRow--
		}
		copy(result.data[targetRow*result.cols:targetRow*result.cols+col], l.data[i*l.cols:i*l.cols+col])
		copy(result.data[targetRow*result.cols+col:(targetRow+1)*result.cols], l.data[i*l.cols+col+1:(i+1)*l.cols])
	}
	return result
}

func (m Matrix) Determinant() float64 {
	if m.rows == 2 && m.cols == 2 {
		return m.data[0]*m.data[3] - m.data[1]*m.data[2]
	}
	det := 0.0
	for col := 0; col < m.cols; col++ {
		det = det + m.data[col]*m.Cofactor(0, col)
	}
	return det
}

func (m Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

func (m Matrix) Cofactor(row, col int) float64 {
	retval := m.Minor(row, col)
	if (row+col)%2 == 1 {
		retval = -retval
	}
	return retval
}

func (m Matrix) Inverse() (Matrix, error) {
	d := m.Determinant()
	if utils.FloatEqual(d, 0.0) {
		return m, fmt.Errorf("Matrix is not invertable")
	}
	result := New(m.rows, m.cols)
	for row := 0; row < m.rows; row++ {
		for col := 0; col < m.cols; col++ {
			c := m.Cofactor(row, col)
			result.data[col*m.rows+row] = c / d
		}

	}
	return result, nil
}
