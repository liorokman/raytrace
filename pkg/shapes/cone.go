package shapes

import (
	"fmt"
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

func (c cone) String() string {
	return fmt.Sprintf("Min: %f Max: %f Closed: %t", c.Min, c.Max, c.Closed)
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

func (c cone) localIntersect(ray Ray, outer Shape) []Intersection {

	A := ray.Direction.X()*ray.Direction.X() - ray.Direction.Y()*ray.Direction.Y() + ray.Direction.Z()*ray.Direction.Z()
	B := 2*ray.Origin.X()*ray.Direction.X() - 2*ray.Origin.Y()*ray.Direction.Y() + 2*ray.Origin.Z()*ray.Direction.Z()
	C := ray.Origin.X()*ray.Origin.X() - ray.Origin.Y()*ray.Origin.Y() + ray.Origin.Z()*ray.Origin.Z()

	if utils.FloatEqual(A, 0.0) {
		if !utils.FloatEqual(B, 0.0) {
			retval := c.intersectCaps(ray, outer)
			return append(retval, Intersection{T: -C / (2.0 * B), Shape: outer})
		} else {
			return c.intersectCaps(ray, outer)
		}
	}
	disc := B*B - 4.0*A*C
	if disc < 0.0 {
		return c.intersectCaps(ray, outer)
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

func (c cone) intersectCaps(ray Ray, outer Shape) []Intersection {
	checkCaps := func(t, r float64) bool {
		x := ray.Origin.X() + t*ray.Direction.X()
		z := ray.Origin.Z() + t*ray.Direction.Z()
		return (x*x + z*z) <= math.Abs(r)
	}

	if !c.Closed || utils.FloatEqual(ray.Direction.Y(), 0.0) {
		return []Intersection{}
	}

	retval := []Intersection{}
	t := (c.Min - ray.Origin.Y()) / ray.Direction.Y()
	if checkCaps(t, c.Min) {
		retval = append(retval, Intersection{T: t, Shape: outer})
	}
	t = (c.Max - ray.Origin.Y()) / ray.Direction.Y()
	if checkCaps(t, c.Max) {
		retval = append(retval, Intersection{T: t, Shape: outer})
	}
	return retval
}
