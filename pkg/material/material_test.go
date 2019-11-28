package material

import (
	"fmt"
	"math"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func TestLighting(t *testing.T) {
	g := NewGomegaWithT(t)
	m := Default()
	pos := tuple.NewPoint(0, 0, 0)
	identity := matrix.NewIdentity()

	eyev := tuple.NewVector(0, 0, -1)
	normalv := tuple.NewVector(0, 0, -1)
	l := fixtures.NewPointLight(tuple.NewPoint(0, 0, -10), tuple.NewColor(1, 1, 1))
	r := m.Lighting(identity, l, pos, eyev, normalv, false)
	g.Expect(r.Equals(tuple.NewColor(1.9, 1.9, 1.9))).To(BeTrue())

	eyev = tuple.NewVector(0, math.Sqrt(2.0)/2.0, -math.Sqrt(2.0)/2.0)
	normalv = tuple.NewVector(0, 0, -1)
	l = fixtures.NewPointLight(tuple.NewPoint(0, 0, -10), tuple.NewColor(1, 1, 1))
	r = m.Lighting(identity, l, pos, eyev, normalv, false)
	g.Expect(r.Equals(tuple.NewColor(1, 1, 1))).To(BeTrue())

	eyev = tuple.NewVector(0, 0, -1)
	normalv = tuple.NewVector(0, 0, -1)
	l = fixtures.NewPointLight(tuple.NewPoint(0, 10, -10), tuple.NewColor(1, 1, 1))
	r = m.Lighting(identity, l, pos, eyev, normalv, false)
	g.Expect(r.Equals(tuple.NewColor(0.7364, 0.7364, 0.7364))).To(BeTrue())

	eyev = tuple.NewVector(0, -math.Sqrt(2.0)/2.0, -math.Sqrt(2.0)/2.0)
	normalv = tuple.NewVector(0, 0, -1)
	l = fixtures.NewPointLight(tuple.NewPoint(0, 10, -10), tuple.NewColor(1, 1, 1))
	r = m.Lighting(identity, l, pos, eyev, normalv, false)
	g.Expect(r.Equals(tuple.NewColor(1.6364, 1.6364, 1.6364))).To(BeTrue())

	eyev = tuple.NewVector(0, 0, -1)
	normalv = tuple.NewVector(0, 0, -1)
	l = fixtures.NewPointLight(tuple.NewPoint(0, 0, -10), tuple.NewColor(1, 1, 1))
	r = m.Lighting(identity, l, pos, eyev, normalv, false)
	fmt.Printf("%#v\n", r)
	g.Expect(r.Equals(tuple.NewColor(1.9, 1.9, 1.9))).To(BeTrue())

	eyev = tuple.NewVector(0, 0, -1)
	normalv = tuple.NewVector(0, 0, -1)
	l = fixtures.NewPointLight(tuple.NewPoint(0, 0, -10), tuple.NewColor(1, 1, 1))
	r = m.Lighting(identity, l, pos, eyev, normalv, true)
	g.Expect(r.Equals(tuple.NewColor(0.1, 0.1, 0.1))).To(BeTrue())
}
