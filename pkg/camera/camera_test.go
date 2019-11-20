package camera

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/world"
)

func defaultWorld() *world.World {
	w := world.New()

	mBuilder := material.NewDefaultBuilder()
	mBuilder.
		WithColor(tuple.NewColor(0.8, 1, 0.6)).
		WithDiffuse(0.7).
		WithSpecular(0.2)
	w.AddShapes(
		shapes.NewSphere().WithMaterial(mBuilder.Build()),
		shapes.NewSphere().WithTransform(matrix.NewScale(0.5, 0.5, 0.5)),
	)
	return &w
}

func TestViewTransformation(t *testing.T) {
	g := NewGomegaWithT(t)

	from := tuple.NewPoint(0, 0, 0)
	to := tuple.NewPoint(0, 0, -1)
	up := tuple.NewVector(0, 1, 0)
	v := ViewTransformation(from, to, up)

	g.Expect(v.Equals(matrix.NewIdentity())).To(BeTrue())

	from = tuple.NewPoint(0, 0, 0)
	to = tuple.NewPoint(0, 0, 1)
	up = tuple.NewVector(0, 1, 0)
	v = ViewTransformation(from, to, up)

	g.Expect(v.Equals(matrix.NewScale(-1, 1, -1))).To(BeTrue())

	from = tuple.NewPoint(0, 0, 8)
	to = tuple.NewPoint(0, 0, 0)
	up = tuple.NewVector(0, 1, 0)
	v = ViewTransformation(from, to, up)

	g.Expect(v.Equals(matrix.NewTranslation(0, 0, -8))).To(BeTrue())

	from = tuple.NewPoint(1, 3, 2)
	to = tuple.NewPoint(4, -2, 8)
	up = tuple.NewVector(1, 1, 0)
	v = ViewTransformation(from, to, up)

	g.Expect(v.Equals(matrix.Matrix{
		{-0.50709, 0.50709, 0.67612, -2.36643},
		{0.76772, 0.60609, 0.12122, -2.82843},
		{-0.35857, 0.59761, -0.71714, 0},
		{0, 0, 0, 1},
	})).To(BeTrue())
}

func TestNewCamera(t *testing.T) {
	g := NewGomegaWithT(t)
	c := NewCamera(160, 120, math.Pi/2.0)
	g.Expect(c.HSize()).To(Equal(uint32(160)))
	g.Expect(c.VSize()).To(Equal(uint32(120)))
	g.Expect(c.FieldOfView()).To(Equal(math.Pi / 2.0))
	g.Expect(c.Transform().Equals(matrix.NewIdentity())).To(BeTrue())
}

func TestPixelSize(t *testing.T) {
	g := NewGomegaWithT(t)
	c := NewCamera(200, 125, math.Pi/2.0)
	g.Expect(c.PixelSize()).To(Equal(0.01))

	c = NewCamera(125, 200, math.Pi/2.0)
	g.Expect(c.PixelSize()).To(Equal(0.01))

}

func TestRayForPixel(t *testing.T) {
	g := NewGomegaWithT(t)
	c := NewCamera(201, 101, math.Pi/2)

	// Through the center of the canvas
	r := c.RayForPixel(uint32(100), uint32(50))
	g.Expect(r.Origin.Equals(tuple.NewPoint(0, 0, 0))).To(BeTrue())
	g.Expect(r.Direction.Equals(tuple.NewVector(0, 0, -1))).To(BeTrue())

	// Through the corner of the canvas
	r = c.RayForPixel(uint32(0), uint32(0))
	g.Expect(r.Origin.Equals(tuple.NewPoint(0, 0, 0))).To(BeTrue())
	g.Expect(r.Direction.Equals(tuple.NewVector(0.66519, 0.33259, -0.66851))).To(BeTrue())

	// When the camera is transformed
	c = c.WithTransform(matrix.NewRotateY(math.Pi/4.0).Translate(0, -2, 5))
	r = c.RayForPixel(uint32(100), uint32(50))
	g.Expect(r.Origin.Equals(tuple.NewPoint(0, 2, -5))).To(BeTrue())
	g.Expect(r.Direction.Equals(tuple.NewVector(math.Sqrt(2.0)/2.0, 0.0, -math.Sqrt(2.0)/2.0))).To(BeTrue())
}
func TestRender(t *testing.T) {
	g := NewGomegaWithT(t)
	w := defaultWorld()
	c := NewCamera(11, 11, math.Pi/2)
	from := tuple.NewPoint(0, 0, -5)
	to := tuple.NewPoint(0, 0, 0)
	up := tuple.NewVector(0, 1, 0)
	c = c.WithTransform(ViewTransformation(from, to, up))

	image := c.Render(w)
	pixel, err := image.GetPixel(5, 5)
	g.Expect(err).To(BeNil())
	g.Expect(pixel.Equals(tuple.NewColor(0.38066, 0.47583, 0.2855))).To(BeTrue())

}
