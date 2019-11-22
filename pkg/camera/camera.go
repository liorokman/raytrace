package camera

import (
	"math"

	"github.com/liorokman/raytrace/pkg/canvas"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/ray"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/world"
)

func ViewTransformation(from, to, up tuple.Tuple) matrix.Matrix {

	forward := to.Subtract(from).Normalize()
	upn := up.Normalize()
	left := forward.Cross(upn)
	trueUp := left.Cross(forward)

	orientation := matrix.Matrix{
		{left.X(), left.Y(), left.Z(), 0},
		{trueUp.X(), trueUp.Y(), trueUp.Z(), 0},
		{-forward.X(), -forward.Y(), -forward.Z(), 0},
		{0, 0, 0, 1},
	}
	return orientation.Multiply(matrix.NewTranslation(-from.X(), -from.Y(), -from.Z()))
}

type Camera struct {
	hsize       uint32
	vsize       uint32
	fieldOfView float64
	transform   matrix.Matrix

	halfWidth  float64
	halfHeight float64
	pixelSize  float64
}

func NewCamera(hsize, vsize uint32, fieldOfView float64) Camera {
	cam := Camera{
		hsize:       hsize,
		vsize:       vsize,
		fieldOfView: fieldOfView,
		transform:   matrix.NewIdentity(),
	}
	halfView := math.Tan(fieldOfView / 2.0)
	aspect := float64(hsize) / float64(vsize)
	if aspect >= 1 {
		cam.halfWidth = halfView
		cam.halfHeight = halfView / aspect
	} else {
		cam.halfWidth = halfView * aspect
		cam.halfHeight = halfView
	}
	cam.pixelSize = (cam.halfWidth * 2.0) / float64(cam.hsize)
	return cam
}

func (c Camera) WithTransform(t matrix.Matrix) Camera {
	return Camera{
		hsize:       c.hsize,
		vsize:       c.vsize,
		fieldOfView: c.fieldOfView,
		transform:   t,
		halfWidth:   c.halfWidth,
		halfHeight:  c.halfHeight,
		pixelSize:   c.pixelSize,
	}
}

func (c Camera) RayForPixel(px, py uint32) ray.Ray {

	xOffset := (float64(px) + 0.5) * c.pixelSize
	yOffset := (float64(py) + 0.5) * c.pixelSize

	worldX := c.halfWidth - xOffset
	worldY := c.halfHeight - yOffset

	transformInverse, err := c.transform.Inverse()
	if err != nil {
		panic(err)
	}
	pixel := transformInverse.MultiplyTuple(tuple.NewPoint(worldX, worldY, -1))
	origin := transformInverse.MultiplyTuple(tuple.NewPoint(0, 0, 0))
	direction := pixel.Subtract(origin).Normalize()
	ray, err := ray.New(origin, direction)
	if err != nil {
		panic(err)
	}
	return ray
}

func (c Camera) Render(w *world.World) canvas.Canvas {
	image := canvas.New(c.hsize, c.vsize)
	for y := uint32(0); y < c.hsize; y++ {
		for x := uint32(0); x < c.vsize; x++ {
			ray := c.RayForPixel(x, y)
			color := w.ColorAt(ray)
			image.SetPixel(x, y, color)
		}
	}
	return image
}

func (c Camera) HSize() uint32 {
	return c.hsize
}

func (c Camera) VSize() uint32 {
	return c.vsize
}

func (c Camera) FieldOfView() float64 {
	return c.fieldOfView
}

func (c Camera) Transform() matrix.Matrix {
	return c.transform
}

func (c Camera) HalfWidth() float64 {
	return c.halfWidth
}

func (c Camera) HalfHeight() float64 {
	return c.halfHeight
}

func (c Camera) PixelSize() float64 {
	return c.pixelSize
}
