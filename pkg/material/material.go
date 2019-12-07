package material

import (
	"fmt"
	"math"

	"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

type Material struct {
	Pattern         Pattern
	ambient         float64
	diffuse         float64
	specular        float64
	shininess       float64
	reflective      float64
	transparency    float64
	refractiveIndex float64
}

type MaterialBuilder struct {
	m Material
}

func NewDefaultBuilder() *MaterialBuilder {
	return NewBuilder(Default())
}

func NewBuilder(m Material) *MaterialBuilder {
	return &MaterialBuilder{m}
}

func (b *MaterialBuilder) WithRefractiveIndex(ri float64) *MaterialBuilder {
	if ri < 0 {
		panic("Refractive index can't be negative")
	}
	b.m.refractiveIndex = ri
	return b
}

func (b *MaterialBuilder) WithTransparency(t float64) *MaterialBuilder {
	if t < 0 {
		panic("Transparency can't be negative")
	}
	b.m.transparency = t
	return b
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

func (b *MaterialBuilder) WithReflective(a float64) *MaterialBuilder {
	if a < 0 || a > 1 {
		panic("Reflective parameter should be in (0,1) range")
	}
	b.m.reflective = a
	return b
}

func (b *MaterialBuilder) WithShininess(a float64) *MaterialBuilder {
	if a < 0 {
		panic("Shininess parameter must be non-negative")
	}
	b.m.shininess = a
	return b
}

func (b *MaterialBuilder) WithPattern(p Pattern) *MaterialBuilder {
	b.m.Pattern = p
	return b
}

func (b *MaterialBuilder) WithColor(c tuple.Color) *MaterialBuilder {
	return b.WithPattern(NewSolidPattern(c))
}

func (b *MaterialBuilder) Reset() *MaterialBuilder {
	return b.ResetTo(Default())
}

func (b *MaterialBuilder) ResetTo(m Material) *MaterialBuilder {
	b.m = m
	return b
}
func (b *MaterialBuilder) Build() Material {
	return b.m
}

func New(p Pattern, ambient, diffuse, specular, shininess, reflective, transparency, refractiveIndex float64) Material {
	if (ambient < 0 || ambient > 1) ||
		(diffuse < 0 || diffuse > 1) ||
		(specular < 0 || specular > 1) ||
		(reflective < 0 || reflective > 1) ||
		transparency < 0 ||
		refractiveIndex < 0 ||
		shininess < 0 {
		panic(fmt.Sprintf("Invalid values for reflection parameters: (%f,%f,%f,%f)", ambient, diffuse, specular, shininess))
	}

	return Material{
		Pattern:         p,
		ambient:         ambient,
		diffuse:         diffuse,
		specular:        specular,
		shininess:       shininess,
		reflective:      reflective,
		transparency:    transparency,
		refractiveIndex: refractiveIndex,
	}
}

func Default() Material {
	return New(NewSolidPattern(tuple.White), 0.1, 0.9, 0.9, 200.0, 0.0, 0.0, 1.0)
}

func Glass() Material {
	return New(NewSolidPattern(tuple.White), 0.1, 0.1, 0.9, 200.0, 0.0, 1.0, 1.5)
}

func (m Material) Transparency() float64 {
	return m.transparency
}

func (m Material) Reflective() float64 {
	return m.reflective
}

func (m Material) RefractiveIndex() float64 {
	return m.refractiveIndex
}

func (m Material) Lighting(objTransform matrix.Matrix, l fixtures.PointLight, point tuple.Tuple, eyev, normal tuple.Tuple, inShadow bool) tuple.Color {
	effectiveColor := m.Pattern.PatternAtObject(objTransform, point).MultColor(l.Intensity())
	lightV := l.Position().Subtract(point).Normalize()

	ambient := effectiveColor.Mult(m.ambient)
	diffuse := tuple.Black
	specular := tuple.Black

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
	if inShadow {
		return ambient
	}
	return ambient.Add(diffuse).Add(specular)
}
