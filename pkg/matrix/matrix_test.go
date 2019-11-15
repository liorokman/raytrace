package matrix

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestNewMatrix(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		[]float64{1, 2, 3, 4},
		[]float64{5.5, 6.5, 7.5, 8.5},
		[]float64{9, 10, 11, 12},
		[]float64{13.5, 14.5, 15.5, 16.5},
	}

	g.Expect(m[0][0]).To(Equal(1.0))
	g.Expect(m[0][3]).To(Equal(4.0))
	g.Expect(m[1][0]).To(Equal(5.5))
	g.Expect(m[1][2]).To(Equal(7.5))
	g.Expect(m[2][2]).To(Equal(11.0))
	g.Expect(m[3][0]).To(Equal(13.5))
	g.Expect(m[3][2]).To(Equal(15.5))

	m = Matrix{
		[]float64{-3, 5},
		[]float64{1, -2},
	}
	g.Expect(m[0][0]).To(Equal(-3.0))
	g.Expect(m[0][1]).To(Equal(5.0))
	g.Expect(m[1][0]).To(Equal(1.0))
	g.Expect(m[1][1]).To(Equal(-2.0))
}

func TestEquality(t *testing.T) {

	g := NewGomegaWithT(t)

	m1 := Matrix{
		[]float64{1, 2, 3, 4},
		[]float64{5, 6, 7, 8},
		[]float64{9, 8, 7, 6},
		[]float64{5, 4, 3, 2},
	}
	m2 := Matrix{
		[]float64{1, 2, 3, 4},
		[]float64{9, 8, 7, 6},
		[]float64{5, 6, 7, 8},
		[]float64{5, 4, 3, 2},
	}

	g.Expect(m1.Equals(m1)).To(BeTrue())
	g.Expect(m1.Equals(m2)).To(BeFalse())

	m3 := Matrix{
		[]float64{1, 2},
		[]float64{3, 4},
	}

	g.Expect(m1.Equals(m3)).To(BeFalse())
}

func TestMultiplication(t *testing.T) {

	g := NewGomegaWithT(t)
	m1 := Matrix{
		[]float64{1, 2, 3, 4},
		[]float64{5, 6, 7, 8},
		[]float64{9, 8, 7, 6},
		[]float64{5, 4, 3, 2},
	}

	m2 := Matrix{
		[]float64{-2, 1, 2, 3},
		[]float64{3, 2, 1, -1},
		[]float64{4, 3, 6, 5},
		[]float64{1, 2, 7, 8},
	}

	m3 := Matrix{
		[]float64{20, 22, 50, 48},
		[]float64{44, 54, 114, 108},
		[]float64{40, 58, 110, 102},
		[]float64{16, 26, 46, 42},
	}

	g.Expect(m1.Multiply(m2).Equals(m3)).To(BeTrue())
}

func TestMultiplyTuple(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		[]float64{1, 2, 3, 4},
		[]float64{2, 4, 4, 2},
		[]float64{8, 6, 4, 1},
		[]float64{0, 0, 0, 1},
	}

	t1 := tuple.Tuple{1, 2, 3, 1}
	r1 := tuple.Tuple{18, 24, 33, 1}
	g.Expect(m1.MultiplyTuple(t1)).To(Equal(r1))
}

func TestTranspose(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		[]float64{0, 9, 3, 0},
		[]float64{9, 8, 0, 8},
		[]float64{1, 8, 5, 3},
		[]float64{0, 0, 5, 8},
	}

	m2 := Matrix{
		[]float64{0, 9, 1, 0},
		[]float64{9, 8, 8, 0},
		[]float64{3, 0, 5, 5},
		[]float64{0, 8, 3, 8},
	}

	g.Expect(m1.Transpose().Equals(m2)).To(BeTrue())
}

