package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestCylinderNormal(t *testing.T) {
	tests := []struct {
		point  tuple.Tuple
		normal tuple.Tuple
	}{
		{tuple.NewPoint(1, 0, 0), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(0, 5, -1), tuple.NewVector(0, 0, -1)},
		{tuple.NewPoint(0, -2, 1), tuple.NewVector(0, 0, 1)},
		{tuple.NewPoint(-1, 1, 0), tuple.NewVector(-1, 0, 0)},
	}

	g := NewGomegaWithT(t)
	c := NewCylinder()

	for _, curr := range tests {

		n, err := c.NormalAt(curr.point)
		g.Expect(err).To(BeNil())
		g.Expect(n).To(Equal(curr.normal))
	}
}

func TestClosedCylinderNormal(t *testing.T) {
	tests := []struct {
		point  tuple.Tuple
		normal tuple.Tuple
	}{
		{tuple.NewPoint(0, 1, 0), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0.5, 1, 0), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0, 1, 0.5), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(0, 2, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0.5, 2, 0), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0, 2, 0.5), tuple.NewVector(0, 1, 0)},
	}

	g := NewGomegaWithT(t)
	c := NewConstrainedCylinder(1, 2, true)

	for _, curr := range tests {
		n, err := c.NormalAt(curr.point)
		g.Expect(err).To(BeNil())
		g.Expect(n).To(Equal(curr.normal))
	}
}

func TestCylinderIntersections(t *testing.T) {

	tests := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		t         []float64
	}{
		// No hits
		{tuple.NewPoint(1, 0, 0), tuple.NewVector(0, 1, 0), []float64{}},
		{tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0), []float64{}},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(1, 1, 1), []float64{}},

		// Intersects
		{tuple.NewPoint(1, 0, -5), tuple.NewVector(0, 0, 1), []float64{5, 5}},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), []float64{4, 6}},
		{tuple.NewPoint(0.5, 0, -5), tuple.NewVector(0.1, 1, 1), []float64{6.80798191702732, 7.088723439378861}},
	}
	g := NewGomegaWithT(t)
	c := NewCylinder()

	for _, curr := range tests {

		r, err := NewRay(curr.origin, curr.direction.Normalize())
		g.Expect(err).To(BeNil())
		xs := c.LocalIntersect(r)

		g.Expect(len(xs)).To(Equal(len(curr.t)))
		for i := range xs {
			g.Expect(xs[i].T).To(BeNumerically("~", curr.t[i]))
		}

	}
}

func TestConstrainedCylinderIntersections(t *testing.T) {

	tests := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		num       int
	}{
		{tuple.NewPoint(0, 1.5, 0), tuple.NewVector(0.1, 1, 0), 0},
		{tuple.NewPoint(0, 3, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1), 0},
		{tuple.NewPoint(0, 1.5, -2), tuple.NewVector(0, 0, 1), 2},
	}
	g := NewGomegaWithT(t)
	c := NewConstrainedCylinder(1, 2, false)

	for _, curr := range tests {

		r, err := NewRay(curr.origin, curr.direction.Normalize())
		g.Expect(err).To(BeNil())
		xs := c.LocalIntersect(r)

		g.Expect(len(xs)).To(Equal(curr.num))

	}
}

func TestConstrainedClosedCylinderIntersections(t *testing.T) {

	tests := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		num       int
	}{
		{tuple.NewPoint(0, 3, 0), tuple.NewVector(0, -1, 0), 2},
		{tuple.NewPoint(0, 3, -2), tuple.NewVector(0, -1, 2), 2},
		{tuple.NewPoint(0, 4, -2), tuple.NewVector(0, -1, 1), 2},
		{tuple.NewPoint(0, 0, -2), tuple.NewVector(0, 1, 2), 2},
		{tuple.NewPoint(0, -1, -2), tuple.NewVector(0, 1, 1), 2},
	}
	g := NewGomegaWithT(t)
	c := NewConstrainedCylinder(1, 2, true)

	for _, curr := range tests {

		r, err := NewRay(curr.origin, curr.direction.Normalize())
		g.Expect(err).To(BeNil())
		xs := c.LocalIntersect(r)

		g.Expect(len(xs)).To(Equal(curr.num))

	}
}
