package shapes

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type smoothTriangle struct {
	P1, P2, P3 tuple.Tuple
	E1, E2     tuple.Tuple
	N1, N2, N3 tuple.Tuple
}

func NewSmoothTriangle(p1, p2, p3, n1, n2, n3 tuple.Tuple) Shape {
	return newShape(material.Default(), matrix.NewIdentity(), smoothTriangle{
		P1: p1,
		P2: p2,
		P3: p3,
		N1: n1,
		N2: n2,
		N3: n3,
		E1: p2.Subtract(p1),
		E2: p3.Subtract(p1),
	})
}

func (t smoothTriangle) String() string {
	return fmt.Sprintf("P1: %s, P2: %s, P3: %s, N1: %s, N2: %s, N3: %s", t.P1, t.P2, t.P3, t.N1, t.N2, t.N3)
}

func (t smoothTriangle) shapeIdPrefix() string {
	return "ST"
}

func (t smoothTriangle) normalAt(point tuple.Tuple, hit Intersection) tuple.Tuple {
	return t.N2.Mult(hit.U).Add(t.N3.Mult(hit.V)).Add(t.N1.Mult(1.0 - hit.U - hit.V))
}

func (t smoothTriangle) localIntersect(ray Ray, outer Shape) []Intersection {
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

	return []Intersection{
		{
			T:     f * t.E2.Dot(originCrossE1),
			Shape: outer,
			U:     u,
			V:     v,
		},
	}

}
