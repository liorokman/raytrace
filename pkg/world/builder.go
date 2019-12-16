package world

import (
	"fmt"
	"io/ioutil"
	"math"

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
	SPHERE   = "sphere"
	PLANE    = "plane"
	CUBE     = "cube"
	CYLINDER = "cylinder"
	CONE     = "cone"

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

func extractFloatParam(bag map[string]interface{}, name string) (float64, bool, error) {
	if val, ok := bag[name]; ok {
		if fval, ok := val.(float64); ok {
			return fval, true, nil
		}
		return 0, false, fmt.Errorf("%s found but is not a float64 value", name)
	} else {
		return 0, false, nil
	}
}

func newShape(sType string, params map[string]interface{}) (shapes.Shape, error) {
	switch sType {
	case SPHERE:
		return shapes.NewSphere(), nil
	case PLANE:
		return shapes.NewPlane(), nil
	case CUBE:
		return shapes.NewCube(), nil
	case CYLINDER:
		fallthrough
	case CONE:
		closed := false
		min := math.Inf(-1)
		max := math.Inf(1)
		if val, ok := params["closed"]; ok {
			if boolVal, ok := val.(bool); ok {
				closed = boolVal
			} else {
				return nil, fmt.Errorf("Closed param for cylinder/cone should be a bool")
			}
		}
		val, ok, err := extractFloatParam(params, "minimum")
		if err != nil {
			return nil, err
		} else if ok {
			min = val
		}
		val, ok, err = extractFloatParam(params, "maximum")
		if err != nil {
			return nil, err
		} else if ok {
			max = val
		}
		if sType == CONE {
			return shapes.NewConstrainedCone(min, max, closed), nil
		} else {
			return shapes.NewConstrainedCylinder(min, max, closed), nil
		}
	default:
		return nil, fmt.Errorf("Unknown shape %s", sType)
	}
}

const (
	AMBIENT         = "ambient"
	DIFFUSE         = "diffuse"
	SPECULAR        = "specular"
	SHININESS       = "shininess"
	REFLECTIVE      = "reflective"
	TRANSPARENCY    = "transparency"
	REFRACTIVEINDEX = "refractiveIndex"

	DEFAULTMATERIAL = "default"
	GLASSMATERIAL   = "glass"
)

type materialInput struct {
	Pattern pattern
	Params  map[string]interface{} `yaml:",inline"`
}

func (m materialInput) toMaterial() (material.Material, error) {
	pat, err := m.Pattern.toPattern()
	if err != nil {
		return material.Material{}, err
	}
	finalTransform := matrix.NewIdentity()
	for _, t := range m.Pattern.Transform {
		if mat, err := t.toMatrix(); err != nil {
			return material.Material{}, err
		} else {
			finalTransform = finalTransform.Multiply(mat)
		}
	}
	pat = pat.WithTransform(finalTransform)
	mb := material.NewBuilder(material.Default())
	if matType, ok := m.Params["preset"]; ok {
		if str, ok := matType.(string); ok {
			switch str {
			case DEFAULTMATERIAL:
				mb = material.NewBuilder(material.Default())
			case GLASSMATERIAL:
				mb = material.NewBuilder(material.Glass())
			default:
				mb = material.NewBuilder(material.Default())
			}
		}
	}
	mb = mb.WithPattern(pat)
	for k := range m.Params {
		switch k {
		case AMBIENT:
			if val, ok, err := extractFloatParam(m.Params, AMBIENT); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithAmbient(val)
			}
		case DIFFUSE:
			if val, ok, err := extractFloatParam(m.Params, DIFFUSE); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithDiffuse(val)
			}
		case SPECULAR:
			if val, ok, err := extractFloatParam(m.Params, SPECULAR); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithSpecular(val)
			}
		case SHININESS:
			if val, ok, err := extractFloatParam(m.Params, SHININESS); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithShininess(val)
			}
		case REFLECTIVE:
			if val, ok, err := extractFloatParam(m.Params, REFLECTIVE); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithReflective(val)
			}
		case TRANSPARENCY:
			if val, ok, err := extractFloatParam(m.Params, TRANSPARENCY); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithTransparency(val)
			}
		case REFRACTIVEINDEX:
			if val, ok, err := extractFloatParam(m.Params, REFRACTIVEINDEX); err != nil {
				return material.Material{}, err
			} else if ok {
				mb.WithRefractiveIndex(val)
			}
		}
	}

	return mb.Build(), nil
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
	Params    map[string]interface{} `yaml:"params,omitempty"`
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
		s, err := newShape(o.Type, o.Params)
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
		mat, err := o.Material.toMaterial()
		if err != nil {
			return nil, Cam{}, err
		}
		s = s.WithMaterial(mat)

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
