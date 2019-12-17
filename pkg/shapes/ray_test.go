package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

func TestNewRay(t *testing.T) {
	g := NewGomegaWithT(t)

	_, e := NewRay(tuple.NewVector(1, 2, 3), tuple.NewVector(1, 2, 3))
	g.Expect(e).ToNot(BeNil())
	_, e = NewRay(tuple.NewPoint(1, 2, 3), tuple.NewPoint(1, 2, 3))
	g.Expect(e).ToNot(BeNil())
	r, e := NewRay(tuple.NewPoint(1, 2, 3), tuple.NewVector(1, 2, 3))
	g.Expect(e).To(BeNil())
	g.Expect(r.Origin.Equals(tuple.NewPoint(1, 2, 3))).To(BeTrue())
	g.Expect(r.Direction.Equals(tuple.NewVector(1, 2, 3))).To(BeTrue())
}

func TestRayPosition(t *testing.T) {
	g := NewGomegaWithT(t)

	r, e := NewRay(tuple.NewPoint(2, 3, 4), tuple.NewVector(1, 0, 0))
	g.Expect(e).To(BeNil())
	p := r.Position(0)
	g.Expect(p.Equals(tuple.NewPoint(2, 3, 4))).To(BeTrue())
	p = r.Position(1)
	g.Expect(p.Equals(tuple.NewPoint(3, 3, 4))).To(BeTrue())
	p = r.Position(-1)
	g.Expect(p.Equals(tuple.NewPoint(1, 3, 4))).To(BeTrue())
	p = r.Position(2.5)
	g.Expect(p.Equals(tuple.NewPoint(4.5, 3, 4))).To(BeTrue())
}

func TestRaySphereIntersect(t *testing.T) {
	g := NewGomegaWithT(t)

	r, e := NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	s := NewSphere()

	xs := r.Intersect(s)
	g.Expect(len(xs)).To(Equal(2))
	g.Expect(utils.FloatEqual(xs[0].T, 4.0)).To(BeTrue())
	g.Expect(utils.FloatEqual(xs[1].T, 6.0)).To(BeTrue())

	r, e = NewRay(tuple.NewPoint(0, 1, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	xs = r.Intersect(s)
	g.Expect(len(xs)).To(Equal(2))
	g.Expect(utils.FloatEqual(xs[0].T, 5.0)).To(BeTrue())
	g.Expect(utils.FloatEqual(xs[1].T, 5.0)).To(BeTrue())

	r, e = NewRay(tuple.NewPoint(0, 2, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	xs = r.Intersect(s)
	g.Expect(len(xs)).To(Equal(0))

	r, e = NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	xs = r.Intersect(s)
	g.Expect(len(xs)).To(Equal(2))
	g.Expect(utils.FloatEqual(xs[0].T, -1.0)).To(BeTrue())
	g.Expect(utils.FloatEqual(xs[1].T, 1.0)).To(BeTrue())

	r, e = NewRay(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	xs = r.Intersect(s)
	g.Expect(len(xs)).To(Equal(2))
	g.Expect(utils.FloatEqual(xs[0].T, -6.0)).To(BeTrue())
	g.Expect(utils.FloatEqual(xs[1].T, -4.0)).To(BeTrue())
}

func TestTransformRay(t *testing.T) {
	g := NewGomegaWithT(t)

	r, e := NewRay(tuple.NewPoint(1, 2, 3), tuple.NewVector(0, 1, 0))
	g.Expect(e).To(BeNil())
	tr := r.Transform(matrix.NewTranslation(3, 4, 5))
	g.Expect(tr.Origin.Equals(tuple.NewPoint(4, 6, 8))).To(BeTrue())
	g.Expect(tr.Direction.Equals(tuple.NewVector(0, 1, 0))).To(BeTrue())

	tr = r.Transform(matrix.NewScale(2, 3, 4))
	g.Expect(tr.Origin.Equals(tuple.NewPoint(2, 6, 12))).To(BeTrue())
	g.Expect(tr.Direction.Equals(tuple.NewVector(0, 3, 0))).To(BeTrue())
}
