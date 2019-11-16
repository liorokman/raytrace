package shapes

import (
	"github.com/liorokman/raytrace/pkg/matrix"
)

type Shape interface {
	ID() string
	GetTransform() matrix.Matrix

	WithTransform(matrix.Matrix) Shape
}
