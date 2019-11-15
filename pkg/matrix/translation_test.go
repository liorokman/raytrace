package matrix

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestTranslate(t *testing.T) {

	g := NewGomegaWithT(t)

	tm := NewTranslation(5, -3, 2)
	p := tuple.NewPoint(-3, 4, 5)

	pt := tm.MultiplyTuple(p)

	g.Expect(pt).To(Equal(tuple.NewPoint(2, 1, 7)))

	invTM, err := tm.Inverse()
	g.Expect(err).To(BeNil())

	pt = invTM.MultiplyTuple(p)
	g.Expect(pt).To(Equal(tuple.NewPoint(-8, 7, 3)))

	v := tuple.NewVector(-3, 4, 5)
	g.Expect(tm.MultiplyTuple(v)).To(Equal(v))
}

func TestScale(t *testing.T) {
	g := NewGomegaWithT(t)

	tm := NewScale(2, 3, 4)
	p := tuple.NewPoint(-4, 6, 8)

	pt := tm.MultiplyTuple(p)

	g.Expect(pt).To(Equal(tuple.NewPoint(-8, 18, 32)))

	invTM, err := tm.Inverse()
	g.Expect(err).To(BeNil())

	pt = invTM.MultiplyTuple(p)
	g.Expect(pt).To(Equal(tuple.NewPoint(-2, 2, 2)))

	v := tuple.NewVector(-4, 6, 8)
	g.Expect(tm.MultiplyTuple(v)).To(Equal(tuple.NewVector(-8, 18, 32)))

	// Reflection

	tm = NewScale(-1, 1, 1)
	p = tuple.NewPoint(2, 3, 4)
	pt = tm.MultiplyTuple(p)
	g.Expect(pt).To(Equal(tuple.NewPoint(-2, 3, 4)))

}

func TestRotateX(t *testing.T) {
	g := NewGomegaWithT(t)

	p := tuple.NewPoint(0, 1, 0)
	halfQuarter := NewRotateX(math.Pi / 4)
	fullQuarter := NewRotateX(math.Pi / 2)
	revHalfQuarter, err := halfQuarter.Inverse()
	g.Expect(err).To(BeNil())

	g.Expect(halfQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(0, math.Sqrt(2)/2, math.Sqrt(2)/2))).To(BeTrue())
	g.Expect(fullQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(0, 0, 1))).To(BeTrue())
	g.Expect(revHalfQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))).To(BeTrue())
}

func TestRotateY(t *testing.T) {
	g := NewGomegaWithT(t)

	p := tuple.NewPoint(0, 0, 1)
	halfQuarter := NewRotateY(math.Pi / 4)
	fullQuarter := NewRotateY(math.Pi / 2)

	g.Expect(halfQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(math.Sqrt(2)/2, 0, math.Sqrt(2)/2))).To(BeTrue())
	g.Expect(fullQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(1, 0, 0))).To(BeTrue())
}

func TestRotateZ(t *testing.T) {
	g := NewGomegaWithT(t)

	p := tuple.NewPoint(0, 1, 0)
	halfQuarter := NewRotateZ(math.Pi / 4)
	fullQuarter := NewRotateZ(math.Pi / 2)

	g.Expect(halfQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(-math.Sqrt(2)/2, math.Sqrt(2)/2, 0))).To(BeTrue())
	g.Expect(fullQuarter.MultiplyTuple(p).Equals(tuple.NewPoint(-1, 0, 0))).To(BeTrue())
}

func TestShear(t *testing.T) {
	g := NewGomegaWithT(t)

	tm := NewShear(0, 1, 0, 0, 0, 0)
	p := tuple.NewPoint(2, 3, 4)

	g.Expect(tm.MultiplyTuple(p).Equals(tuple.NewPoint(6, 3, 4))).To(BeTrue())

	tm = NewShear(0, 0, 1, 0, 0, 0)
	g.Expect(tm.MultiplyTuple(p).Equals(tuple.NewPoint(2, 5, 4))).To(BeTrue())

	tm = NewShear(0, 0, 0, 1, 0, 0)
	g.Expect(tm.MultiplyTuple(p).Equals(tuple.NewPoint(2, 7, 4))).To(BeTrue())

	tm = NewShear(0, 0, 0, 0, 1, 0)
	g.Expect(tm.MultiplyTuple(p).Equals(tuple.NewPoint(2, 3, 6))).To(BeTrue())

	tm = NewShear(0, 0, 0, 0, 0, 1)
	g.Expect(tm.MultiplyTuple(p).Equals(tuple.NewPoint(2, 3, 7))).To(BeTrue())
}

func TestChain(t *testing.T) {
	g := NewGomegaWithT(t)

	chain := NewIdentity().RotateX(math.Pi/2).Scale(5, 5, 5).Translate(10, 5, 7)
	p := tuple.NewPoint(1, 0, 1)

	g.Expect(chain.MultiplyTuple(p).Equals(tuple.NewPoint(15, 0, 7)))
}
