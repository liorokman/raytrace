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

func (p plane) localIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64 {
	if math.Abs(direction.Y()) < utils.EPSILON {
		return []float64{}
	}
	return []float64{-origin.Y() / direction.Y()}
}
