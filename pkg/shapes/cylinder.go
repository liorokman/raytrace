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

func (c cylinder) localIntersect(ray Ray, outer Shape) []Intersection {
	A := ray.Direction.X()*ray.Direction.X() + ray.Direction.Z()*ray.Direction.Z()

	if utils.FloatEqual(A, 0.0) {
		return c.intersectCaps(ray, outer)
	}

	B := 2*ray.Origin.X()*ray.Direction.X() +
		2*ray.Origin.Z()*ray.Direction.Z()
	C := ray.Origin.X()*ray.Origin.X() + ray.Origin.Z()*ray.Origin.Z() - 1.0

	disc := B*B - 4*A*C

	if disc < 0 {
		return []Intersection{}
	}

	disc = math.Sqrt(disc)
	t0 := (-B - disc) / (2 * A)
	t1 := (-B + disc) / (2 * A)

	retval := []Intersection{}
	testLimit := func(t float64) bool {
		val := ray.Origin.Y() + t*ray.Direction.Y()
		return c.Min < val && val < c.Max
	}
	if testLimit(t0) {
		retval = append(retval, Intersection{T: t0, Shape: outer})
	}
	if testLimit(t1) {
		retval = append(retval, Intersection{T: t1, Shape: outer})
	}
	retval = append(retval, c.intersectCaps(ray, outer)...)
	return retval
}

func (c cylinder) intersectCaps(ray Ray, outer Shape) []Intersection {
	checkCaps := func(t float64) bool {
		x := ray.Origin.X() + t*ray.Direction.X()
		z := ray.Origin.Z() + t*ray.Direction.Z()
		return (x*x + z*z) <= 1
	}

	if !c.Closed || utils.FloatEqual(ray.Direction.Y(), 0.0) {
		return []Intersection{}
	}

	retval := []Intersection{}
	t := (c.Min - ray.Origin.Y()) / ray.Direction.Y()
	if checkCaps(t) {
		retval = append(retval, Intersection{T: t, Shape: outer})
	}
	t = (c.Max - ray.Origin.Y()) / ray.Direction.Y()
	if checkCaps(t) {
		retval = append(retval, Intersection{T: t, Shape: outer})
	}
	return retval
}
