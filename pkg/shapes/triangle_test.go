package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestNewTriangle(t *testing.T) {
	g := NewGomegaWithT(t)

	t1 := NewTriangle(tuple.NewPoint(0, 1, 0), tuple.NewPoint(-1, 0, 0), tuple.NewPoint(1, 0, 0))

	g.Expect(t1.InnerShape().(triangle).E1).To(Equal(tuple.NewVector(-1, -1, 0)))
	g.Expect(t1.InnerShape().(triangle).E2).To(Equal(tuple.NewVector(1, -1, 0)))
	g.Expect(t1.InnerShape().(triangle).Normal).To(Equal(tuple.NewVector(0, 0, -1)))

	n, err := t1.NormalAt(tuple.NewPoint(-0.5, 0.75, 0), Intersection{})
	g.Expect(err).To(BeNil())
	g.Expect(t1.InnerShape().(triangle).Normal).To(Equal(n))
}

func TestTriangleIntersect(t *testing.T) {
	g := NewGomegaWithT(t)

	t1 := NewTriangle(tuple.NewPoint(0, 1, 0), tuple.NewPoint(-1, 0, 0), tuple.NewPoint(1, 0, 0))

	// ray is parallel to t1
	r, err := NewRay(tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 1, 0))
	g.Expect(err).To(BeNil())
	xs := t1.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(0))

	// ray misses the edges
	r, err = NewRay(tuple.NewPoint(1, 1, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = t1.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(0))

	r, err = NewRay(tuple.NewPoint(-1, 1, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = t1.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(0))

	r, err = NewRay(tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = t1.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(0))

	// Ray strikes the triangle
	r, err = NewRay(tuple.NewPoint(0, 0.5, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = t1.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(1))
	g.Expect(xs[0].T).To(Equal(2.0))

}
