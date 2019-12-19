package types

import (
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type ShapeTransformer interface {
	GetTransform() matrix.Matrix
	WorldToObject(point tuple.Tuple) (tuple.Tuple, error)
}
