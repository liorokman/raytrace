package shapes

import (
	"fmt"
	"sort"

	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Group struct {
	content map[string]Shape
}

func (g Group) String() string {
	retval := "\n"
	for _, s := range g.content {
		retval = retval + "\t" + s.String() + "\n"
	}
	return retval
}

func NewGroup() Shape {
	return newShape(material.Default(), matrix.NewIdentity(), Group{
		content: map[string]Shape{},
	})
}

func Connect(group, child Shape) (Shape, error) {
	if g, ok := group.InnerShape().(Group); !ok {
		return nil, fmt.Errorf("Can't connect a shape to a non-group object")
	} else {
		ret := child.SetParent(group)
		g.Add(ret)
		return ret, nil
	}
}

func (g Group) Add(s Shape) {
	g.content[s.ID()] = s
}

func (g Group) Size() int {
	return len(g.content)
}

func (g Group) shapeIdPrefix() string {
	return "G"
}

func (g Group) normalAt(point tuple.Tuple, hit Intersection) tuple.Tuple {
	panic("group NormalAt should never be called")
}

func (g Group) localIntersect(ray Ray, outer Shape) []Intersection {
	retval := []Intersection{}
	for _, s := range g.content {
		retval = append(retval, ray.Intersect(s)...)
	}
	sort.Sort(ByTime(retval))
	return retval
}
