package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type cube struct {
}

func NewCube() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), cube{})
}

func (c cube) shapeIdPrefix() string {
	return "C"
}

func (c cube) normalAt(point tuple.Tuple, hit Intersection) tuple.Tuple {
	maxc := math.Max(math.Max(math.Abs(point.X()), math.Abs(point.Y())), math.Abs(point.Z()))
	if maxc == math.Abs(point.X()) {
		return tuple.NewVector(point.X(), 0, 0)
	} else if maxc == math.Abs(point.Y()) {
		return tuple.NewVector(0, point.Y(), 0)
	}
	return tuple.NewVector(0, 0, point.Z())
}

func (c cube) localIntersect(ray Ray, outer Shape) []Intersection {
	xtmin, xtmax := checkAxis(ray.Origin.X(), ray.Direction.X())
	ytmin, ytmax := checkAxis(ray.Origin.Y(), ray.Direction.Y())
	ztmin, ztmax := checkAxis(ray.Origin.Z(), ray.Direction.Z())

	min := math.Max(math.Max(xtmin, ytmin), ztmin)
	max := math.Min(math.Min(xtmax, ytmax), ztmax)
	if min > max {
		return []Intersection{}
	}
	return []Intersection{
		{T: min, Shape: outer},
		{T: max, Shape: outer},
	}
}

func checkAxis(origin float64, direction float64) (float64, float64) {
	tminNum := (-1.0 - origin)
	tmaxNum := (1.0 - origin)

	var tmin, tmax float64
	if math.Abs(direction) >= utils.EPSILON {
		tmin = tminNum / direction
		tmax = tmaxNum / direction
	} else {
		var sign int
		if math.Signbit(tminNum) {
			sign = 1
		} else {
			sign = -1
		}
		tmin = math.Inf(sign)
		if math.Signbit(tmaxNum) {
			sign = 1
		} else {
			sign = -1
		}
		tmax = math.Inf(sign)
	}
	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}
	return tmin, tmax
}
