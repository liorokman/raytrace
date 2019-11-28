package world

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/ray"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func defaultWorld() *World {
	w := New()

	mBuilder := material.NewDefaultBuilder()
	mBuilder.
		WithColor(tuple.NewColor(0.8, 1, 0.6)).
		WithDiffuse(0.7).
		WithSpecular(0.2)
	w.AddShapes(
		shapes.NewSphere().WithMaterial(mBuilder.Build()),
		shapes.NewSphere().WithTransform(matrix.NewScale(0.5, 0.5, 0.5)),
	)
	return w
}

func TestDefaultWorld(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()
	g.Expect(w.NumObjects()).To(Equal(2))
}

func TestIntersectWorld(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	r, e := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	xs := w.IntersectRay(r)
	g.Expect(len(xs)).To(Equal(4))
	vals := []float64{4, 4.5, 5.5, 6}
	for i := range xs {
		g.Expect(xs[i].T).To(Equal(vals[i]))
	}
}

func TestShadeWorld(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	r, e := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	i := ray.Intersection{4, w.Shape(0)}
	comps := i.PrepareComputation(r)
	c := w.ShadeHit(comps)
	g.Expect(c.Equals(tuple.NewColor(0.38066, 0.47583, 0.2855))).To(BeTrue())

	w.Lights[0] = fixtures.NewPointLight(tuple.NewPoint(0, 0.25, 0), tuple.NewColor(1, 1, 1))
	r, e = ray.New(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i = ray.Intersection{0.5, w.Shape(1)}
	comps = i.PrepareComputation(r)
	c = w.ShadeHit(comps)
	g.Expect(c.Equals(tuple.NewColor(0.90498, 0.90498, 0.90498))).To(BeTrue())
}

func TestColorAt(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	// Miss
	r, e := ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0))
	g.Expect(e).To(BeNil())
	c := w.ColorAt(r)
	g.Expect(c).To(Equal(tuple.NewColor(0, 0, 0)))

	// Hit
	r, e = ray.New(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	c = w.ColorAt(r)
	g.Expect(c.Equals(tuple.NewColor(0.38066, 0.47583, 0.2855))).To(BeTrue())

	// Intersection is behind the ray
	for j := 0; j < w.NumObjects(); j++ {
		s := w.Shape(j)
		b := material.NewBuilder(s.GetMaterial())
		b.WithAmbient(1)
		w.SetShape(j, s.WithMaterial(b.Build()))
	}
	r, e = ray.New(tuple.NewPoint(0, 0, 0.75), tuple.NewVector(0, 0, -1))
	g.Expect(e).To(BeNil())
	c = w.ColorAt(r)
	g.Expect(c.Equals(w.Shape(1).GetMaterial().Pattern.ColorAt(tuple.NewPoint(0, 0, 0.75)))).To(BeTrue())
}

func TestIsShadowed(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	g.Expect(w.IsShadowed(tuple.NewPoint(0, 10, 0), 0)).To(BeFalse())
	g.Expect(w.IsShadowed(tuple.NewPoint(10, -10, 10), 0)).To(BeTrue())
	g.Expect(w.IsShadowed(tuple.NewPoint(-20, 20, -20), 0)).To(BeFalse())
	g.Expect(w.IsShadowed(tuple.NewPoint(-2, 2, -2), 0)).To(BeFalse())
}

func TestShadeHit(t *testing.T) {

	g := NewGomegaWithT(t)
	w := New()

	w.Lights[0] = fixtures.NewPointLight(tuple.NewPoint(0, 0, -10), tuple.NewColor(1, 1, 1))

	w.AddShapes(shapes.NewSphere(),
		shapes.NewSphere().WithTransform(matrix.NewTranslation(0, 0, 10)))

	r, err := ray.New(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	i := ray.Intersection{4, w.objects[1]}

	comps := i.PrepareComputation(r)
	c := w.ShadeHit(comps)

	g.Expect(c.Equals(tuple.NewColor(0.1, 0.1, 0.1))).To(BeTrue())

}
