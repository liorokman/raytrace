package ray

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestHit(t *testing.T) {
	g := NewGomegaWithT(t)

	s := shapes.NewSphere()

	i, ok := Hit(Intersection{-1, s}, Intersection{1, s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{1, s}))

	i, ok = Hit(Intersection{1, s}, Intersection{2, s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{1, s}))

	i, ok = Hit(Intersection{-1, s}, Intersection{-2, s})
	g.Expect(ok).To(BeFalse())

	i, ok = Hit(
		Intersection{5, s},
		Intersection{7, s},
		Intersection{-3, s},
		Intersection{2, s},
	)
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{2, s}))
}

func TestPrecomputation(t *testing.T) {

	g := NewGomegaWithT(t)

	r, e := New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i := Intersection{
		T:     4,
		Shape: shapes.NewSphere(),
	}
	c := i.PrepareComputation(r)

	g.Expect(c.T).To(Equal(i.T))
	g.Expect(c.Shape).To(Equal(i.Shape))
	g.Expect(c.Point).To(Equal(tuple.NewPoint(0, 0, -1)))
	g.Expect(c.EyeV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.NormalV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.Inside).To(BeFalse())

	r, e = New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i = Intersection{
		T:     1,
		Shape: shapes.NewSphere(),
	}
	c = i.PrepareComputation(r)
	g.Expect(c.Point).To(Equal(tuple.NewPoint(0, 0, 1)))
	g.Expect(c.EyeV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.NormalV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.Inside).To(BeTrue())
}
