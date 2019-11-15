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
