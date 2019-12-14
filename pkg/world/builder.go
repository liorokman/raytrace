package world

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type world struct {
	Objects  []object
	Fixtures []fixture
	Camera   Cam
}

const (
	//shapes
	SPHERE = "sphere"
	PLANE  = "plane"
	CUBE   = "cube"

	// patterns
	SOLID    = "solid"
	GRADIENT = "gradient"
	RING     = "ring"
	CHECKERS = "checker"

	// translations
	TRANSLATE = "translate"
	IDENTITY  = "identity"
	SCALE     = "scale"
	ROTATEX   = "rotatex"
	ROTATEY   = "rotatey"
	ROTATEZ   = "rotatez"
	SHEAR     = "shear"

	// fixtures
	POINTLIGHT = "pointlight"
)

type Cam struct {
	Hsize       uint32
	Vsize       uint32
	FieldOfView float64 `yaml:"fieldOfView"`
	From        Point
	To          Point
	Up          Vector
}

type fixture struct {
	Type     string
	Position Point `yaml:",flow"`
	Color    color `yaml:",flow"`
}

func (f fixture) toFixture() (fixtures.PointLight, error) {
	switch f.Type {
	case POINTLIGHT:
		return fixtures.NewPointLight(f.Position.ToPoint(), f.Color.toColor()), nil
	default:
		return fixtures.PointLight{}, fmt.Errorf("Unsupported fixture %s", f.Type)
	}
}

func newShape(sType string) (shapes.Shape, error) {
	switch sType {
	case SPHERE:
		return shapes.NewSphere(), nil
	case PLANE:
		return shapes.NewPlane(), nil
	case CUBE:
		return shapes.NewCube(), nil
	default:
		return nil, fmt.Errorf("Unknown shape %s", sType)
	}
}

type materialInput struct {
	Pattern         pattern
	Ambient         float64
	Diffuse         float64
	Specular        float64
	Shininess       float64
	Reflective      float64
	Transparency    float64
	RefractiveIndex float64 `yaml:"refractiveIndex"`
}

type pattern struct {
	Type      string
	Colors    []color     `yaml:",flow"`
	Transform []transform `yaml:",flow"`
}

type color [3]float64

func (c color) toColor() tuple.Color {
	return tuple.NewColor(c[0], c[1], c[2])
}

type Point [3]float64

func (c Point) ToPoint() tuple.Tuple {
	return tuple.NewPoint(c[0], c[1], c[2])
}

type Vector [3]float64

func (c Vector) ToVector() tuple.Tuple {
	return tuple.NewVector(c[0], c[1], c[2])
}

func (p pattern) toPattern() (material.Pattern, error) {
	switch p.Type {
	case SOLID:
		if len(p.Colors) != 1 {
			return material.Pattern{}, fmt.Errorf("Solid pattern requires exactly one parameter. Have %d parameters.", len(p.Colors))
		}
		return material.NewSolidPattern(p.Colors[0].toColor()), nil
	case GRADIENT:
		if len(p.Colors) != 2 {
			return material.Pattern{}, fmt.Errorf("Gradient pattern requires exactly two parameters. Have %d parameters.", len(p.Colors))
		}
		return material.NewGradientPattern(p.Colors[0].toColor(), p.Colors[1].toColor()), nil
	case RING:
		if len(p.Colors) != 2 {
			return material.Pattern{}, fmt.Errorf("Ring pattern requires exactly two parameters. Have %d parameters.", len(p.Colors))
		}
		return material.NewRingPattern(p.Colors[0].toColor(), p.Colors[1].toColor()), nil
	case CHECKERS:
		if len(p.Colors) != 2 {
			return material.Pattern{}, fmt.Errorf("Checkers pattern requires exactly two parameters. Have %d parameters.", len(p.Colors))
		}
		return material.NewCheckerPattern(p.Colors[0].toColor(), p.Colors[1].toColor()), nil
	default:
		return material.Pattern{}, fmt.Errorf("Unrecognized pattern %s", p.Type)
	}
}

type object struct {
	Type      string      `yaml:"type"`
	Transform []transform `yaml:",flow"`
	Material  materialInput
}

type transform struct {
	Type   string
	Params []float64 `yaml:",flow"`
}

func (t transform) toMatrix() (matrix.Matrix, error) {
	switch t.Type {
	case TRANSLATE:
		if len(t.Params) != 3 {
			return nil, fmt.Errorf("Translate transform requires (x,y,z) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewTranslation(t.Params[0], t.Params[1], t.Params[2]), nil
	case IDENTITY:
		if len(t.Params) != 0 {
			return nil, fmt.Errorf("Identity transform should not have parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewIdentity(), nil
	case SCALE:
		if len(t.Params) != 3 {
			return nil, fmt.Errorf("Scale transform requires (x,y,z) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewScale(t.Params[0], t.Params[1], t.Params[2]), nil
	case ROTATEX:
		if len(t.Params) != 1 {
			return nil, fmt.Errorf("RotateX transform requires (radians) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewRotateX(t.Params[0]), nil
	case ROTATEY:
		if len(t.Params) != 1 {
			return nil, fmt.Errorf("RotateY transform requires (radians) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewRotateY(t.Params[0]), nil
	case ROTATEZ:
		if len(t.Params) != 1 {
			return nil, fmt.Errorf("RotateZ transform requires (radians) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewRotateZ(t.Params[0]), nil
	case SHEAR:
		if len(t.Params) != 6 {
			return nil, fmt.Errorf("Shear transform requires (xy, xz, yx, yz, zx, zy) parameters. Have %d params instead.", len(t.Params))
		}
		return matrix.NewShear(t.Params[0], t.Params[1], t.Params[2], t.Params[3], t.Params[4], t.Params[5]), nil
	default:
		return nil, fmt.Errorf("Unsupported transform '%s'", t.Type)
	}
}

func NewWorld(file string) (*World, Cam, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, Cam{}, err
	}
	var w world
	err = yaml.Unmarshal(data, &w)
	if err != nil {
		return nil, Cam{}, err
	}
	retval := &World{
		objects: []shapes.Shape{},
		Lights:  []fixtures.PointLight{},
	}
	for _, o := range w.Objects {
		s, err := newShape(o.Type)
		if err != nil {
			return nil, Cam{}, err
		}
		finalTransform := matrix.NewIdentity()
		for _, t := range o.Transform {
			if mat, err := t.toMatrix(); err != nil {
				return nil, Cam{}, err
			} else {
				finalTransform = finalTransform.Multiply(mat)
			}
		}
		s = s.WithTransform(finalTransform)

		pat, err := o.Material.Pattern.toPattern()
		if err != nil {
			return nil, Cam{}, err
		}
		finalTransform = matrix.NewIdentity()
		for _, t := range o.Material.Pattern.Transform {
			if mat, err := t.toMatrix(); err != nil {
				return nil, Cam{}, err
			} else {
				finalTransform = finalTransform.Multiply(mat)
			}
		}
		pat = pat.WithTransform(finalTransform)
		s = s.WithMaterial(material.New(pat, o.Material.Ambient, o.Material.Diffuse,
			o.Material.Specular, o.Material.Shininess, o.Material.Reflective,
			o.Material.Transparency, o.Material.RefractiveIndex))

		retval.AddShapes(s)
	}
	for _, f := range w.Fixtures {
		fix, err := f.toFixture()
		if err != nil {
			return nil, Cam{}, err
		}
		retval.Lights = append(retval.Lights, fix)
	}

	return retval, w.Camera, nil
}
