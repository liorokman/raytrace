package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
)

func TestUniqueSpheres(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()
	s2 := NewSphere()

	g.Expect(s1.ID()).ToNot(Equal(s2.ID()))
}

func TestSphereTransform(t *testing.T) {
	g := NewGomegaWithT(t)

	s1 := NewSphere()
	g.Expect(s1.GetTransform()).To(Equal(matrix.NewIdentity()))

	s2 := s1.WithTransform(matrix.NewTranslation(2, 3, 4))
	g.Expect(s2.GetTransform()).To(Equal(matrix.NewTranslation(2, 3, 4)))
}
