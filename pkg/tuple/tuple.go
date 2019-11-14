package tuple

import (
	"math"
)

const (
	XPos int = iota
	YPos int = iota
	ZPos int = iota
	WPos int = iota
)

type Tuple struct {
	Values [4]float64
}

func (t Tuple) X() float64 {
	return t.Values[XPos]
}

func (t Tuple) Y() float64 {
	return t.Values[YPos]
}

func (t Tuple) Z() float64 {
	return t.Values[ZPos]
}

func (t Tuple) W() float64 {
	return t.Values[WPos]
}

func (t Tuple) IsVector() bool {
	return t.Values[WPos] == 0.0
}

func (t Tuple) IsPoint() bool {
	return t.Values[WPos] == 1.0
}

func NewPoint(x, y, z float64) Tuple {
	return Tuple{[4]float64{x, y, z, 1.0}}
}

func NewVector(x, y, z float64) Tuple {
	return Tuple{[4]float64{x, y, z, 0.0}}
}

func (t Tuple) Add(right Tuple) Tuple {
	return Tuple{[4]float64{
		t.X() + right.X(),
		t.Y() + right.Y(),
		t.Z() + right.Z(),
		t.W() + right.W(),
	}}
}

func (t Tuple) Subtract(right Tuple) Tuple {
	return Tuple{[4]float64{
		t.X() - right.X(),
		t.Y() - right.Y(),
		t.Z() - right.Z(),
		t.W() - right.W(),
	}}
}

func (t Tuple) Negate() Tuple {
	return Tuple{[4]float64{
		-t.X(),
		-t.Y(),
		-t.Z(),
		-t.W(),
	}}
}

func (t Tuple) Mult(scalar float64) Tuple {
	return Tuple{[4]float64{
		scalar * t.X(),
		scalar * t.Y(),
		scalar * t.Z(),
		scalar * t.W(),
	}}
}

func (t Tuple) Div(scalar float64) Tuple {
	return Tuple{[4]float64{
		t.X() / scalar,
		t.Y() / scalar,
		t.Z() / scalar,
		t.W() / scalar,
	}}
}

func (t Tuple) Magnitude() float64 {
	return math.Sqrt(t.X()*t.X() +
		t.Y()*t.Y() +
		t.Z()*t.Z() +
		t.W()*t.W())
}

func (t Tuple) Normalize() Tuple {
	n := t.Magnitude()
	return t.Div(n)
}

func (t Tuple) Dot(r Tuple) float64 {
	res := 0.0
	for i, _ := range t.Values {
		res = res + t.Values[i]*r.Values[i]
	}
	return res
}

func (t Tuple) Cross(r Tuple) Tuple {
	if !t.IsVector() || !r.IsVector() {
		panic("cross product only supported for two vectors")
	}
	return NewVector(
		t.Y()*r.Z()-t.Z()*r.Y(),
		t.Z()*r.X()-t.X()*r.Z(),
		t.X()*r.Y()-t.Y()*r.X())
}
