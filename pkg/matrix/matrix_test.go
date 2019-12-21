package matrix

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestNewMatrix(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			1, 2, 3, 4,
			5.5, 6.5, 7.5, 8.5,
			9, 10, 11, 12,
			13.5, 14.5, 15.5, 16.5},
	}

	g.Expect(m.data[0]).To(Equal(1.0))
	g.Expect(m.data[3]).To(Equal(4.0))
	g.Expect(m.data[1*m.rows]).To(Equal(5.5))
	g.Expect(m.data[1*m.rows+2]).To(Equal(7.5))
	g.Expect(m.data[2*m.rows+2]).To(Equal(11.0))
	g.Expect(m.data[3*m.rows+0]).To(Equal(13.5))
	g.Expect(m.data[3*m.rows+2]).To(Equal(15.5))

	m = Matrix{
		rows: 2,
		cols: 2,
		data: []float64{
			-3, 5,
			1, -2},
	}
	g.Expect(m.data[0]).To(Equal(-3.0))
	g.Expect(m.data[1]).To(Equal(5.0))
	g.Expect(m.data[1*m.rows]).To(Equal(1.0))
	g.Expect(m.data[1*m.rows+1]).To(Equal(-2.0))
}

func TestEquality(t *testing.T) {

	g := NewGomegaWithT(t)

	m1 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			9, 8, 7, 6,
			5, 4, 3, 2},
	}
	m2 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			1, 2, 3, 4,
			9, 8, 7, 6,
			5, 6, 7, 8,
			5, 4, 3, 2},
	}

	g.Expect(m1.Equals(m1)).To(BeTrue())
	g.Expect(m1.Equals(m2)).To(BeFalse())

	m3 := Matrix{
		rows: 2,
		cols: 2,
		data: []float64{
			1, 2,
			3, 4},
	}

	g.Expect(m1.Equals(m3)).To(BeFalse())
}

func TestMultiplication(t *testing.T) {

	g := NewGomegaWithT(t)
	m1 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			1, 2, 3, 4,
			5, 6, 7, 8,
			9, 8, 7, 6,
			5, 4, 3, 2},
	}

	m2 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			-2, 1, 2, 3,
			3, 2, 1, -1,
			4, 3, 6, 5,
			1, 2, 7, 8},
	}

	m3 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			20, 22, 50, 48,
			44, 54, 114, 108,
			40, 58, 110, 102,
			16, 26, 46, 42},
	}

	g.Expect(m1.Multiply(m2).Equals(m3)).To(BeTrue())
}

func TestMultiplyTuple(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			1, 2, 3, 4,
			2, 4, 4, 2,
			8, 6, 4, 1,
			0, 0, 0, 1},
	}

	t1 := tuple.Tuple{1, 2, 3, 1}
	r1 := tuple.Tuple{18, 24, 33, 1}
	g.Expect(m1.MultiplyTuple(t1)).To(Equal(r1))
}

func TestTranspose(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			0, 9, 3, 0,
			9, 8, 0, 8,
			1, 8, 5, 3,
			0, 0, 5, 8},
	}

	m2 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			0, 9, 1, 0,
			9, 8, 8, 0,
			3, 0, 5, 5,
			0, 8, 3, 8},
	}

	g.Expect(m1.Transpose().Equals(m2)).To(BeTrue())
}

func TestSubmatrix(t *testing.T) {
	g := NewGomegaWithT(t)

	m1 := Matrix{
		rows: 3,
		cols: 3,
		data: []float64{
			1, 5, 0,
			-3, 2, 7,
			0, 6, 3},
	}

	g.Expect(m1.Submatrix(0, 2).Equals(Matrix{rows: 2, cols: 2, data: []float64{-3, 2, 0, 6}})).To(BeTrue())
	g.Expect(m1.Submatrix(1, 1).Equals(Matrix{rows: 2, cols: 2, data: []float64{1, 0, 0, 3}})).To(BeTrue())
	g.Expect(m1.Submatrix(0, 0).Equals(Matrix{rows: 2, cols: 2, data: []float64{2, 7, 6, 3}})).To(BeTrue())

	m2 := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			-6, 1, 1, 6,
			-8, 5, 8, 6,
			-1, 0, 8, 2,
			-7, 1, -1, 1},
	}
	g.Expect(m2.Submatrix(2, 1).Equals(Matrix{rows: 3, cols: 3, data: []float64{-6, 1, 6, -8, 8, 6, -7, -1, 1}})).To(BeTrue())
}

