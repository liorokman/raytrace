package shapes

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestConeNormal(t *testing.T) {
	tests := []struct {
		point  tuple.Tuple
		normal tuple.Tuple
	}{
		{tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 0)},
		{tuple.NewPoint(1, 1, 1), tuple.NewVector(1, -math.Sqrt(2.0), 1)},
		{tuple.NewPoint(-1, -1, 0), tuple.NewVector(-1, 1, 0)},
	}

	g := NewGomegaWithT(t)
	c := NewCone()

	for _, curr := range tests {

		n, err := c.NormalAt(curr.point)
		g.Expect(err).To(BeNil())
		g.Expect(n).To(Equal(curr.normal.Normalize()))
	}

}

func TestConeIntersect(t *testing.T) {
	tests := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		t         []float64
	}{
		// Intersect a cone with a ray
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1), []float64{5, 5}},
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(1, 1, 1), []float64{8.660254037844386, 8.660254037844386}},
		{tuple.NewPoint(1, 1, -5), tuple.NewVector(-0.5, -1, 1), []float64{4.550055679356349, 49.449944320643645}},

		// Intersect a cone with a ray that is parallel to itself
		{tuple.NewPoint(0, 0, -1), tuple.NewVector(0, 1, 1), []float64{0.3535533905932738}},
	}

	g := NewGomegaWithT(t)
	c := NewCone()

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

func TestConstrainedConeIntersections(t *testing.T) {
	tests := []struct {
		origin    tuple.Tuple
		direction tuple.Tuple
		num       int
	}{
		{tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0), 0},
		{tuple.NewPoint(0, 0, -0.25), tuple.NewVector(0, 1, 1), 2},
		{tuple.NewPoint(0, 0, -0.25), tuple.NewVector(0, 1, 0), 4},
	}
	g := NewGomegaWithT(t)
	c := NewConstrainedCone(-0.5, 0.5, true)

	for _, curr := range tests {
		r, err := NewRay(curr.origin, curr.direction.Normalize())
		g.Expect(err).To(BeNil())
		xs := c.LocalIntersect(r)
		g.Expect(len(xs)).To(Equal(curr.num))
	}
}
