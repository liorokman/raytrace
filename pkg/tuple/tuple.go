package tuple

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/utils"
)

const (
	XPos int = iota
	YPos int = iota
	ZPos int = iota
	WPos int = iota
)

type Tuple [4]float64

func (t Tuple) String() string {
	return fmt.Sprintf("(%f,%f,%f,%f)", t[0], t[1], t[2], t[3])
}

func (t Tuple) X() float64 {
	return t[XPos]
}

func (t Tuple) Y() float64 {
	return t[YPos]
}

func (t Tuple) Z() float64 {
	return t[ZPos]
}

func (t Tuple) W() float64 {
	return t[WPos]
}

func (t Tuple) IsVector() bool {
	return t[WPos] == 0.0
}

func (t Tuple) IsPoint() bool {
	return t[WPos] == 1.0
}

func NewPoint(x, y, z float64) Tuple {
	return Tuple{x, y, z, 1.0}
}

func NewVector(x, y, z float64) Tuple {
	return Tuple{x, y, z, 0.0}
}

func (t Tuple) Equals(o Tuple) bool {
	return utils.FloatEqual(t.X(), o.X()) &&
		utils.FloatEqual(t.Y(), o.Y()) &&
		utils.FloatEqual(t.Z(), o.Z()) &&
		utils.FloatEqual(t.W(), o.W())

}

func (t Tuple) Add(right Tuple) Tuple {
	return Tuple{
		t.X() + right.X(),
		t.Y() + right.Y(),
		t.Z() + right.Z(),
		t.W() + right.W(),
	}
}

func (t Tuple) Subtract(right Tuple) Tuple {
	return Tuple{
		t.X() - right.X(),
		t.Y() - right.Y(),
		t.Z() - right.Z(),
		t.W() - right.W(),
	}
}

func (t Tuple) Negate() Tuple {
	return Tuple{
		-t.X(),
		-t.Y(),
		-t.Z(),
		-t.W(),
	}
}

func (t Tuple) Mult(scalar float64) Tuple {
	return Tuple{
		scalar * t.X(),
		scalar * t.Y(),
		scalar * t.Z(),
		scalar * t.W(),
	}
}

func (t Tuple) Div(scalar float64) Tuple {
	return Tuple{
		t.X() / scalar,
		t.Y() / scalar,
		t.Z() / scalar,
		t.W() / scalar,
	}
}

func (t Tuple) Magnitude() float64 {
	return math.Sqrt(t.X()*t.X() +
		t.Y()*t.Y() +
		t.Z()*t.Z() +
		t.W()*t.W())
}

func (t Tuple) Normalize() Tuple {
	n := t.Magnitude()
	if n == 0 {
		return t
	}
	return t.Div(n)
}

func (t Tuple) Dot(r Tuple) float64 {
	res := 0.0
	for i, _ := range t {
		res = res + t[i]*r[i]
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

func (t Tuple) Reflect(normal Tuple) Tuple {
	if !normal.IsVector() || !t.IsVector() {
		panic("Reflection is between two vectors")
	}
	return t.Subtract(normal.Mult(2.0).Mult(t.Dot(normal)))
}
