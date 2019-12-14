package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type cylinder struct {
	Min    float64
	Max    float64
	Closed bool
}

func NewCylinder() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), cylinder{
		Min: math.Inf(-1),
		Max: math.Inf(1),
	})
}

func NewConstrainedCylinder(min, max float64, closed bool) Shape {
	return newShape(material.Default(), matrix.NewIdentity(), cylinder{
		Min:    min,
		Max:    max,
		Closed: closed,
	})
}

func (c cylinder) shapeIdPrefix() string {
	return "CY"
}

func (c cylinder) normalAt(point tuple.Tuple) tuple.Tuple {
	dist := point.X()*point.X() + point.Z()*point.Z()
	if dist < 1 && point.Y() >= (c.Max-utils.EPSILON) {
		return tuple.NewVector(0, 1, 0)
	}
	if dist < 1 && point.Y() <= (c.Min+utils.EPSILON) {
		return tuple.NewVector(0, -1, 0)
	}

	return tuple.NewVector(point.X(), 0, point.Z())
}

func (c cylinder) localIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64 {
	A := direction.X()*direction.X() + direction.Z()*direction.Z()

	if utils.FloatEqual(A, 0.0) {
		return c.intersectCaps(direction, origin)
	}

	B := 2*origin.X()*direction.X() +
		2*origin.Z()*direction.Z()
	C := origin.X()*origin.X() + origin.Z()*origin.Z() - 1.0

	disc := B*B - 4*A*C

	if disc < 0 {
		return []float64{}
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

func (c cylinder) intersectCaps(direction tuple.Tuple, origin tuple.Tuple) []float64 {
	checkCaps := func(t float64) bool {
		x := origin.X() + t*direction.X()
		z := origin.Z() + t*direction.Z()
		return (x*x + z*z) <= 1
	}

	if !c.Closed || utils.FloatEqual(direction.Y(), 0.0) {
		return []float64{}
	}

	retval := []float64{}
	t := (c.Min - origin.Y()) / direction.Y()
	if checkCaps(t) {
		retval = append(retval, t)
	}
	t = (c.Max - origin.Y()) / direction.Y()
	if checkCaps(t) {
		retval = append(retval, t)
	}
	return retval
}
