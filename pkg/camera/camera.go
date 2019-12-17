package camera

import (
	"math"
	"runtime"
	"sync"

	"github.com/liorokman/raytrace/pkg/canvas"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
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

func (c Camera) RayForPixel(px, py uint32) shapes.Ray {

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
	ray, err := shapes.NewRay(origin, direction)
	if err != nil {
		panic(err)
	}
	return ray
}

type unitOfWork struct {
	x, y uint32
}

type queue chan unitOfWork

func (c Camera) Render(w *world.World) canvas.Canvas {
	image := canvas.New(c.hsize, c.vsize)

	wg := sync.WaitGroup{}
	q := make(queue, runtime.NumCPU())

	for cpu := 0; cpu < runtime.NumCPU()/len(w.Lights); cpu++ {
		wg.Add(1)
		go func() {
			for {
				unit, ok := <-q
				if !ok {
					wg.Done()
					return
				}
				ray := c.RayForPixel(unit.x, unit.y)
				color := w.ColorAt(ray, 4)
				image.SetPixel(unit.x, unit.y, color)
			}
		}()
	}

	for y := uint32(0); y < c.hsize; y++ {
		for x := uint32(0); x < c.vsize; x++ {
			q <- unitOfWork{x: x, y: y}
		}
	}
	close(q)
	wg.Wait()
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
