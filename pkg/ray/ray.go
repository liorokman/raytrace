package ray

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Ray struct {
	Origin    tuple.Tuple
	Direction tuple.Tuple
}

func New(origin, direction tuple.Tuple) (Ray, error) {
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

func (r Ray) Intersect(shape shapes.Shape) []Intersection {
	// sr is "sphere-to-ray"
	invShapeTransform, err := shape.GetTransform().Inverse()
	if err != nil {
		panic(err)
	}
	tr := r.Transform(invShapeTransform)
	sr := tr.Origin.Subtract(tuple.NewPoint(0, 0, 0))
	a := tr.Direction.Dot(tr.Direction)
	b := 2 * tr.Direction.Dot(sr)
	c := sr.Dot(sr) - 1.0

	// Solve "a*t^2 + b*t + c" for t to get the intersections
	disc := b*b - 4*a*c
	if disc < 0 {
		return []Intersection{}
	}

	rootOfDisc := math.Sqrt(disc)
	return []Intersection{
		{
			T:     (-b - rootOfDisc) / (2 * a),
			Shape: shape,
		},
		{
			T:     (-b + rootOfDisc) / (2 * a),
			Shape: shape,
		},
	}
}

func (r Ray) Transform(m matrix.Matrix) Ray {

	return Ray{
		Origin:    m.MultiplyTuple(r.Origin),
		Direction: m.MultiplyTuple(r.Direction),
	}
}
