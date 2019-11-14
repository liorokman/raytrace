package utils

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFloatEqual(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(FloatEqual(3, 2)).To(BeFalse())
	g.Expect(FloatEqual(3, 3.000001)).To(BeTrue())
}
