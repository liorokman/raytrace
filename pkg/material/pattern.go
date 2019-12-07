package material

import (
	"math"

	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type Pattern struct {
	transform matrix.Matrix
	pat       pattern
}

func newPattern(pat pattern) Pattern {
	return Pattern{
		transform: matrix.NewIdentity(),
		pat:       pat,
	}
}

func (p Pattern) WithTransform(t matrix.Matrix) Pattern {
	return Pattern{
		transform: t,
		pat:       p.pat,
	}
}

func (p Pattern) PatternAtObject(objTransform matrix.Matrix, point tuple.Tuple) tuple.Color {

	objTransInv, err := objTransform.Inverse()
	if err != nil {
		panic(err)
	}
	objPoint := objTransInv.MultiplyTuple(point)

	patTransInv, err := p.transform.Inverse()
	if err != nil {
		panic(err)
	}
	patPoint := patTransInv.MultiplyTuple(objPoint)
	return p.pat.ColorAt(patPoint)
}

func (p Pattern) ColorAt(point tuple.Tuple) tuple.Color {
	return p.pat.ColorAt(point)
}

type pattern interface {
	ColorAt(tuple.Tuple) tuple.Color
}

type solidPattern tuple.Color

func NewSolidPattern(c tuple.Color) Pattern {
	return newPattern(solidPattern(c))
}

func (p solidPattern) ColorAt(point tuple.Tuple) tuple.Color {
	return tuple.Color(p)
}

type stripePattern struct {
	colors []tuple.Color
}

func NewStripePattern(c1, c2 tuple.Color) Pattern {
	return newPattern(stripePattern{
		colors: []tuple.Color{c1, c2},
	})
}

func (p stripePattern) ColorAt(point tuple.Tuple) tuple.Color {
	return p.colors[int(math.Abs(math.Floor(point.X())))%2]
}

type gradientPattern struct {
	base     tuple.Color
	distance tuple.Color
}

func NewGradientPattern(c1, c2 tuple.Color) Pattern {
	return newPattern(gradientPattern{
		base:     c1,
		distance: c2.Subtract(c1),
	})
}

func (p gradientPattern) ColorAt(point tuple.Tuple) tuple.Color {
	a := math.Abs(point.X())
	fraction := a - math.Floor(a)
	//	fmt.Printf("%#v %#v %#v\n ", point, fraction, p.base.Add(p.distance.Mult(fraction)))
	return p.base.Add(p.distance.Mult(fraction))
}

type ringPattern struct {
	colors []tuple.Color
}

func NewRingPattern(c1, c2 tuple.Color) Pattern {
	return newPattern(ringPattern{
		colors: []tuple.Color{c1, c2},
	})
}

func (p ringPattern) ColorAt(point tuple.Tuple) tuple.Color {
	c := math.Floor(math.Sqrt((point.X() * point.X()) + (point.Z() * point.Z())))
	return p.colors[int(c)%2]
}

type checkerPattern struct {
	colors []tuple.Color
}

func NewCheckerPattern(c1, c2 tuple.Color) Pattern {
	return newPattern(checkerPattern{
		colors: []tuple.Color{c1, c2},
	})
}

func (p checkerPattern) ColorAt(point tuple.Tuple) tuple.Color {
	c := math.Abs(math.Floor(point.X()+utils.EPSILON) + math.Floor(point.Y()+utils.EPSILON) + math.Floor(point.Z()+utils.EPSILON))
	return p.colors[int(c)%2]
}

type testPattern struct{}

func NewTestPattern() Pattern {
	return newPattern(testPattern{})
}

func (p testPattern) ColorAt(point tuple.Tuple) tuple.Color {
	return tuple.NewColor(point.X(), point.Y(), point.Z())
}
