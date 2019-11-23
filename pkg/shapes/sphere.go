package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type sphere struct {
}

func NewSphere() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), sphere{})
}

func (s sphere) shapeIdPrefix() string {
	return "S"
}

func (s sphere) normalAt(point tuple.Tuple) tuple.Tuple {
	return point.Subtract(tuple.NewPoint(0, 0, 0))
}

func (s sphere) localIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64 {

	sr := origin.Subtract(tuple.NewPoint(0, 0, 0))
	a := direction.Dot(direction)
	b := 2 * direction.Dot(sr)
	c := sr.Dot(sr) - 1.0

	// Solve "a*t^2 + b*t + c" for t to get the intersections
	disc := b*b - 4*a*c
	if disc < 0 {
		return []float64{}
	}
	rootOfDisc := math.Sqrt(disc)
	return []float64{
		(-b - rootOfDisc) / (2 * a),
		(-b + rootOfDisc) / (2 * a),
	}
}
