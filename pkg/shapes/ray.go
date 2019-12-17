package shapes

import (
	"fmt"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Ray struct {
	Origin    tuple.Tuple
	Direction tuple.Tuple
}

func NewRay(origin, direction tuple.Tuple) (Ray, error) {
	if !origin.IsPoint() {
		return Ray{}, fmt.Errorf("Origin is not a point")
	}
	if !direction.IsVector() {
		return Ray{}, fmt.Errorf("Direction is not a vector")
	}
	return Ray{origin, direction}, nil
}

func (r Ray) Position(time float64) tuple.Tuple {
	return r.Origin.Add(r.Direction.Mult(time))
}

func (r Ray) Intersect(shape Shape) []Intersection {
	invShapeTransform, err := shape.GetTransform().Inverse()
	if err != nil {
		panic(err)
	}
	tr := r.Transform(invShapeTransform)
	ints := shape.LocalIntersect(tr)

	retval := make([]Intersection, len(ints))
	for i := range ints {
		retval[i] = Intersection{ints[i].T, ints[i].Shape}
	}
	return retval
}

func (r Ray) Transform(m matrix.Matrix) Ray {

	return Ray{
		Origin:    m.MultiplyTuple(r.Origin),
		Direction: m.MultiplyTuple(r.Direction),
	}
}
