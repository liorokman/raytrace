package ray

import (
	"math"

	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Intersection struct {
	T     float64
	Shape shapes.Shape
}

type Computation struct {
	Intersection
	Point   tuple.Tuple
	EyeV    tuple.Tuple
	NormalV tuple.Tuple
	Inside  bool
}

type ByTime []Intersection

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].T < a[j].T }

func Hit(i ...Intersection) (Intersection, bool) {

	if len(i) == 0 {
		return Intersection{}, false
	}
	curr := Intersection{math.MaxFloat64, nil}
	for _, j := range i {
		if j.T >= 0 && j.T < curr.T {
			curr = j
		}
	}

	if curr.Shape != nil {
		return curr, true
	}

	return Intersection{}, false
}

func (i Intersection) PrepareComputation(r Ray) Computation {
	retval := Computation{
		Intersection: i,
		Point:        r.Position(i.T),
		EyeV:         r.Direction.Mult(-1),
	}
	retval.NormalV = i.Shape.NormalAt(retval.Point)
	// Check if the intersection happens from the inside of the shape
	if retval.NormalV.Dot(retval.EyeV) < 0 {
		retval.Inside = true
		retval.NormalV = retval.NormalV.Mult(-1)
	}

	return retval
}
