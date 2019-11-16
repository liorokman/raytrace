package ray

import (
	"github.com/liorokman/raytrace/pkg/shapes"
	"math"
)

type Intersection struct {
	T     float64
	Shape shapes.Shape
}

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
