package canvas

import (
	"fmt"
	"github.com/liorokman/raytrace/pkg/tuple"
	"io"
	"math"
	"text/template"
)

type Canvas interface {
	SetPixel(x, y uint32, c tuple.Color) error
	GetPixel(x, y uint32) (tuple.Color, error)
	Width() uint32
	Height() uint32
}

type canvas struct {
	Data      []tuple.Color
	TheWidth  uint32
	TheHeight uint32
}

func New(width, height uint32) *canvas {
	return &canvas{
		Data:      make([]tuple.Color, height*width),
		TheWidth:  width,
		TheHeight: height,
	}
}

func (c *canvas) Width() uint32 {
	return c.TheWidth
}

func (c *canvas) Height() uint32 {
	return c.TheHeight
}

func (c *canvas) calcPos(x, y uint32) uint32 {
	return y*c.Width() + x
}

func (c *canvas) GetPixel(x, y uint32) (tuple.Color, error) {
	if x >= c.Width() || y >= c.Height() {
		return tuple.NewColor(0, 0, 0), fmt.Errorf("(%v, %v) is outside of the canvas (%v, %v)", x, y, c.Width(), c.Height())
	}
	return c.Data[c.calcPos(x, y)], nil
}

func (c *canvas) SetPixel(x, y uint32, v tuple.Color) error {
	if x >= c.Width() || y >= c.Height() {
		return fmt.Errorf("(%v, %v) is outside of the canvas (%v, %v)", x, y, c.Width(), c.Height())
	}
	c.Data[c.calcPos(x, y)] = v
	return nil
}

func (c *canvas) WritePPM(wr io.Writer) error {
	ppmTemplate := `P3
{{ .Width }} {{ .Height }}
255{{- $pos := 0 }}{{$w := .Width}}
{{ range $index, $element := .Data }}{{- $scaled := $element.Mult 255  }}
{{- ceil $scaled.Red | clamp }} {{ ceil $scaled.Green | clamp }} {{ ceil $scaled.Blue | clamp -}}
{{- $pos = newpos $pos $scaled -}}{{ if or (ge $pos 56) (isnewline $index $w) }}{{ $pos = 0 }}
{{ else }} {{ end }}{{end}}

	`
	funcMap := template.FuncMap{
		"isnewline": func(loc int, width uint32) bool {
			return (loc+1)%int(width) == 0
		},
		"ceil": func(v float64) int { return int(math.Ceil(v)) },
		"clamp": func(v int) int {
			if v > 255 {
				return 255
			} else if v < 0 {
				return 0
			}
			return v
		},
		"newpos": func(pos int, element tuple.Color) int {
			tmp1 := func(f float64) int {
				if f <= 9.0 {
					return 1
				} else if f <= 99.0 {
					return 2
				} else {
					return 3
				}
			}
			x := 2 // two spaces between the three numbers
			x = x + tmp1(element.Red())
			x = x + tmp1(element.Green())
			x = x + tmp1(element.Blue())
			return x + pos
		},
	}

	tmpl, err := template.New("ppm").Funcs(funcMap).Parse(ppmTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, c)
}
