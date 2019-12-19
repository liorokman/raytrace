package shapes

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestGroupWorldToObject(t *testing.T) {
	g := NewGomegaWithT(t)
	group1 := NewGroup().WithTransform(matrix.NewRotateY(math.Pi / 2.0))
	group2 := NewGroup().WithTransform(matrix.NewScale(2, 2, 2))
	group2, err := Connect(group1, group2)
	g.Expect(err).To(BeNil())
	g.Expect(group2.Parent().ID()).To(Equal(group1.ID()))
	sphere := NewSphere().WithTransform(matrix.NewTranslation(5, 0, 0))
	sphere, err = Connect(group2, sphere)
	g.Expect(err).To(BeNil())
	g.Expect(sphere.Parent().ID()).To(Equal(group2.ID()))

	p, err := sphere.WorldToObject(tuple.NewPoint(-2, 0, -10))
	g.Expect(err).To(BeNil())
	g.Expect(p.Equals(tuple.NewPoint(0, 0, -1))).To(BeTrue())

	v, err := sphere.NormalToWorld(tuple.NewVector(math.Sqrt(3.0)/3.0, math.Sqrt(3.0)/3.0, math.Sqrt(3.0)/3.0))
	g.Expect(err).To(BeNil())
	g.Expect(v.Equals(tuple.NewVector(0.2857, 0.4286, -0.8571)))

}

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
