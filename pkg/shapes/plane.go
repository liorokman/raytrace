package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type plane struct {
}

func NewPlane() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), plane{})
}

func (p plane) shapeIdPrefix() string {
	return "P"
}

func (p plane) normalAt(tuple.Tuple) tuple.Tuple {
	return tuple.NewVector(0, 1, 0)
}

func (p plane) localIntersect(ray Ray, outer Shape) []Intersection {
	if math.Abs(ray.Direction.Y()) < utils.EPSILON {
		return []Intersection{}
	}
	return []Intersection{
		{T: -ray.Origin.Y() / ray.Direction.Y(), Shape: outer},
	}
}
