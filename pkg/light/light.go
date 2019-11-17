package light

import (
	"github.com/liorokman/raytrace/pkg/tuple"
)

type PointLight struct {
	position  tuple.Tuple
	intensity tuple.Color
}

func New(position tuple.Tuple, intensity tuple.Color) PointLight {
	if !position.IsPoint() {
		panic("Pointlight can't be located at a vector")
	}
	return PointLight{
		position:  position,
		intensity: intensity,
	}
}

func (p PointLight) Position() tuple.Tuple {
	return p.position
}

func (p PointLight) Intensity() tuple.Color {
	return p.intensity
}
