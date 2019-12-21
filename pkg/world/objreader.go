package world

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type objReader struct {
	ignoredLines   int
	vertices       []tuple.Tuple
	verticeNormals []tuple.Tuple
	groups         map[string]shapes.Shape
}

func newObjReader() *objReader {
	return &objReader{
		vertices:       []tuple.Tuple{tuple.NewPoint(math.NaN(), math.NaN(), math.NaN())},
		verticeNormals: []tuple.Tuple{tuple.NewVector(math.NaN(), math.NaN(), math.NaN())},
		groups: map[string]shapes.Shape{
			"defaultGroup": shapes.NewGroup(),
		},
	}
}

const (
	vertexInObj  = "v"
	vertexNormal = "vn"
	faceInObj    = "f"
	groupInObj   = "g"
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

func toIntSlice(in []string, section int) ([]int, error) {
	r := make([]int, len(in))
	for i := range in {
		var err error
		parts := strings.Split(in[i], "/")
		if len(parts) < section {
			return []int{}, nil
		}
		r[i], err = strconv.Atoi(parts[section])
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

	whitespaceSqueezer := regexp.MustCompile("(\\s)\\s*")
	currentGroup := "defaultGroup"
	scan := bufio.NewScanner(in)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if ind := strings.Index(line, "#"); ind >= 0 {
			line = line[0:ind]
		}
		line = whitespaceSqueezer.ReplaceAllString(line, "$1")
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		switch parts[0] {
		case vertexInObj, vertexNormal:
			coords, err := toFloat64Slice(parts[1:])
			if err != nil {
				return err
			}
			if len(coords) != 3 {
				o.ignoredLines++
				continue
			}
			if parts[0] == vertexInObj {
				o.vertices = append(o.vertices, tuple.NewPoint(coords[0], coords[1], coords[2]))
			} else {
				o.verticeNormals = append(o.verticeNormals, tuple.NewPoint(coords[0], coords[1], coords[2]))
			}
		case faceInObj:
			vertices, err := toIntSlice(parts[1:], 0)
			if err != nil {
				return err
			}
			if len(vertices) < 3 {
				o.ignoredLines++
				continue
			}
			vNormals, err := toIntSlice(parts[1:], 2)
			if err != nil {
				return err
			}
			currGroup := o.groups[currentGroup]
			for i := 1; i < len(vertices)-1; i++ {
				var err error
				var tri shapes.Shape
				if len(vNormals) == 0 {
					tri = shapes.NewTriangle(
						o.vertices[vertices[0]],
						o.vertices[vertices[i]],
						o.vertices[vertices[i+1]])
				} else {
					tri = shapes.NewSmoothTriangle(
						o.vertices[vertices[0]],
						o.vertices[vertices[i]],
						o.vertices[vertices[i+1]],
						o.verticeNormals[vNormals[0]],
						o.verticeNormals[vNormals[i]],
						o.verticeNormals[vNormals[i+1]],
					)

				}
				if _, err = shapes.Connect(currGroup, tri); err != nil {
					return err
				}
			}
			o.groups[currentGroup] = currGroup
		case groupInObj:
			o.groups[parts[1]] = shapes.NewGroup()
			currentGroup = parts[1]
		}

	}
	if err := scan.Err(); err != nil {
		return err
	}
	return nil
}
