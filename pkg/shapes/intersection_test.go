package shapes

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

func TestHit(t *testing.T) {
	g := NewGomegaWithT(t)

	s := NewSphere()

	i, ok := Hit(Intersection{T: -1, Shape: s}, Intersection{T: 1, Shape: s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{T: 1, Shape: s}))

	i, ok = Hit(Intersection{T: 1, Shape: s}, Intersection{T: 2, Shape: s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{T: 1, Shape: s}))

	i, ok = Hit(Intersection{T: -1, Shape: s}, Intersection{T: -2, Shape: s})
	g.Expect(ok).To(BeFalse())

	i, ok = Hit(
		Intersection{T: 5, Shape: s},
		Intersection{T: 7, Shape: s},
		Intersection{T: -3, Shape: s},
		Intersection{T: 2, Shape: s},
	)
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{T: 2, Shape: s}))
}

func TestPrecomputation(t *testing.T) {

	g := NewGomegaWithT(t)

	r, e := NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i := Intersection{
		T:     4,
		Shape: NewSphere(),
	}
	c, e := i.PrepareComputation(r, i)
	g.Expect(e).To(BeNil())

	g.Expect(c.T).To(Equal(i.T))
	g.Expect(c.Shape).To(Equal(i.Shape))
	g.Expect(c.Point).To(Equal(tuple.NewPoint(0, 0, -1)))
	g.Expect(c.EyeV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.NormalV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.Inside).To(BeFalse())

	r, e = NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i = Intersection{
		T:     1,
		Shape: NewSphere(),
	}
	c, e = i.PrepareComputation(r, i)
	g.Expect(e).To(BeNil())
	g.Expect(c.Point).To(Equal(tuple.NewPoint(0, 0, 1)))
	g.Expect(c.EyeV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.NormalV).To(Equal(tuple.NewVector(0, 0, -1)))
	g.Expect(c.Inside).To(BeTrue())
}

func TestPrepareReflectVector(t *testing.T) {
	g := NewGomegaWithT(t)

	r, err := NewRay(tuple.NewPoint(0, 1, -1), tuple.NewVector(0, -math.Sqrt(2.0)/2.0, math.Sqrt(2)/2.0))
	g.Expect(err).To(BeNil())

	i := Intersection{
		T:     math.Sqrt(2),
		Shape: NewPlane(),
	}
	comps, err := i.PrepareComputation(r, i)
	g.Expect(err).To(BeNil())
	g.Expect(comps.ReflectV.Equals(tuple.NewVector(0, math.Sqrt(2)/2, math.Sqrt(2)/2))).To(BeTrue())
}

func TestOverUnderZ(t *testing.T) {
	g := NewGomegaWithT(t)

	r, e := NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	s := NewSphere().WithTransform(matrix.NewTranslation(0, 0, 1))
	i := Intersection{T: 5, Shape: s}

	comps, e := i.PrepareComputation(r, i)
	g.Expect(e).To(BeNil())
	g.Expect(comps.OverPoint.Z()).To(BeNumerically("<", -utils.EPSILON/2))
	g.Expect(comps.Point.Z()).To(BeNumerically(">", comps.OverPoint.Z()))

	g.Expect(comps.UnderPoint.Z()).To(BeNumerically(">", utils.EPSILON/2))
	g.Expect(comps.Point.Z()).To(BeNumerically("<", comps.UnderPoint.Z()))

}

func TestFindingN1N2(t *testing.T) {
	g := NewGomegaWithT(t)
	mb := material.NewDefaultBuilder()
	A := NewGlassSphere().WithTransform(matrix.NewScale(2, 2, 2)).WithMaterial(mb.WithRefractiveIndex(1.5).Build())
	B := NewGlassSphere().WithTransform(matrix.NewTranslation(0, 0, -0.25)).WithMaterial(mb.WithRefractiveIndex(2.0).Build())
	C := NewGlassSphere().WithTransform(matrix.NewTranslation(0, 0, 0.25)).WithMaterial(mb.WithRefractiveIndex(2.5).Build())
	r, err := NewRay(tuple.NewPoint(0, 0, -4), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs := []Intersection{{T: 2, Shape: A}, {T: 2.75, Shape: B}, {T: 3.25, Shape: C}, {T: 4.75, Shape: B}, {T: 5.25, Shape: C}, {T: 6, Shape: A}}
	expected := []struct{ n1, n2 float64 }{{1.0, 1.5}, {1.5, 2.0}, {2.0, 2.5}, {2.5, 2.5}, {2.5, 1.5}, {1.5, 1.0}}
	for i := range xs {
		comps, err := xs[i].PrepareComputation(r, xs...)
		g.Expect(err).To(BeNil())
		g.Expect(comps.N1).To(BeNumerically("==", expected[i].n1))
		g.Expect(comps.N2).To(BeNumerically("==", expected[i].n2))
	}
}

func TestSchlick(t *testing.T) {
	g := NewGomegaWithT(t)

	// Total internal reflection
	sphere := NewGlassSphere()
	r, err := NewRay(tuple.NewPoint(0, 0, math.Sqrt(2)/2.0), tuple.NewVector(0, 1, 0))
	g.Expect(err).To(BeNil())
	xs := []Intersection{
		{T: -math.Sqrt(2) / 2, Shape: sphere},
		{T: math.Sqrt(2) / 2, Shape: sphere},
	}
	comps, err := xs[1].PrepareComputation(r, xs...)
	g.Expect(err).To(BeNil())
	reflectence := comps.Schlick()
	g.Expect(reflectence).To(Equal(1.0))

	// Reflectence of perpendicular viewing angle
	r, err = NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0))
	g.Expect(err).To(BeNil())
	xs = []Intersection{{T: -1, Shape: sphere}, {T: 1, Shape: sphere}}
	comps, err = xs[1].PrepareComputation(r, xs...)
	g.Expect(err).To(BeNil())
	reflectence = comps.Schlick()
	g.Expect(reflectence).To(BeNumerically("~", 0.04))

	// Reflectence when n2 > n1
	r, err = NewRay(tuple.NewPoint(0, 0.99, -2), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = []Intersection{{T: 1.8589, Shape: sphere}}
	comps, err = xs[0].PrepareComputation(r, xs...)
	g.Expect(err).To(BeNil())
	reflectence = comps.Schlick()
	g.Expect(reflectence).To(BeNumerically("~", 0.4887308101))

}
