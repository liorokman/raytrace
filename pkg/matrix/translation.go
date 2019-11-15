package matrix

import (
	"math"
)

func NewIdentity() Matrix {
	return Matrix{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func NewTranslation(x, y, z float64) Matrix {
	return Matrix{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func (m Matrix) Translate(x, y, z float64) Matrix {
	return m.Multiply(NewTranslation(x, y, z))
}

func NewScale(x, y, z float64) Matrix {
	return Matrix{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) Scale(x, y, z float64) Matrix {
	return m.Multiply(NewScale(x, y, z))
}

func NewRotateX(radians float64) Matrix {
	return Matrix{
		{1, 0, 0, 0},
		{0, math.Cos(radians), -math.Sin(radians), 0},
		{0, math.Sin(radians), math.Cos(radians), 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) RotateX(rad float64) Matrix {
	return m.Multiply(NewRotateX(rad))
}

func NewRotateY(radians float64) Matrix {
	return Matrix{
		{math.Cos(radians), 0, math.Sin(radians), 0},
		{0, 1, 0, 0},
		{-math.Sin(radians), 0, math.Cos(radians), 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) RotateY(rad float64) Matrix {
	return m.Multiply(NewRotateY(rad))
}

func NewRotateZ(radians float64) Matrix {
	return Matrix{
		{math.Cos(radians), -math.Sin(radians), 0, 0},
		{math.Sin(radians), math.Cos(radians), 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) RotateZ(rad float64) Matrix {
	return m.Multiply(NewRotateZ(rad))
}

func NewShear(xy, xz, yx, yz, zx, zy float64) Matrix {
	return Matrix{
		{1, xy, xz, 0},
		{yx, 1, yz, 0},
		{zx, zy, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) Shear(xy, xz, yx, yz, zx, zy float64) Matrix {
	return m.Multiply(NewShear(xy, xz, yx, yz, zx, zy))
}
