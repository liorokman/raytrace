package shapes

import (
	"testing"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	. "github.com/onsi/gomega"
)

func TestNewCSG(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()
	s2 := NewSphere()

	g.Expect(s1.Parent()).To(BeNil())
	g.Expect(s2.Parent()).To(BeNil())
	csg := NewCSG(&s1, &s2, UnionOp)
	g.Expect(s2.Parent()).To(Equal(csg))
	g.Expect(s1.Parent()).To(Equal(csg))
}

func TestIntersectionAllowed(t *testing.T) {
	tests := []struct {
		op             CSGOp
		lhit, inL, inR bool
		result         bool
	}{
		{UnionOp, true, true, true, false},
		{UnionOp, true, true, false, true},
		{UnionOp, true, false, true, false},
		{UnionOp, true, false, false, true},
		{UnionOp, false, true, true, false},
		{UnionOp, false, true, false, false},
		{UnionOp, false, false, true, true},
		{UnionOp, false, false, false, true},

		{IntersectOp, true, true, true, true},
		{IntersectOp, true, true, false, false},
		{IntersectOp, true, false, true, true},
		{IntersectOp, true, false, false, false},
		{IntersectOp, false, true, true, true},
		{IntersectOp, false, true, false, true},
		{IntersectOp, false, false, true, false},
		{IntersectOp, false, false, false, false},

		{DifferenceOp, true, true, true, false},
		{DifferenceOp, true, true, false, true},
		{DifferenceOp, true, false, true, false},
		{DifferenceOp, true, false, false, true},
		{DifferenceOp, false, true, true, true},
		{DifferenceOp, false, true, false, true},
		{DifferenceOp, false, false, true, false},
		{DifferenceOp, false, false, false, false},
	}

	g := NewGomegaWithT(t)
	s1 := NewSphere()
	s2 := NewSphere()

	for _, test := range tests {
		c := NewCSG(&s1, &s2, test.op)
		r := c.InnerShape().(csg).intersectionAllowed(test.lhit, test.inL, test.inR)
		g.Expect(r).To(Equal(test.result))
	}
}

func TestFilterIntersections(t *testing.T) {
	tests := []struct {
		op     CSGOp
		x0, x1 int
	}{
		{UnionOp, 0, 3},
		{IntersectOp, 1, 2},
		{DifferenceOp, 0, 1},
	}

	s1 := NewSphere()
	s2 := NewCube()
	xs := []Intersection{
		{T: 1, Shape: s1},
		{T: 2, Shape: s2},
		{T: 3, Shape: s1},
		{T: 4, Shape: s2},
	}

	g := NewGomegaWithT(t)

	for _, test := range tests {
		c := NewCSG(&s1, &s2, test.op)
		result := c.InnerShape().(csg).filterIntersections(xs)
		g.Expect(len(result)).To(Equal(2))
		g.Expect(result[0]).To(Equal(xs[test.x0]))
		g.Expect(result[1]).To(Equal(xs[test.x1]))
	}
}

func TestCSGIntersections(t *testing.T) {

	g := NewGomegaWithT(t)

	s1, s2 := NewSphere(), NewCube()

	// test misses
	c := NewCSG(&s1, &s2, UnionOp)
	r, err := NewRay(tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())

	xs := r.Intersect(c)
	g.Expect(len(xs)).To(Equal(0))

	// test hits
	s2 = NewSphere().WithTransform(matrix.NewTranslation(0, 0, 0.5))
	c = NewCSG(&s1, &s2, UnionOp)
	r, err = NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = r.Intersect(c)
	g.Expect(len(xs)).To(Equal(2))
	g.Expect(xs[0].T).To(Equal(4.0))
	g.Expect(xs[0].Shape.ID()).To(Equal(s1.ID()))
	g.Expect(xs[1].T).To(Equal(6.5))
	g.Expect(xs[1].Shape.ID()).To(Equal(s2.ID()))

}
