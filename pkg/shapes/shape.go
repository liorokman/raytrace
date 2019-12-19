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
	NormalAt(tuple.Tuple) (tuple.Tuple, error)
	LocalIntersect(ray Ray) []Intersection
	WorldToObject(point tuple.Tuple) (tuple.Tuple, error)
	NormalToWorld(vector tuple.Tuple) (tuple.Tuple, error)

	SetParent(p Shape) Shape
	Parent() Shape
	InnerShape() ShapeDetails

	String() string
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

func (s shapeCore) String() string {
	return fmt.Sprintf("ID: %s, transform: %#v, material: %#v, shape details: %s\n", s.ID(), s.transform, s.material, s.shape)
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

func (s shapeCore) NormalAt(point tuple.Tuple) (tuple.Tuple, error) {
	if !point.IsPoint() {
		return tuple.Tuple{}, fmt.Errorf("Can't compute a normal at a Vector")
	}

	localPoint, err := s.WorldToObject(point)
	if err != nil {
		return tuple.Tuple{}, err
	}
	localNormal := s.shape.normalAt(localPoint)
	normal, err := s.NormalToWorld(localNormal)
	if err != nil {
		return tuple.Tuple{}, err
	}
	return normal, nil
}

func (s shapeCore) NormalToWorld(vector tuple.Tuple) (tuple.Tuple, error) {
	inv, err := s.transform.Inverse()
	if err != nil {
		return tuple.Tuple{}, err
	}
	retval := inv.Transpose().MultiplyTuple(vector)
	retval[tuple.WPos] = 0
	retval = retval.Normalize()

	if s.Parent() != nil {
		retval, err = s.Parent().NormalToWorld(retval)
		if err != nil {
			return tuple.Tuple{}, err
		}
	}
	return retval, nil
}

func (s shapeCore) WorldToObject(point tuple.Tuple) (tuple.Tuple, error) {
	var err error
	retval := point
	if s.Parent() != nil {
		retval, err = s.Parent().WorldToObject(point)
		if err != nil {
			return retval, err
		}
	}
	inv, err := s.transform.Inverse()
	if err != nil {
		return retval, err
	}
	return inv.MultiplyTuple(retval), nil
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
