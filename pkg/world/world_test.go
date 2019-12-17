package world

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
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

func TestReflection(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()
	r, e := shapes.NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	mb := material.NewBuilder(w.Shape(1).GetMaterial())
	mb.WithAmbient(1)
	w.SetShape(1, w.Shape(1).WithMaterial(mb.Build()))
	i := shapes.Intersection{
		T:     1,
		Shape: w.Shape(1),
	}
	comps := i.PrepareComputation(r)
	color := w.ReflectedColor(comps, 5)
	g.Expect(color).To(Equal(tuple.Black))

	mb.Reset().WithReflective(0.5)
	s := shapes.NewPlane().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, -1, 0))
	w.AddShapes(s)
	r, e = shapes.NewRay(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt(2)/2.0, math.Sqrt(2.0)/2.0))
	g.Expect(e).To(BeNil())
	i = shapes.Intersection{
		T:     math.Sqrt(2),
		Shape: s,
	}
	comps = i.PrepareComputation(r)
	color = w.ReflectedColor(comps, 5)
	g.Expect(color.Equals(tuple.NewColor(0.1903323, 0.2379154, 0.14274924))).To(BeTrue())

	color = w.ShadeHit(comps, 5)
	g.Expect(color.Equals(tuple.NewColor(0.8767577, 0.9243407, 0.82917462))).To(BeTrue())
}

func TestNoInfiniteRecursionInReflection(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	mb := material.NewDefaultBuilder().WithReflective(0.5)
	w.AddShapes(shapes.NewPlane().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, -1, 0)))
	r, e := shapes.NewRay(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt(2)/2.0, math.Sqrt(2)/2.0))
	g.Expect(e).To(BeNil())
	i := shapes.Intersection{
		T:     math.Sqrt(2.0),
		Shape: w.Shape(2),
	}
	comps := i.PrepareComputation(r)
	g.Expect(w.ReflectedColor(comps, 0)).To(Equal(tuple.Black))

	w = New()
	w.Lights[0] = fixtures.NewPointLight(tuple.NewPoint(0, 0, 0), tuple.White)
	mb.Reset().WithReflective(1)
	w.AddShapes(shapes.NewPlane().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, -1, 0)))
	w.AddShapes(shapes.NewPlane().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, 1, 0)))
	r, e = shapes.NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 1, 0))
	g.Expect(e).To(BeNil())
	g.Expect(w.ColorAt(r, 4)).To(Equal(tuple.NewColor(9.5, 9.5, 9.5)))
}

func TestIntersectWorld(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	r, e := shapes.NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
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

	r, e := shapes.NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())

	i := shapes.Intersection{T: 4, Shape: w.Shape(0)}
	comps := i.PrepareComputation(r)
	c := w.ShadeHit(comps, 5)
	g.Expect(c.Equals(tuple.NewColor(0.38066, 0.47583, 0.2855))).To(BeTrue())

	w.Lights[0] = fixtures.NewPointLight(tuple.NewPoint(0, 0.25, 0), tuple.NewColor(1, 1, 1))
	r, e = shapes.NewRay(tuple.NewPoint(0, 0, 0), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	i = shapes.Intersection{T: 0.5, Shape: w.Shape(1)}
	comps = i.PrepareComputation(r)
	c = w.ShadeHit(comps, 5)
	g.Expect(c.Equals(tuple.NewColor(0.90498, 0.90498, 0.90498))).To(BeTrue())
}

func TestColorAt(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	// Miss
	r, e := shapes.NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 1, 0))
	g.Expect(e).To(BeNil())
	c := w.ColorAt(r, 5)
	g.Expect(c).To(Equal(tuple.NewColor(0, 0, 0)))

	// Hit
	r, e = shapes.NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(e).To(BeNil())
	c = w.ColorAt(r, 5)
	g.Expect(c.Equals(tuple.NewColor(0.38066, 0.47583, 0.2855))).To(BeTrue())

	// Intersection is behind the ray
	for j := 0; j < w.NumObjects(); j++ {
		s := w.Shape(j)
		b := material.NewBuilder(s.GetMaterial())
		b.WithAmbient(1)
		w.SetShape(j, s.WithMaterial(b.Build()))
	}
	r, e = shapes.NewRay(tuple.NewPoint(0, 0, 0.75), tuple.NewVector(0, 0, -1))
	g.Expect(e).To(BeNil())
	c = w.ColorAt(r, 5)
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

	r, err := shapes.NewRay(tuple.NewPoint(0, 0, 5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())
	i := shapes.Intersection{T: 4, Shape: w.objects[1]}

	comps := i.PrepareComputation(r)
	c := w.ShadeHit(comps, 5)

	g.Expect(c.Equals(tuple.NewColor(0.1, 0.1, 0.1))).To(BeTrue())

}

