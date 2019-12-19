package material

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestStripePattern(t *testing.T) {
	g := NewGomegaWithT(t)
	sp := NewStripePattern(tuple.White, tuple.Black)
	// Stable in Y()
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 0, 0))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 1, 0))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 2, 0))).To(Equal(tuple.White))

	// Stable in Z()
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 0, 1))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 0, 2))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 0, 3))).To(Equal(tuple.White))

	// Alternates in X()
	g.Expect(sp.ColorAt(tuple.NewPoint(0, 0, 1))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(0.9, 0, 2))).To(Equal(tuple.White))
	g.Expect(sp.ColorAt(tuple.NewPoint(1, 0, 3))).To(Equal(tuple.Black))
	g.Expect(sp.ColorAt(tuple.NewPoint(-0.1, 0, 3))).To(Equal(tuple.Black))
	g.Expect(sp.ColorAt(tuple.NewPoint(-1, 0, 3))).To(Equal(tuple.Black))
	g.Expect(sp.ColorAt(tuple.NewPoint(-1.1, 0, 3))).To(Equal(tuple.White))
}

func TestRingPattern(t *testing.T) {
	g := NewGomegaWithT(t)

	p := NewRingPattern(tuple.White, tuple.Black)
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(1, 0, 0))).To(Equal(tuple.Black))
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 1))).To(Equal(tuple.Black))
	// 0.708 = just slightly more than √2/2
	g.Expect(p.ColorAt(tuple.NewPoint(0.708, 0, 0.708))).To(Equal(tuple.Black))
}

func TestCheckerPattern(t *testing.T) {
	g := NewGomegaWithT(t)

	p := NewCheckerPattern(tuple.White, tuple.Black)
	// Repeat in X
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(0.99, 0, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(1.01, 0, 0))).To(Equal(tuple.Black))
	// Repeat in Y
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0.99, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(0, 1.01, 0))).To(Equal(tuple.Black))
	// Repeat in Z
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 0))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 0.99))).To(Equal(tuple.White))
	g.Expect(p.ColorAt(tuple.NewPoint(0, 0, 1.01))).To(Equal(tuple.Black))
}

func TestPatternAt(t *testing.T) {
	g := NewGomegaWithT(t)

	transform := testShape{matrix.NewScale(2, 2, 2)}

	p := NewStripePattern(tuple.White, tuple.Black)
	g.Expect(p.PatternAtObject(transform, tuple.NewPoint(1.5, 0, 0))).To(Equal(tuple.White))

	p = p.WithTransform(matrix.NewScale(2, 2, 2))
	g.Expect(p.PatternAtObject(testShape{matrix.NewIdentity()}, tuple.NewPoint(1.5, 0, 0))).To(Equal(tuple.White))

	g.Expect(p.PatternAtObject(transform, tuple.NewPoint(1.5, 0, 0))).To(Equal(tuple.White))
}
