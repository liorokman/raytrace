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

func NewGlassSphere() Shape {
	return newShape(material.Glass(), matrix.NewIdentity(), sphere{})
}

func (s sphere) shapeIdPrefix() string {
	return "S"
}

func (s sphere) normalAt(point tuple.Tuple) tuple.Tuple {
	return point.Subtract(tuple.NewPoint(0, 0, 0))
}

func (s sphere) localIntersect(ray Ray, outer Shape) []Intersection {

	sr := ray.Origin.Subtract(tuple.NewPoint(0, 0, 0))
	a := ray.Direction.Dot(ray.Direction)
	b := 2.0 * ray.Direction.Dot(sr)
	c := sr.Dot(sr) - 1.0

	// Solve "a*t^2 + b*t + c" for t to get the intersections
	disc := b*b - 4.0*a*c
	if disc < 0 {
		return []Intersection{}
	}
	rootOfDisc := math.Sqrt(disc)
	return []Intersection{
		{(-b - rootOfDisc) / (2.0 * a), outer},
		{(-b + rootOfDisc) / (2.0 * a), outer},
	}
}
