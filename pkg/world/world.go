package world

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type World struct {
	objects []shapes.Shape
	Lights  []fixtures.PointLight
}

func New() *World {
	return &World{
		objects: []shapes.Shape{},
		Lights:  []fixtures.PointLight{fixtures.NewPointLight(tuple.NewPoint(-10, 10, -10), tuple.NewColor(1, 1, 1))},
	}
}

func (w *World) String() string {
	retval := fmt.Sprintf("Lights: %#v\n", w.Lights)
	for _, o := range w.objects {
		retval = retval + "\n" + o.String()
	}
	return retval
}

func (w *World) NumObjects() int {
	return len(w.objects)
}

func (w *World) Shape(i int) shapes.Shape {
	if i < len(w.objects) {
		return w.objects[i]
	}
	return nil
}

func (w *World) SetShape(i int, s shapes.Shape) *World {
	if i >= len(w.objects) {
		fmt.Printf("Warn: attempt to modify a non-existant shape %d\n", i)
		return w
	}
	w.objects[i] = s
	return w
}

func (w *World) AddShapes(s ...shapes.Shape) *World {
	w.objects = append(w.objects, s...)
	return w
}

func (w *World) IntersectRay(r shapes.Ray) []shapes.Intersection {
	retval := []shapes.Intersection{}
	for _, o := range w.objects {
		retval = append(retval, r.Intersect(o)...)
	}

	sort.Sort(shapes.ByTime(retval))
	return retval
}

func (w *World) ShadeHit(comps shapes.Computation, depth int) tuple.Color {
	wg := sync.WaitGroup{}
	colorFromL := make([]tuple.Color, len(w.Lights))
	for i, l := range w.Lights {
		wg.Add(1)
		go func(ind int, light fixtures.PointLight) {
			shadowed := w.IsShadowed(comps.OverPoint, ind)
			colorFromL[ind] = comps.Shape.GetMaterial().Lighting(comps.Shape, light, comps.Point, comps.EyeV, comps.NormalV, shadowed)
			// TODO: Actually do something with these errors
			reflect, _ := w.ReflectedColor(comps, depth)
			refract, _ := w.RefractedColor(comps, depth)
			if comps.Shape.GetMaterial().Reflective() > 0.0 && comps.Shape.GetMaterial().Transparency() > 0.0 {
				reflectence := comps.Schlick()
				colorFromL[ind] = colorFromL[ind].Add(reflect.Mult(reflectence)).Add(refract.Mult((1 - reflectence)))
			} else {
				colorFromL[ind] = colorFromL[ind].Add(reflect).Add(refract)
			}
			wg.Done()
		}(i, l)
	}
	wg.Wait()
	retval := tuple.Black
	for _, c := range colorFromL {
		retval = retval.Add(c)
	}
	return retval
}

func (w *World) RefractedColor(comps shapes.Computation, depth int) (tuple.Color, error) {
	if depth == 0 || comps.Shape.GetMaterial().Transparency() == 0 {
		return tuple.Black, nil
	}
	nRatio := comps.N1 / comps.N2
	cosI := comps.EyeV.Dot(comps.NormalV)
	sin2t := math.Pow(nRatio, 2) * (1 - math.Pow(cosI, 2))
	if sin2t > 1 {
		// This means that this case is total internal reflection
		return tuple.Black, nil
	}
	cosT := math.Sqrt(1.0 - sin2t)
	direction := comps.NormalV.Mult(nRatio*cosI - cosT).Subtract(comps.EyeV.Mult(nRatio))
	refractRay, err := shapes.NewRay(comps.UnderPoint, direction)
	if err != nil {
		return tuple.Color{}, err
	}
	c, err := w.ColorAt(refractRay, depth-1)
	if err != nil {
		return tuple.Color{}, err
	}
	return c.Mult(comps.Shape.GetMaterial().Transparency()), nil
}

func (w *World) ReflectedColor(comps shapes.Computation, depth int) (tuple.Color, error) {
	if comps.Shape.GetMaterial().Reflective() == 0 || depth <= 0 {
		return tuple.Black, nil
	}
	reflectedRay, err := shapes.NewRay(comps.OverPoint, comps.ReflectV)
	if err != nil {
		return tuple.Color{}, err
	}
	c, err := w.ColorAt(reflectedRay, depth-1)
	if err != nil {
		return tuple.Color{}, err
	}
	return c.Mult(comps.Shape.GetMaterial().Reflective()), nil
}

func (w *World) ColorAt(r shapes.Ray, depth int) (tuple.Color, error) {
	xs := w.IntersectRay(r)
	if h, ok := shapes.Hit(xs...); ok {
		comps, err := h.PrepareComputation(r, xs...)
		if err != nil {
			return tuple.Color{}, err
		}
		return w.ShadeHit(comps, depth), nil
	} else {
		return tuple.Black, nil
	}
}

func (w *World) IsShadowed(p tuple.Tuple, lightIndex int) bool {
	if !p.IsPoint() {
		panic("Expecting a point, not a vector")
	}
	if lightIndex < 0 || lightIndex >= len(w.Lights) {
		panic("No such light source in world")
	}
	v := w.Lights[lightIndex].Position().Subtract(p)
	direction := v.Normalize()

	r, err := shapes.NewRay(p, direction)
	if err != nil {
		panic(err)
	}
	intersections := w.IntersectRay(r)
	if hit, ok := shapes.Hit(intersections...); ok {
		distance := v.Magnitude()
		return hit.T < distance
	}
	return false
}
