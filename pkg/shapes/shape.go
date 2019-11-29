package shapes

import (
	"fmt"
	"sync/atomic"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Shape interface {
	ID() string
	GetTransform() matrix.Matrix
	GetMaterial() material.Material

	WithTransform(matrix.Matrix) Shape
	WithMaterial(material.Material) Shape
	NormalAt(tuple.Tuple) tuple.Tuple
	LocalIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64
}

type ShapeList []Shape

func (s ShapeList) Find(needle Shape) int {
	for i := range s {
		if s[i].ID() == needle.ID() {
			return i
		}
	}
	return -1
}

type shapeDetails interface {
	shapeIdPrefix() string
	normalAt(tuple.Tuple) tuple.Tuple
	localIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64
}

var shapeCounter int32 = 0

func GetMaxCreatedShapes() int32 {
	return atomic.LoadInt32(&shapeCounter)
}

type shapeCore struct {
	id        int32
	transform matrix.Matrix
	material  material.Material
	shape     shapeDetails
}

func (s shapeCore) GetTransform() matrix.Matrix {
	return s.transform
}

func (s shapeCore) GetMaterial() material.Material {
	return s.material
}

func newShape(m material.Material, t matrix.Matrix, s shapeDetails) shapeCore {
	return shapeCore{
		id:        atomic.AddInt32(&shapeCounter, 1),
		material:  m,
		transform: t,
		shape:     s,
	}
}

func (s shapeCore) NormalAt(point tuple.Tuple) tuple.Tuple {
	if !point.IsPoint() {
		panic("Can't compute a normal at a Vector")
	}
	shapeInverseTransform, err := s.GetTransform().Inverse()
	if err != nil {
		panic(err)
	}
	objPoint := shapeInverseTransform.MultiplyTuple(point)

	objNormal := s.shape.normalAt(objPoint)

	worldNormal := shapeInverseTransform.Transpose().MultiplyTuple(objNormal)
	worldNormal[tuple.WPos] = 0
	return worldNormal.Normalize()
}

func (s shapeCore) LocalIntersect(direction tuple.Tuple, origin tuple.Tuple) []float64 {
	return s.shape.localIntersect(direction, origin)
}

func (s shapeCore) ID() string {
	return fmt.Sprintf("%s%d", s.shape.shapeIdPrefix(), s.id)
}

func (s shapeCore) WithTransform(t matrix.Matrix) Shape {
	return newShape(s.material, t, s.shape)
}

func (s shapeCore) WithMaterial(m material.Material) Shape {
	return newShape(m, s.transform, s.shape)
}
