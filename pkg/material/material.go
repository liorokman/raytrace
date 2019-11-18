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

type MaterialBuilder struct {
	m Material
}

func NewBuilder() *MaterialBuilder {
	return &MaterialBuilder{Default()}
}

func (b *MaterialBuilder) WithAmbient(a float64) *MaterialBuilder {
	if a < 0 || a > 1 {
		panic("Ambient parameter should be in (0,1) range")
	}
	b.m.ambient = a
	return b
}
func (b *MaterialBuilder) WithDiffuse(a float64) *MaterialBuilder {
	if a < 0 || a > 1 {
		panic("Diffuse parameter should be in (0,1) range")
	}
	b.m.diffuse = a
	return b
}
func (b *MaterialBuilder) WithSpecular(a float64) *MaterialBuilder {
	if a < 0 || a > 1 {
		panic("Specular parameter should be in (0,1) range")
	}
	b.m.specular = a
	return b
}
func (b *MaterialBuilder) WithShininess(a float64) *MaterialBuilder {
	if a < 0 {
		panic("Shininess parameter must be non-negative")
	}
	b.m.shininess = a
	return b
}
func (b *MaterialBuilder) WithColor(c tuple.Color) *MaterialBuilder {
	b.m.Color = c
	return b
}
func (b *MaterialBuilder) Reset() *MaterialBuilder {
	b.m = Default()
	return b
}
func (b *MaterialBuilder) Build() Material {
	return b.m
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
