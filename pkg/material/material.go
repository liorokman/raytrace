package material

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/light"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Material struct {
	Color     tuple.Color
	ambient   float64
	diffuse   float64
	specular  float64
	shininess float64
}

func New(c tuple.Color, ambient, diffuse, specular, shininess float64) Material {
	if (ambient < 0 || ambient > 1) ||
		(diffuse < 0 || diffuse > 1) ||
		(specular < 0 || specular > 1) ||
		shininess < 0 {
		panic(fmt.Sprintf("Invalid values for reflection parameters: (%f,%f,%f,%f)", ambient, diffuse, specular, shininess))
	}

	return Material{
		Color:     c,
		ambient:   ambient,
		diffuse:   diffuse,
		specular:  specular,
		shininess: shininess,
	}
}

func Default() Material {
	return New(tuple.NewColor(1, 1, 1), 0.1, 0.9, 0.9, 200.0)
}

func (m Material) Lighting(l light.PointLight, point tuple.Tuple, eyev, normal tuple.Tuple) tuple.Color {
	effectiveColor := m.Color.MultColor(l.Intensity())
	lightV := l.Position().Subtract(point).Normalize()

	ambient := effectiveColor.Mult(m.ambient)
	diffuse := tuple.NewColor(0, 0, 0)
	specular := tuple.NewColor(0, 0, 0)

	lightDotNormal := lightV.Dot(normal)
	if lightDotNormal >= 0 {
		diffuse = effectiveColor.Mult(m.diffuse * lightDotNormal)
		reflectV := lightV.Mult(-1).Reflect(normal)
		reflectDotEye := reflectV.Dot(eyev)
		if reflectDotEye > 0 {
			factor := math.Pow(reflectDotEye, m.shininess)
			specular = l.Intensity().Mult(m.specular * factor)
		}
	}
	return ambient.Add(diffuse).Add(specular)
}
