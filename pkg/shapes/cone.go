package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type cone struct {
	Min    float64
	Max    float64
	Closed bool
}

func NewCone() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), cone{
		Min: math.Inf(-1),
		Max: math.Inf(1),
	})
}

func (c cone) shapeIdPrefix() string {
	return "CO"
}

func NewConstrainedCone(min, max float64, closed bool) Shape {
	return newShape(material.Default(), matrix.NewIdentity(), cone{
		Min:    min,
		Max:    max,
		Closed: closed,
	})
}

func (c cone) normalAt(point tuple.Tuple) tuple.Tuple {
	dist := point.X()*point.X() + point.Z()*point.Z()
	if dist < 1 && point.Y() >= (c.Max-utils.EPSILON) {
		return tuple.NewVector(0, 1, 0)
	}
	if dist < 1 && point.Y() <= (c.Min+utils.EPSILON) {
		return tuple.NewVector(0, -1, 0)
	}
	y := math.Sqrt(point.X()*point.X() + point.Z()*point.Z())
	if point.Y() > 0 {
		y = -y
	}
	return tuple.NewVector(point.X(), y, point.Z())
}

func (c cone) localIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64 {

	A := direction.X()*direction.X() - direction.Y()*direction.Y() + direction.Z()*direction.Z()
	B := 2*origin.X()*direction.X() - 2*origin.Y()*direction.Y() + 2*origin.Z()*direction.Z()
	C := origin.X()*origin.X() - origin.Y()*origin.Y() + origin.Z()*origin.Z()

	if utils.FloatEqual(A, 0.0) {
		if !utils.FloatEqual(B, 0.0) {
			retval := c.intersectCaps(direction, origin)
			return append(retval, -C/(2.0*B))
		} else {
			return c.intersectCaps(direction, origin)
		}
	}
	disc := B*B - 4.0*A*C
	if disc < 0.0 {
		return c.intersectCaps(direction, origin)
	}
	disc = math.Sqrt(disc)
	t0 := (-B - disc) / (2 * A)
	t1 := (-B + disc) / (2 * A)

	retval := []float64{}
	testLimit := func(t float64) bool {
		val := origin.Y() + t*direction.Y()
		return c.Min < val && val < c.Max
	}
	if testLimit(t0) {
		retval = append(retval, t0)
	}
	if testLimit(t1) {
		retval = append(retval, t1)
	}
	retval = append(retval, c.intersectCaps(direction, origin)...)
	return retval
}

func (c cone) intersectCaps(direction tuple.Tuple, origin tuple.Tuple) []float64 {
	checkCaps := func(t, r float64) bool {
		x := origin.X() + t*direction.X()
		z := origin.Z() + t*direction.Z()
		return (x*x + z*z) <= math.Abs(r)
	}

	if !c.Closed || utils.FloatEqual(direction.Y(), 0.0) {
		return []float64{}
	}

	retval := []float64{}
	t := (c.Min - origin.Y()) / direction.Y()
	if checkCaps(t, c.Min) {
		retval = append(retval, t)
	}
	t = (c.Max - origin.Y()) / direction.Y()
	if checkCaps(t, c.Max) {
		retval = append(retval, t)
	}
	return retval
}
