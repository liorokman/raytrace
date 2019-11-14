package tuple

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/utils"
)

func TestColor(t *testing.T) {
	g := NewGomegaWithT(t)

	c := NewColor(-0.5, 0.4, 1.7)
	g.Expect(c.Red()).To(Equal(-0.5))
	g.Expect(c.Green()).To(Equal(0.4))
	g.Expect(c.Blue()).To(Equal(1.7))
}

func TestColorOps(t *testing.T) {

	g := NewGomegaWithT(t)

	g.Expect(NewColor(0.9, 0.6, 0.75).Add(NewColor(0.7, 0.1, 0.25))).To(Equal(NewColor(1.6, 0.7, 1.0)))
	sub := NewColor(0.9, 0.6, 0.75).Subtract(NewColor(0.7, 0.1, 0.25))
	res := NewColor(0.2, 0.5, 0.5)
	g.Expect(utils.FloatEqual(sub.Red(), res.Red())).To(BeTrue())
	g.Expect(utils.FloatEqual(sub.Green(), res.Green())).To(BeTrue())
	g.Expect(utils.FloatEqual(sub.Blue(), res.Blue())).To(BeTrue())

	g.Expect(NewColor(0.2, 0.3, 0.4).Mult(2.0)).To(Equal(NewColor(0.4, 0.6, 0.8)))

	sub = NewColor(1, 0.2, 0.4).MultColor(NewColor(0.9, 1, 0.1))
	res = NewColor(0.9, 0.2, 0.04)
	g.Expect(utils.FloatEqual(sub.Red(), res.Red())).To(BeTrue())
	g.Expect(utils.FloatEqual(sub.Green(), res.Green())).To(BeTrue())
	g.Expect(utils.FloatEqual(sub.Blue(), res.Blue())).To(BeTrue())

}
