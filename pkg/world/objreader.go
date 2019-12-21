package world

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type objReader struct {
	ignoredLines int
	vertices     []tuple.Tuple
	groups       map[string]shapes.Shape
}

func newObjReader() *objReader {
	return &objReader{
		vertices: []tuple.Tuple{tuple.NewPoint(math.NaN(), math.NaN(), math.NaN())},
		groups: map[string]shapes.Shape{
			"defaultGroup": shapes.NewGroup(),
		},
	}
}

const (
	vertex = "v"
	face   = "f"
	group  = "g"
)

func toFloat64Slice(in []string) ([]float64, error) {
	r := make([]float64, len(in))
	for i := range in {
		var err error
		r[i], err = strconv.ParseFloat(in[i], 64)
		if err != nil {
			return []float64{}, err
		}
	}
	return r, nil
}

func toIntSlice(in []string) ([]int, error) {
	r := make([]int, len(in))
	for i := range in {
		var err error
		r[i], err = strconv.Atoi(in[i])
		if err != nil {
			return []int{}, err
		}
	}
	return r, nil
}

func (o *objReader) AsGroup() shapes.Shape {
	g := shapes.NewGroup()
	for i := range o.groups {
		shapes.Connect(g, o.groups[i])
	}
	return g
}

func (o *objReader) ReadObj(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}

	currentGroup := "defaultGroup"
	scan := bufio.NewScanner(in)
	for scan.Scan() {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		// TODO: remove comments from after the first position
		parts := strings.Split(line, " ")
		switch parts[0] {
		case vertex:
			coords, err := toFloat64Slice(parts[1:])
			if err != nil {
				return err
			}
			if len(coords) != 3 {
				o.ignoredLines++
				continue
			}
			o.vertices = append(o.vertices, tuple.NewPoint(coords[0], coords[1], coords[2]))
		case face:
			vertices, err := toIntSlice(parts[1:])
			if err != nil {
				return err
			}
			if len(vertices) < 3 {
				o.ignoredLines++
				continue
			}
			currGroup := o.groups[currentGroup]
			for i := 1; i < len(vertices)-1; i++ {
				var err error
				tri := shapes.NewTriangle(
					o.vertices[vertices[0]],
					o.vertices[vertices[i]],
					o.vertices[vertices[i+1]])
				if _, err = shapes.Connect(currGroup, tri); err != nil {
					return err
				}
			}
			o.groups[currentGroup] = currGroup
		case group:
			o.groups[parts[1]] = shapes.NewGroup()
			currentGroup = parts[1]
		}

	}
	if err := scan.Err(); err != nil {
		return err
	}
	return nil
}
