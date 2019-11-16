package shapes

import (
	"fmt"

	"github.com/liorokman/raytrace/pkg/matrix"
)

const sphereId = "S"

var sphereCounter = 0

type Sphere struct {
	id        int
	transform matrix.Matrix
}

func NewSphere() Sphere {
	sphereCounter++
	return Sphere{
		id:        sphereCounter,
		transform: matrix.NewIdentity(),
	}
}

func (s Sphere) ID() string {
	return fmt.Sprintf("%s%d", sphereId, s.id)
}

func (s Sphere) GetTransform() matrix.Matrix {
	return s.transform
}

func (s Sphere) WithTransform(t matrix.Matrix) Shape {
	res := NewSphere()
	res.transform = t
	return res
}
