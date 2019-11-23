package shapes

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestUniqueSpheres(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()
	s2 := NewSphere()

	g.Expect(s1.ID()).ToNot(Equal(s2.ID()))
}

func TestShapeTransform(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()
	g.Expect(s1.GetTransform()).To(Equal(matrix.NewIdentity()))

	s2 := s1.WithTransform(matrix.NewTranslation(2, 3, 4))
	g.Expect(s2.GetTransform()).To(Equal(matrix.NewTranslation(2, 3, 4)))
}

func TestSphereNormal(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()

	n := s1.NormalAt(tuple.NewPoint(1, 0, 0))
	g.Expect(n).To(Equal(tuple.NewVector(1, 0, 0)))

	n = s1.NormalAt(tuple.NewPoint(0, 1, 0))
	g.Expect(n).To(Equal(tuple.NewVector(0, 1, 0)))

	n = s1.NormalAt(tuple.NewPoint(0, 0, 1))
	g.Expect(n).To(Equal(tuple.NewVector(0, 0, 1)))

	v := math.Sqrt(3.0) / 3.0
	n = s1.NormalAt(tuple.NewPoint(v, v, v))
	g.Expect(n).To(Equal(tuple.NewVector(v, v, v)))

	g.Expect(n.Normalize().Equals(n)).To(BeTrue())
}

func TestTransformedSphereNormal(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere().WithTransform(matrix.NewTranslation(0, 1, 0))

	n := s1.NormalAt(tuple.NewPoint(0, 1.70711, -0.70711))
	g.Expect(n.Equals(tuple.NewVector(0, 0.70711, -0.70711))).To(BeTrue())

	s1 = NewSphere().WithTransform(matrix.NewScale(1, 0.5, 1).RotateZ(math.Pi / 5))
	v := math.Sqrt(2.0) / 2
	n = s1.NormalAt(tuple.NewPoint(0, v, -v))
	g.Expect(n.Equals(tuple.NewVector(0, 0.97014, -0.24254))).To(BeTrue())

}
