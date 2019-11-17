package shapes

import (
	"fmt"
	"sync/atomic"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

const sphereId = "S"

var sphereCounter int32 = 0

type Sphere struct {
	id        int32
	transform matrix.Matrix
	material  material.Material
}

func NewSphere() Sphere {

	return Sphere{
		id:        atomic.AddInt32(&sphereCounter, 1),
		transform: matrix.NewIdentity(),
		material:  material.Default(),
	}
}

func (s Sphere) ID() string {
	return fmt.Sprintf("%s%d", sphereId, s.id)
}

func (s Sphere) GetTransform() matrix.Matrix {
	return s.transform
}

func (s Sphere) GetMaterial() material.Material {
	return s.material
}

func (s Sphere) WithTransform(t matrix.Matrix) Shape {
	res := Sphere{
		id:        atomic.AddInt32(&sphereCounter, 1),
		transform: t,
		material:  s.material,
	}
	return res
}

func (s Sphere) WithMaterial(m material.Material) Shape {
	res := Sphere{
		id:        atomic.AddInt32(&sphereCounter, 1),
		transform: s.transform,
		material:  m,
	}
	return res
}

func (s Sphere) NormalAt(point tuple.Tuple) tuple.Tuple {
	if !point.IsPoint() {
		panic("Can't compute a normal at a Vector")
	}
	shapeInverseTransform, err := s.GetTransform().Inverse()
	if err != nil {
		panic(err)
	}
	objPoint := shapeInverseTransform.MultiplyTuple(point)
	objNormal := objPoint.Subtract(tuple.NewPoint(0, 0, 0))
	worldNormal := shapeInverseTransform.Transpose().MultiplyTuple(objNormal)
	worldNormal[tuple.WPos] = 0
	return worldNormal.Normalize()
}
