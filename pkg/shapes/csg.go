package shapes

import (
	"fmt"
	"sort"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type CSGOp int

const (
	UnionOp      CSGOp = iota
	IntersectOp  CSGOp = iota
	DifferenceOp CSGOp = iota
)

type csg struct {
	operation   CSGOp
	left, right Shape
}

func NewCSG(left, right *Shape, op CSGOp) Shape {
	r := newShape(material.Default(), matrix.NewIdentity(), csg{
		left:      *left,
		right:     *right,
		operation: op,
	})
	*left = (*left).SetParent(r)
	*right = (*right).SetParent(r)
	return r
}

func (c csg) shapeIdPrefix() string {
	return "CSG"
}

func (c csg) String() string {
	return fmt.Sprintf("left: %s, right: %s, operation: %d", c.left, c.right, c.operation)
}

func (c csg) intersectionAllowed(leftHit, inL, inR bool) bool {
	switch c.operation {
	case UnionOp:
		return (leftHit && !inR) || (!leftHit && !inL)
	case IntersectOp:
		return (leftHit && inR) || (!leftHit && inL)
	case DifferenceOp:
		return (leftHit && !inR) || (!leftHit && inL)
	default:
		return false
	}
}

func includes(a, b Shape) bool {
	switch inner := a.InnerShape().(type) {
	case csg:
		return includes(inner.left, b) || includes(inner.right, b)
	case Group:
		for _, curr := range inner.content {
			if includes(curr, b) {
				return true
			}
		}
		return false
	default:
		return (a.ID() == b.ID())
	}
}

func (c csg) filterIntersections(hits []Intersection) []Intersection {
	inL, inR := false, false
	result := []Intersection{}

	for _, i := range hits {
		lhit := includes(c.left, i.Shape)
		if c.intersectionAllowed(lhit, inL, inR) {
			result = append(result, i)
		}
		if lhit {
			inL = !inL
		} else {
			inR = !inR
		}
	}

	return result

}

func (c csg) normalAt(point tuple.Tuple, hit Intersection) tuple.Tuple {
	panic("CSG normalAt should never be called")
}

func (c csg) localIntersect(ray Ray, outer Shape) []Intersection {
	hits := ray.Intersect(c.left)
	hits = append(hits, ray.Intersect(c.right)...)
	sort.Sort(ByTime(hits))

	return c.filterIntersections(hits)
}
