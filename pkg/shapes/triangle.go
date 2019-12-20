package shapes

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type triangle struct {
	P1, P2, P3 tuple.Tuple
	E1, E2     tuple.Tuple
	Normal     tuple.Tuple
}

func NewTriangle(p1, p2, p3 tuple.Tuple) Shape {

	t := triangle{
		P1: p1,
		P2: p2,
		P3: p3,
		E1: p2.Subtract(p1),
		E2: p3.Subtract(p1),
	}
	t.Normal = t.E2.Cross(t.E1).Normalize()
	return newShape(material.Default(), matrix.NewIdentity(), t)
}

func (t triangle) String() string {
	return fmt.Sprintf("P1: %s, P2: %s, P3: %s", t.P1, t.P2, t.P3)
}

func (t triangle) shapeIdPrefix() string {
	return "T"
}

func (t triangle) normalAt(point tuple.Tuple) tuple.Tuple {
	return t.Normal
}

func (t triangle) localIntersect(ray Ray, outer Shape) []Intersection {
	dirCrossE2 := ray.Direction.Cross(t.E2)
	det := t.E1.Dot(dirCrossE2)

	if math.Abs(det) < utils.EPSILON {
		return []Intersection{}
	}

	f := 1.0 / det
	p1ToOrigin := ray.Origin.Subtract(t.P1)
	u := f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		return []Intersection{}
	}

	originCrossE1 := p1ToOrigin.Cross(t.E1)
	v := f * ray.Direction.Dot(originCrossE1)
	if v < 0 || (u+v) > 1 {
		return []Intersection{}
	}

	return []Intersection{{f * t.E2.Dot(originCrossE1), outer}}
}
