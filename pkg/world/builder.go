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
	Objects   []object
	Fixtures  []fixture
	Materials []materialInput
	Camera    Cam
}

const (
	//shapes
	SPHERE   = "sphere"
	PLANE    = "plane"
	CUBE     = "cube"
	CYLINDER = "cylinder"
	CONE     = "cone"
	GROUP    = "group"
	TRIANGLE = "triangle"

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

func extractFloatSliceParam(bag map[string]interface{}, name string) ([]float64, bool, error) {
	if val, ok := bag[name]; ok {
		if arr, ok := val.([]interface{}); ok {
			r := make([]float64, len(arr))
			for i := range arr {
				switch fval := arr[i].(type) {
				case float64:
					r[i] = fval
				case int:
					r[i] = float64(fval)
				}
			}
			return r, true, nil
		} else {
			return []float64{}, false, fmt.Errorf("%s found but is not an array", name)
		}
	} else {
		return []float64{}, false, nil
	}
}

func sliceToPoint(s []float64) (Point, error) {

	if len(s) == 3 {
		return Point{s[0], s[1], s[2]}, nil
	}
	return Point{}, fmt.Errorf("Mismatched number of indices for a point")
}

func extractFloatParam(bag map[string]interface{}, name string) (float64, bool, error) {
	if val, ok := bag[name]; ok {
		switch fval := val.(type) {
		case float64:
			return fval, true, nil
		case int:
			return float64(fval), true, nil
		}
		return 0, false, fmt.Errorf("%s found but is not a float64 value", name)
	} else {
		return 0, false, nil
	}
}

func newShape(sType string, params map[string]interface{}, cache map[string]material.Material) (shapes.Shape, error) {
	switch sType {
	case SPHERE:
		return shapes.NewSphere(), nil
	case PLANE:
		return shapes.NewPlane(), nil
	case CUBE:
		return shapes.NewCube(), nil
	case TRIANGLE:
		p1Raw, ok, err := extractFloatSliceParam(params, "p1")
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, fmt.Errorf("triangle needs a p1 coordinate, none supplied")
		}
		p1, err := sliceToPoint(p1Raw)
		if err != nil {
			return nil, err
		}
		p2Raw, ok, err := extractFloatSliceParam(params, "p2")
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, fmt.Errorf("triangle needs a p1 coordinate, none supplied")
		}
		p2, err := sliceToPoint(p2Raw)
		if err != nil {
			return nil, err
		}

		p3Raw, ok, err := extractFloatSliceParam(params, "p3")
		if err != nil {
			return nil, err
		} else if !ok {
			return nil, fmt.Errorf("triangle needs a p1 coordinate, none supplied")
		}
		p3, err := sliceToPoint(p3Raw)
		if err != nil {
			return nil, err
		}
		return shapes.NewTriangle(p1.ToPoint(), p2.ToPoint(), p3.ToPoint()), nil

	case GROUP:
		g := shapes.NewGroup()
		if val, ok := params["content"]; ok {
			asYaml, _ := yaml.Marshal(val)
			var content []object = []object{}
			if err := yaml.Unmarshal(asYaml, &content); err != nil {
				return nil, err
			}
			for _, o := range content {
				s, err := translater(o, cache)
				if err != nil {
					return nil, err
				}
				if _, err = shapes.Connect(g, s); err != nil {
					return nil, err
				}
			}
		} else if val, ok := params["objfile"]; ok {
			if strVal, ok := val.(string); ok {
				objIn := newObjReader()
				if err := objIn.ReadObj(strVal); err != nil {
					return nil, err
				}
				return objIn.AsGroup(), nil
			} else {
				return nil, fmt.Errorf("group parameter objfile isn't a string")
			}
		}
		return g, nil
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

func (m materialInput) toMaterial(cache map[string]material.Material) (material.Material, error) {
	mb := material.NewBuilder(material.Default())
	if matType, ok := m.Params["preset"]; ok {
		if str, ok := matType.(string); ok {
			switch str {
			case DEFAULTMATERIAL:
				mb = material.NewBuilder(material.Default())
			case GLASSMATERIAL:
				mb = material.NewBuilder(material.Glass())
			default:
				if cache != nil {
					if cached, ok := cache[str]; ok {
						return cached, nil
					}
				}
				mb = material.NewBuilder(material.Default())
			}
		}
	}
	{
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
		mb = mb.WithPattern(pat)
	}
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
	if cache != nil {
		if cacheName, ok := m.Params["name"]; ok {
			cache[fmt.Sprintf("%s", cacheName)] = mb.Build()
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
	Material  *materialInput
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

func translater(o object, cache map[string]material.Material) (shapes.Shape, error) {
	s, err := newShape(o.Type, o.Params, cache)
	if err != nil {
		return nil, err
	}
	finalTransform := matrix.NewIdentity()
	for _, t := range o.Transform {
		if mat, err := t.toMatrix(); err != nil {
			return nil, err
		} else {
			finalTransform = finalTransform.Multiply(mat)
		}
	}
	s = s.WithTransform(finalTransform)
	if o.Material != nil {
		mat, err := o.Material.toMaterial(cache)
		if err != nil {
			return nil, err
		}
		s = s.WithMaterial(mat)
	}
	return s, nil
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
	materialCache := map[string]material.Material{}
	for _, m := range w.Materials {
		if _, err := m.toMaterial(materialCache); err != nil {
			return nil, Cam{}, err
		}
	}
	for _, o := range w.Objects {
		s, err := translater(o, materialCache)
		if err != nil {
			return nil, Cam{}, err
		}
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