func TestRefractions(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()

	r, err := shapes.NewRay(tuple.NewPoint(0, 0, -5), tuple.NewVector(0, 0, 1))
	g.Expect(err).To(BeNil())

	xs := []shapes.Intersection{
		{T: 4, Shape: w.Shape(0)},
		{T: 6, Shape: w.Shape(0)},
	}
	comps := xs[0].PrepareComputation(r, xs...)
	c := w.RefractedColor(comps, 5)
	g.Expect(c).To(Equal(tuple.Black))

	w.SetShape(0, shapes.NewGlassSphere())
	xs = []shapes.Intersection{
		{T: 4, Shape: w.Shape(0)},
		{T: 6, Shape: w.Shape(0)},
	}
	comps = xs[0].PrepareComputation(r, xs...)
	c = w.RefractedColor(comps, 0)
	g.Expect(c).To(Equal(tuple.Black))

	// Test total internal reflection
	r, err = shapes.NewRay(tuple.NewPoint(0, 0, math.Sqrt(2.0)/2.0), tuple.NewVector(0, 1, 0))
	g.Expect(err).To(BeNil())
	xs = []shapes.Intersection{
		{T: -math.Sqrt(2.0) / 2.0, Shape: w.Shape(0)},
		{T: math.Sqrt(2.0) / 2.0, Shape: w.Shape(0)},
	}
	comps = xs[1].PrepareComputation(r, xs...)
	c = w.RefractedColor(comps, 5)
	g.Expect(c).To(Equal(tuple.Black))

	// Test a refrected ray
	mb := material.NewDefaultBuilder()
	w.SetShape(0, shapes.NewSphere().WithMaterial(mb.WithAmbient(1.0).WithPattern(material.NewTestPattern().WithTransform(matrix.NewIdentity())).Build()))
	w.SetShape(1, shapes.NewGlassSphere())
	r, err = shapes.NewRay(tuple.NewPoint(0, 0, 0.1), tuple.NewVector(0, 1, 0))
	g.Expect(err).To(BeNil())
	xs = []shapes.Intersection{
		{T: -0.9899, Shape: w.Shape(0)},
		{T: -0.4899, Shape: w.Shape(1)},
		{T: 0.4899, Shape: w.Shape(1)},
		{T: 0.9899, Shape: w.Shape(0)},
	}
	comps = xs[2].PrepareComputation(r, xs...)
	c = w.RefractedColor(comps, 5)
	g.Expect(c.Equals(tuple.NewColor(0, 0.99888, 0.04721))).To(BeTrue())
}

func TestRefractedShader(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()
	mb := material.NewDefaultBuilder()
	floor := shapes.NewPlane().WithTransform(matrix.NewTranslation(0, -1, 0)).WithMaterial(mb.WithTransparency(0.5).WithRefractiveIndex(1.5).Build())
	ball := shapes.NewSphere().WithTransform(matrix.NewTranslation(0, -3.5, -0.5)).WithMaterial(mb.Reset().WithColor(tuple.Red).WithAmbient(0.5).Build())
	w.AddShapes(floor, ball)

	r, err := shapes.NewRay(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt(2.0)/2.0, math.Sqrt(2.0)/2.0))
	g.Expect(err).To(BeNil())

	xs := []shapes.Intersection{
		{T: math.Sqrt(2.0), Shape: floor},
	}
	comps := xs[0].PrepareComputation(r, xs...)
	c := w.ShadeHit(comps, 5)
	g.Expect(c.Equals(tuple.NewColor(0.93642, 0.68642, 0.68642))).To(BeTrue())
}

func TestSchlickEnabledShader(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()
	mb := material.NewDefaultBuilder()
	floor := shapes.NewPlane().WithTransform(matrix.NewTranslation(0, -1, 0)).WithMaterial(mb.WithTransparency(0.5).WithRefractiveIndex(1.5).WithReflective(0.5).Build())
	ball := shapes.NewSphere().WithTransform(matrix.NewTranslation(0, -3.5, -0.5)).WithMaterial(mb.Reset().WithColor(tuple.Red).WithAmbient(0.5).Build())
	w.AddShapes(floor, ball)

	r, err := shapes.NewRay(tuple.NewPoint(0, 0, -3), tuple.NewVector(0, -math.Sqrt(2.0)/2.0, math.Sqrt(2.0)/2.0))
	g.Expect(err).To(BeNil())

	xs := []shapes.Intersection{
		{T: math.Sqrt(2.0), Shape: floor},
	}
	comps := xs[0].PrepareComputation(r, xs...)
	c := w.ShadeHit(comps, 5)
	g.Expect(c.Equals(tuple.NewColor(0.93391, 0.69643, 0.69243))).To(BeTrue())

}