func TestMinor(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Matrix{
		rows: 3,
		cols: 3,
		data: []float64{
			3, 5, 0,
			2, -1, -7,
			6, -1, 5},
	}
	g.Expect(m.Minor(1, 0)).To(Equal(25.0))
	g.Expect(m.Minor(0, 0)).To(Equal(-12.0))
	g.Expect(m.Minor(1, 1)).To(Equal(15.0))
	g.Expect(m.Minor(2, 2)).To(Equal(-13.0))
}

func TestCofactor(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Matrix{
		rows: 3,
		cols: 3,
		data: []float64{
			3, 5, 0,
			2, -1, -7,
			6, -1, 5},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(-12.0))
	g.Expect(m.Cofactor(1, 0)).To(Equal(-25.0))
}

func TestDeterminant(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		rows: 3,
		cols: 3,
		data: []float64{
			1, 2, 6,
			-5, 8, -4,
			2, 6, 4},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(56.0))
	g.Expect(m.Cofactor(0, 1)).To(Equal(12.0))
	g.Expect(m.Cofactor(0, 2)).To(Equal(-46.0))
	g.Expect(m.Determinant()).To(Equal(-196.0))

	m = Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			-2, -8, 3, 5,
			-3, 1, 7, 3,
			1, 2, -9, 6,
			-6, 7, 7, -9},
	}
	g.Expect(m.Cofactor(0, 0)).To(Equal(690.0))
	g.Expect(m.Cofactor(0, 1)).To(Equal(447.0))
	g.Expect(m.Cofactor(0, 2)).To(Equal(210.0))
	g.Expect(m.Cofactor(0, 3)).To(Equal(51.0))
	g.Expect(m.Determinant()).To(Equal(-4071.0))
}

func BenchmarkInverse(b *testing.B) {

	for bb := 0; bb < b.N; bb++ {

		m := Matrix{
			rows: 4,
			cols: 4,
			data: []float64{
				8, -5, 9, 2,
				7, 5, 6, 1,
				-6, 0, 9, 6,
				-3, 0, -9, -4},
		}

		m.Inverse()
	}

}

func TestInverse(t *testing.T) {
	g := NewGomegaWithT(t)

	m := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			8, -5, 9, 2,
			7, 5, 6, 1,
			-6, 0, 9, 6,
			-3, 0, -9, -4},
	}

	invM := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			-0.15385, -0.15385, -0.28205, -0.53846,
			-0.07692, 0.12308, 0.02564, 0.03077,
			0.35897, 0.35897, 0.43590, 0.92308,
			-0.69231, -0.69231, -0.76923, -1.92308},
	}
	calcInv, err := m.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(calcInv.Equals(invM)).To(BeTrue())

	m = Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			9, 3, 0, 9,
			-5, -2, -6, -3,
			-4, 9, 6, 4,
			-7, 6, 6, 2},
	}
	invM = Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			-0.04074, -0.07778, 0.14444, -0.22222,
			-0.07778, 0.03333, 0.36667, -0.33333,
			-0.02901, -0.14630, -0.10926, 0.12963,
			0.17778, 0.06667, -0.26667, 0.33333},
	}
	calcInv, err = m.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(calcInv.Equals(invM)).To(BeTrue())

	a := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			3, -9, 7, 3,
			3, -8, 2, -9,
			-4, 4, 4, 1,
			-6, 5, -1, 1},
	}
	b := Matrix{
		rows: 4,
		cols: 4,
		data: []float64{
			8, 2, 2, 2,
			3, -1, 7, 0,
			7, 0, 5, 4,
			6, -2, 0, 5},
	}

	c := a.Multiply(b)
	invB, err := b.Inverse()
	g.Expect(err).To(BeNil())
	g.Expect(c.Multiply(invB).Equals(a)).To(BeTrue())
}
