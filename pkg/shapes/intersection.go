package shapes

import (
	"math"

	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/utils"
)

type Intersection struct {
	T     float64
	Shape Shape
}

type Computation struct {
	Intersection
	Point      tuple.Tuple
	EyeV       tuple.Tuple
	NormalV    tuple.Tuple
	OverPoint  tuple.Tuple
	UnderPoint tuple.Tuple
	Inside     bool
	ReflectV   tuple.Tuple
	N1, N2     float64
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

func (i Intersection) Equals(other Intersection) bool {
	return utils.FloatEqual(i.T, other.T) && i.Shape.ID() == other.Shape.ID()
}

func (i Intersection) PrepareComputation(r Ray, xs ...Intersection) Computation {
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
	retval.OverPoint = retval.Point.Add(retval.NormalV.Mult(utils.EPSILON))
	retval.UnderPoint = retval.Point.Subtract(retval.NormalV.Mult(utils.EPSILON))
	retval.ReflectV = r.Direction.Reflect(retval.NormalV)

	containers := ShapeList{}
	for currXS := range xs {
		if xs[currXS].Equals(i) {
			if len(containers) == 0 {
				retval.N1 = 1.0
			} else {
				retval.N1 = containers[len(containers)-1].GetMaterial().RefractiveIndex()
			}
		}
		if place := containers.Find(xs[currXS].Shape); place >= 0 {
			copy(containers[place:], containers[place+1:])
			containers = containers[:len(containers)-1]
		} else {
			containers = append(containers, xs[currXS].Shape)
		}
		if xs[currXS].Equals(i) {
			if len(containers) == 0 {
				retval.N2 = 1.0
			} else {
				retval.N2 = containers[len(containers)-1].GetMaterial().RefractiveIndex()
			}
			break
		}
	}
	return retval
}

func (c Computation) Schlick() float64 {
	cos := c.EyeV.Dot(c.NormalV)
	if c.N1 > c.N2 {
		n := c.N1 / c.N2
		sin2t := n * n * (1.0 - cos*cos)
		if sin2t > 1.0 {
			return 1.0
		}
		cos = math.Sqrt(1.0 - sin2t)
	}
	r0 := math.Pow((c.N1-c.N2)/(c.N1+c.N2), 2)
	return r0 + (1-r0)*math.Pow(1-cos, 5)
}
