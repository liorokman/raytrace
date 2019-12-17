package shapes

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestGroups(t *testing.T) {

	g := NewGomegaWithT(t)

	group := NewGroup()

	s := NewSphere()
	g.Expect(s.Parent()).To(BeNil())

	s, err := Connect(group, s)
	g.Expect(err).To(BeNil())

	g.Expect(group.InnerShape().(Group).Size()).To(Equal(1))
	g.Expect(s.Parent()).To(Equal(group))
}

func TestIntersects(t *testing.T) {
	g := NewGomegaWithT(t)

	group := NewGroup()

	r, err := NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs := group.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(0))

	s1 := NewSphere()
	s2 := NewSphere().WithTransform(matrix.NewTranslation(0, 0, -3))
	s3 := NewSphere().WithTransform(matrix.NewTranslation(5, 0, 0))

	s1, err = Connect(group, s1)
	g.Expect(err).To(BeNil())
	s2, err = Connect(group, s2)
	g.Expect(err).To(BeNil())
	s3, err = Connect(group, s3)
	g.Expect(err).To(BeNil())

	r, err = NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	xs = group.LocalIntersect(r)
	g.Expect(len(xs)).To(Equal(4))
	g.Expect(xs[0].Shape.ID()).To(Equal(s2.ID()))
	g.Expect(xs[1].Shape.ID()).To(Equal(s2.ID()))
	g.Expect(xs[2].Shape.ID()).To(Equal(s1.ID()))
	g.Expect(xs[3].Shape.ID()).To(Equal(s1.ID()))
}
