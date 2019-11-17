package shapes

import (
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Shape interface {
	ID() string
	GetTransform() matrix.Matrix
	GetMaterial() material.Material

	WithTransform(matrix.Matrix) Shape
	NormalAt(tuple.Tuple) tuple.Tuple
}
