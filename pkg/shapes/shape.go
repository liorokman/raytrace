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
	LocalIntersect(ray Ray) []Intersection

	SetParent(p Shape) Shape
	Parent() Shape
	InnerShape() ShapeDetails
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

type ShapeDetails interface {
	shapeIdPrefix() string
	normalAt(tuple.Tuple) tuple.Tuple
	localIntersect(ray Ray, outer Shape) []Intersection
}

var shapeCounter int32 = 0

func GetMaxCreatedShapes() int32 {
	return atomic.LoadInt32(&shapeCounter)
}

type shapeCore struct {
	id        int32
	transform matrix.Matrix
	material  material.Material
	shape     ShapeDetails
	parent    Shape
}

func (s shapeCore) InnerShape() ShapeDetails {
	return s.shape
}

func (s shapeCore) SetParent(p Shape) Shape {
	s.parent = p
	return s
}

func (s shapeCore) Parent() Shape {
	return s.parent
}

func (s shapeCore) GetTransform() matrix.Matrix {
	return s.transform
}

func (s shapeCore) GetMaterial() material.Material {
	return s.material
}

func newShape(m material.Material, t matrix.Matrix, s ShapeDetails) shapeCore {
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

func (s shapeCore) LocalIntersect(ray Ray) []Intersection {
	return s.shape.localIntersect(ray, s)
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
