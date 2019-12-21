package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestCubeNormal(t *testing.T) {
	tests := []struct {
		point  tuple.Tuple
		normal tuple.Tuple
	}{
		{tuple.NewPoint(1, 0.5, -0.8), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(-1, -0.2, 0.9), tuple.NewVector(-1, 0, 0)},
		{tuple.NewPoint(-0.4, 1, -0.1), tuple.NewVector(0, 1, 0)},
		{tuple.NewPoint(0.3, -1, -0.7), tuple.NewVector(0, -1, 0)},
		{tuple.NewPoint(-0.6, 0.3, 1), tuple.NewVector(0, 0, 1)},
		{tuple.NewPoint(0.4, 0.4, -1), tuple.NewVector(0, 0, -1)},
		{tuple.NewPoint(1, 1, 1), tuple.NewVector(1, 0, 0)},
		{tuple.NewPoint(-1, -1, -1), tuple.NewVector(-1, 0, 0)},
	}

	g := NewGomegaWithT(t)
	c := cube{}

	for i := range tests {
		norm := c.normalAt(tests[i].point, Intersection{})
		g.Expect(norm.Equals(tests[i].normal)).To(BeTrue())
	}
}

func TestCubeIntersections(t *testing.T) {

	tests := []struct {
		name      string
		origin    tuple.Tuple
		direction tuple.Tuple
		t         []float64
	}{
		{"+x", tuple.NewPoint(5, 0.5, 0), tuple.NewVector(-1, 0, 0), []float64{4, 6}},
		{"-x", tuple.NewPoint(-5, 0.5, 0), tuple.NewVector(1, 0, 0), []float64{4, 6}},
		{"+y", tuple.NewPoint(0.5, 5, 0), tuple.NewVector(0, -1, 0), []float64{4, 6}},
		{"-y", tuple.NewPoint(0.5, -5, 0), tuple.NewVector(0, 1, 0), []float64{4, 6}},
		{"+z", tuple.NewPoint(0.5, 0, 5), tuple.NewVector(0, 0, -1), []float64{4, 6}},
		{"-z", tuple.NewPoint(0.5, 0, -5), tuple.NewVector(0, 0, 1), []float64{4, 6}},
		{"inside", tuple.NewPoint(0, 0.5, 0), tuple.NewVector(0, 0, 1), []float64{-1, 1}},
		{"outside1", tuple.NewPoint(-2, 0, 0), tuple.NewVector(0.2673, 0.5345, 0.8018), []float64{}},
		{"outside2", tuple.NewPoint(0, -2, 0), tuple.NewVector(0.8018, 0.2673, 0.5345), []float64{}},
		{"outside3", tuple.NewPoint(0, 0, -2), tuple.NewVector(0.5345, 0.8018, 0.2673), []float64{}},
		{"outside5", tuple.NewPoint(2, 0, 2), tuple.NewVector(0, 0, -1), []float64{}},
		{"outside4", tuple.NewPoint(0, 2, 2), tuple.NewVector(0, -1, 0), []float64{}},
		{"outside6", tuple.NewPoint(2, 2, 0), tuple.NewVector(-1, 0, 0), []float64{}},
	}

	g := NewGomegaWithT(t)
	c := NewCube()

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