func TestSubmatrix(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		[]float64{1, 5, 0},
		[]float64{-3, 2, 7},
		[]float64{0, 6, 3},
	}

	g.Expect(m1.Submatrix(0, 2).Equals(Matrix{{-3, 2}, {0, 6}})).To(BeTrue())
	g.Expect(m1.Submatrix(1, 1).Equals(Matrix{{1, 0}, {0, 3}})).To(BeTrue())
	g.Expect(m1.Submatrix(0, 0).Equals(Matrix{{2, 7}, {6, 3}})).To(BeTrue())

	m2 := Matrix{
		[]float64{-6, 1, 1, 6},
		[]float64{-8, 5, 8, 6},
		[]float64{-1, 0, 8, 2},
		[]float64{-7, 1, -1, 1},
	}
	g.Expect(m2.Submatrix(2, 1).Equals(Matrix{{-6, 1, 6}, {-8, 8, 6}, {-7, -1, 1}})).To(BeTrue())
}

func TestMinor(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Matrix{
		{3, 5, 0},
		{2, -1, -7},
		{6, -1, 5},
	}
	g.Expect(m.Minor(1, 0)).To(Equal(25.0))
	g.Expect(m.Minor(0, 0)).To(Equal(-12.0))
	g.Expect(m.Minor(1, 1)).To(Equal(15.0))
	g.Expect(m.Minor(2, 2)).To(Equal(-13.0))
}

func TestCofactor(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Matrix{
		{3, 5, 0},
		{2, -1, -7},
		{6, -1, 5},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(-12.0))
	g.Expect(m.Cofactor(1, 0)).To(Equal(-25.0))
}

func TestDeterminant(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		{1, 2, 6},
		{-5, 8, -4},
		{2, 6, 4},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(56.0))
	g.Expect(m.Cofactor(0, 1)).To(Equal(12.0))
	g.Expect(m.Cofactor(0, 2)).To(Equal(-46.0))
	g.Expect(m.Determinant()).To(Equal(-196.0))

	m = Matrix{
		{-2, -8, 3, 5},
		{-3, 1, 7, 3},
		{1, 2, -9, 6},
		{-6, 7, 7, -9},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(690.0))
	g.Expect(m.Cofactor(0, 1)).To(Equal(447.0))
	g.Expect(m.Cofactor(0, 2)).To(Equal(210.0))
	g.Expect(m.Cofactor(0, 3)).To(Equal(51.0))
	g.Expect(m.Determinant()).To(Equal(-4071.0))
}

func TestInverse(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		{8, -5, 9, 2},
		{7, 5, 6, 1},
		{-6, 0, 9, 6},
		{-3, 0, -9, -4},
	}

	invM := Matrix{
		{-0.15385, -0.15385, -0.28205, -0.53846},
		{-0.07692, 0.12308, 0.02564, 0.03077},
		{0.35897, 0.35897, 0.43590, 0.92308},
		{-0.69231, -0.69231, -0.76923, -1.92308},
	}
	calcInv, err := m.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(calcInv.Equals(invM)).To(BeTrue())

	m = Matrix{
		{9, 3, 0, 9},
		{-5, -2, -6, -3},
		{-4, 9, 6, 4},
		{-7, 6, 6, 2},
	}
	invM = Matrix{
		{-0.04074, -0.07778, 0.14444, -0.22222},
		{-0.07778, 0.03333, 0.36667, -0.33333},
		{-0.02901, -0.14630, -0.10926, 0.12963},
		{0.17778, 0.06667, -0.26667, 0.33333},
	}
	calcInv, err = m.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(calcInv.Equals(invM)).To(BeTrue())

	a := Matrix{
		{3, -9, 7, 3},
		{3, -8, 2, -9},
		{-4, 4, 4, 1},
		{-6, 5, -1, 1},
	}
	b := Matrix{
		{8, 2, 2, 2},
		{3, -1, 7, 0},
		{7, 0, 5, 4},
		{6, -2, 0, 5},
	}

	c := a.Multiply(b)
	invB, err := b.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(c.Multiply(invB).Equals(a)).To(BeTrue())
}
