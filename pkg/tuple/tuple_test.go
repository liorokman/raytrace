package tuple

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/utils"
)

func TestIsPoint(t *testing.T) {
	g := NewGomegaWithT(t)

	a := Tuple{4.3, -4.2, 3.1, 1.0}

	g.Expect(a[XPos]).To(Equal(4.3))
	g.Expect(a[YPos]).To(Equal(-4.2))
	g.Expect(a[ZPos]).To(Equal(3.1))
	g.Expect(a[WPos]).To(Equal(1.0))
	g.Expect(a.IsPoint()).To(BeTrue())
	g.Expect(a.IsVector()).To(BeFalse())
}

func TestIsVector(t *testing.T) {
	g := NewGomegaWithT(t)

	a := Tuple{4.3, -4.2, 3.1, 0.0}

	g.Expect(a.X()).To(Equal(4.3))
	g.Expect(a.Y()).To(Equal(-4.2))
	g.Expect(a.Z()).To(Equal(3.1))
	g.Expect(a.W()).To(Equal(0.0))
	g.Expect(a.IsPoint()).To(BeFalse())
	g.Expect(a.IsVector()).To(BeTrue())
}

func TestNewPoint(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewPoint(4, -4, 3)
	g.Expect(a[XPos]).To(Equal(4.0))
	g.Expect(a[YPos]).To(Equal(-4.0))
	g.Expect(a[ZPos]).To(Equal(3.0))
	g.Expect(a[WPos]).To(Equal(1.0))
	g.Expect(a.IsPoint()).To(BeTrue())
}

func TestNewVector(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewVector(4, -4, 3)
	g.Expect(a[XPos]).To(Equal(4.0))
	g.Expect(a[YPos]).To(Equal(-4.0))
	g.Expect(a[ZPos]).To(Equal(3.0))
	g.Expect(a[WPos]).To(Equal(0.0))
	g.Expect(a.IsVector()).To(BeTrue())
}

func TestAddTuple(t *testing.T) {

	g := NewGomegaWithT(t)

	l := NewPoint(3, -2, 5)
	r := NewVector(-2, 3, 1)

	sum := l.Add(r)

	g.Expect(sum.IsPoint()).To(BeTrue())
	g.Expect(sum.IsVector()).To(BeFalse())
	g.Expect(sum.X()).To(Equal(1.0))
	g.Expect(sum.Y()).To(Equal(1.0))
	g.Expect(sum.Z()).To(Equal(6.0))

}

func TestSubtractTuple(t *testing.T) {

	g := NewGomegaWithT(t)

	l := NewPoint(3, 2, 1)
	r := NewPoint(5, 6, 7)

	sum := l.Subtract(r)

	g.Expect(sum.IsPoint()).To(BeFalse())
	g.Expect(sum.IsVector()).To(BeTrue())
	g.Expect(sum.X()).To(Equal(-2.0))
	g.Expect(sum.Y()).To(Equal(-4.0))
	g.Expect(sum.Z()).To(Equal(-6.0))

	v := NewVector(5, 6, 7)

	sum = l.Subtract(v)
	g.Expect(sum.IsPoint()).To(BeTrue())
	g.Expect(sum.IsVector()).To(BeFalse())
	g.Expect(sum.X()).To(Equal(-2.0))
	g.Expect(sum.Y()).To(Equal(-4.0))
	g.Expect(sum.Z()).To(Equal(-6.0))
}

func TestNegateTuple(t *testing.T) {
	g := NewGomegaWithT(t)

	tup := Tuple{1, -2, 3, -4}.Negate()

	g.Expect(tup.IsPoint()).To(BeFalse())
	g.Expect(tup.IsVector()).To(BeFalse())
	g.Expect(tup.X()).To(Equal(-1.0))
	g.Expect(tup.Y()).To(Equal(2.0))
	g.Expect(tup.Z()).To(Equal(-3.0))
	g.Expect(tup.W()).To(Equal(4.0))
}

func TestMultTuple(t *testing.T) {
	g := NewGomegaWithT(t)

	tup := Tuple{1, -2, 3, -4}

	a := tup.Mult(3.5)
	g.Expect(a.IsPoint()).To(BeFalse())
	g.Expect(a.IsVector()).To(BeFalse())
	g.Expect(a.X()).To(Equal(3.5))
	g.Expect(a.Y()).To(Equal(-7.0))
	g.Expect(a.Z()).To(Equal(10.5))
	g.Expect(a.W()).To(Equal(-14.0))

	a = tup.Mult(0.5)
	g.Expect(a.IsPoint()).To(BeFalse())
	g.Expect(a.IsVector()).To(BeFalse())
	g.Expect(a.X()).To(Equal(0.5))
	g.Expect(a.Y()).To(Equal(-1.0))
	g.Expect(a.Z()).To(Equal(1.5))
	g.Expect(a.W()).To(Equal(-2.0))
}

func TestDivTuple(t *testing.T) {
	g := NewGomegaWithT(t)

	tup := Tuple{1, -2, 3, -4}

	a := tup.Div(2)
	g.Expect(a.IsPoint()).To(BeFalse())
	g.Expect(a.IsVector()).To(BeFalse())
	g.Expect(a.X()).To(Equal(0.5))
	g.Expect(a.Y()).To(Equal(-1.0))
	g.Expect(a.Z()).To(Equal(1.5))
	g.Expect(a.W()).To(Equal(-2.0))
}

func TestMagnitude(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(NewVector(0, 1, 0).Magnitude()).To(Equal(1.0))
	g.Expect(NewVector(0, 0, 1).Magnitude()).To(Equal(1.0))
	g.Expect(NewVector(1, 2, 3).Magnitude()).To(Equal(math.Sqrt(14)))
	g.Expect(NewVector(-1, -2, -3).Magnitude()).To(Equal(math.Sqrt(14)))
}

func TestNormalize(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(NewVector(4, 0, 0).Normalize()).To(Equal(NewVector(1, 0, 0)))
	tmp := 1 / math.Sqrt(14)
	n := NewVector(1, 2, 3).Normalize()
	g.Expect(utils.FloatEqual(n.X(), tmp)).To(BeTrue())
	g.Expect(n.Magnitude()).To(Equal(1.0))
}

func TestDot(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(NewVector(1, 2, 3).Dot(NewVector(2, 3, 4))).To(Equal(20.0))
}

func TestCross(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewVector(1, 2, 3)
	b := NewVector(2, 3, 4)
	g.Expect(a.Cross(b)).To(Equal(NewVector(-1, 2, -1)))
	g.Expect(b.Cross(a)).To(Equal(NewVector(1, -2, 1)))

	g.Expect(func() { a.Cross(NewPoint(1, 2, 3)) }).To(Panic())
	g.Expect(func() { NewPoint(1, 2, 3).Cross(a) }).To(Panic())
}

func TestReflect(t *testing.T) {
	g := NewGomegaWithT(t)

	v := math.Sqrt(2.0) / 2.0

	p := NewVector(0, -1, 0)
	n := NewVector(v, v, 0)

	r := p.Reflect(n)
	g.Expect(r.Equals(NewVector(1, 0, 0))).To(BeTrue())

	p = NewVector(1, -1, 0)
	n = NewVector(0, 1, 0)
	r = p.Reflect(n)
	g.Expect(r.Equals(NewVector(1, 1, 0))).To(BeTrue())

}
