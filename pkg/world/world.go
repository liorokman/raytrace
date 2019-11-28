package world

import (
	"fmt"
	"sort"
	"sync"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/ray"
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

func (w *World) IntersectRay(r ray.Ray) []ray.Intersection {
	retval := []ray.Intersection{}
	for _, o := range w.objects {
		retval = append(retval, r.Intersect(o)...)
	}

	sort.Sort(ray.ByTime(retval))
	return retval
}

func (w *World) ShadeHit(comps ray.Computation, depth int) tuple.Color {
	wg := sync.WaitGroup{}
	colorFromL := make([]tuple.Color, len(w.Lights))
	for i, l := range w.Lights {
		wg.Add(1)
		go func(ind int, light fixtures.PointLight) {
			shadowed := w.IsShadowed(comps.OverPoint, ind)
			colorFromL[ind] = comps.Shape.GetMaterial().Lighting(comps.Shape.GetTransform(), light, comps.Point, comps.EyeV, comps.NormalV, shadowed)
			colorFromL[ind] = colorFromL[ind].Add(w.ReflectedColor(comps, depth))
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

func (w *World) ReflectedColor(comps ray.Computation, depth int) tuple.Color {
	if comps.Shape.GetMaterial().Reflective() == 0 || depth <= 0 {
		return tuple.Black
	}
	reflectedRay, err := ray.New(comps.OverPoint, comps.ReflectV)
	if err != nil {
		panic(err)
	}
	c := w.ColorAt(reflectedRay, depth-1)
	return c.Mult(comps.Shape.GetMaterial().Reflective())
}

func (w *World) ColorAt(r ray.Ray, depth int) tuple.Color {
	xs := w.IntersectRay(r)
	if h, ok := ray.Hit(xs...); ok {
		comps := h.PrepareComputation(r)
		return w.ShadeHit(comps, depth)
	} else {
		return tuple.Black
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

	r, err := ray.New(p, direction)
	if err != nil {
		panic(err)
	}
	intersections := w.IntersectRay(r)
	if hit, ok := ray.Hit(intersections...); ok {
		distance := v.Magnitude()
		return hit.T < distance
	}
	return false
}
