package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

func defaultSmoothTriangle() Shape {
	return NewSmoothTriangle(tuple.NewPoint(0, 1, 0), tuple.NewPoint(-1, 0, 0), tuple.NewPoint(1, 0, 0),
		tuple.NewVector(0, 1, 0), tuple.NewVector(-1, 0, 0), tuple.NewVector(1, 0, 0))
}

func TestSmoothTriangleIntersection(t *testing.T) {
	g := NewGomegaWithT(t)

	tri := defaultSmoothTriangle()

	r, err := NewRay(tuple.NewPoint(-0.2, 0.3, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs := tri.LocalIntersect(r)
	g.Expect(utils.FloatEqual(xs[0].U, 0.45)).To(BeTrue())
	g.Expect(utils.FloatEqual(xs[0].V, 0.25)).To(BeTrue())

	n, err := tri.NormalAt(tuple.NewPoint(0, 0, 0), xs[0])
	g.Expect(err).To(BeNil())
	g.Expect(n.Equals(tuple.NewVector(-0.5547, 0.83205, 0))).To(BeTrue())

	xs = []Intersection{Intersection{T: 1, Shape: tri, U: 0.45, V: 0.25}}
	comps, err := xs[0].PrepareComputation(r, xs...)
	g.Expect(err).To(BeNil())
	g.Expect(comps.NormalV.Equals(tuple.NewVector(-0.5547, 0.83205, 0))).To(BeTrue())
}
