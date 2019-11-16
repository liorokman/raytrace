package ray

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/shapes"
)

func TestHit(t *testing.T) {
	g := NewGomegaWithT(t)

	s := shapes.NewSphere()

	i, ok := Hit(Intersection{-1, s}, Intersection{1, s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{1, s}))

	i, ok = Hit(Intersection{1, s}, Intersection{2, s})
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{1, s}))

	i, ok = Hit(Intersection{-1, s}, Intersection{-2, s})
	g.Expect(ok).To(BeFalse())

	i, ok = Hit(
		Intersection{5, s},
		Intersection{7, s},
		Intersection{-3, s},
		Intersection{2, s},
	)
	g.Expect(ok).To(BeTrue())
	g.Expect(i).To(Equal(Intersection{2, s}))
}
