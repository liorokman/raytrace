package canvas

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestCreateCanvas(t *testing.T) {

	g := NewGomegaWithT(t)

	black := tuple.NewColor(0, 0, 0)

	canv := New(10, 20)
	g.Expect(canv.Width()).To(Equal(uint32(10)))
	g.Expect(canv.Height()).To(Equal(uint32(20)))

	for x := uint32(0); x < canv.Width(); x++ {
		for y := uint32(0); y < canv.Height(); y++ {
			g.Expect(canv.GetPixel(x, y)).To(Equal(black))
		}
	}

	canv.SetPixel(2, 3, tuple.NewColor(1, 0, 0))
	g.Expect(canv.GetPixel(2, 3)).To(Equal(tuple.NewColor(1, 0, 0)))
}

func XTestWritePPM(t *testing.T) {
	g := NewGomegaWithT(t)

	canv := New(10, 10)
	canv.SetPixel(1, 1, tuple.NewColor(1, 0, 0))
	canv.SetPixel(2, 3, tuple.NewColor(0, 0.5, 0))
	err := canv.WritePPM(os.Stdout)
	g.Expect(err).To(BeNil())

	canv = New(10, 2)
	for x := uint32(0); x < canv.Width(); x++ {
		for y := uint32(0); y < canv.Height(); y++ {
			canv.SetPixel(x, y, tuple.NewColor(1, 0.8, 0.6))
		}
	}
	err = canv.WritePPM(os.Stdout)
	g.Expect(err).To(BeNil())
}
