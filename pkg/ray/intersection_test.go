package ray

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
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

func TestPrepareReflectVector(t *testing.T) {
	g := NewGomegaWithT(t)

	r, err := New(tuple.NewPoint(0, 1, -1), tuple.NewVector(0, -math.Sqrt(2.0)/2.0, math.Sqrt(2)/2.0))
	g.Expect(err).To(BeNil())

	i := Intersection{
		T:     math.Sqrt(2),
		Shape: shapes.NewPlane(),
	}
	comps := i.PrepareComputation(r)
	fmt.Printf("Comps: %#v\n", comps)
	g.Expect(comps.ReflectV.Equals(tuple.NewVector(0, math.Sqrt(2)/2, math.Sqrt(2)/2))).To(BeTrue())
}

func TestOverZ(t *testing.T) {
	g := NewGomegaWithT(t)

	r, e := New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	s := shapes.NewSphere().WithTransform(matrix.NewTranslation(0, 0, 1))
	i := Intersection{5, s}

	comps := i.PrepareComputation(r)
	g.Expect(comps.OverPoint.Z()).To(BeNumerically("<", -utils.EPSILON/2))
	g.Expect(comps.Point.Z()).To(BeNumerically(">", comps.OverPoint.Z()))

}
