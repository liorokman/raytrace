package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/liorokman/raytrace/pkg/canvas"
	"github.com/liorokman/raytrace/pkg/light"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/ray"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func main() {

	var filename = flag.String("filename", "", "The file to write output to")
	var canvasPixels = flag.Int("width", 100, "Number of pixels on the width of the canvas")
	var wallZ = flag.Float64("wallz", 10.0, "The Z coordinate for the wall")
	var wallSize = flag.Float64("wallsize", 7.0, "The size of the wall")
	var frame = flag.Bool("frame", false, "Frame the canvas")
	flag.Parse()
	if *filename == "" {
		fmt.Printf("Must provide a filename.\n")
		flag.Usage()
		os.Exit(1)
	}

	c := canvas.New(uint32(*canvasPixels), uint32(*canvasPixels))

	//color := tuple.NewColor(1, 0, 0)

	l := light.New(tuple.NewPoint(-10, 10, -10), tuple.NewColor(1, 1, 1))
	eye := tuple.NewPoint(0, 0, -5)

	pixelSize := *wallSize / float64(*canvasPixels)

	halfWall := *wallSize / 2.0

	shape := shapes.NewSphere().WithMaterial(material.New(tuple.NewColor(1, 0.2, 1), 0.1, 0.9, 0.9, 200))

	fmt.Printf("Pixelsize: %#v\n", pixelSize)
	for y := 0; y < *canvasPixels; y++ {
		worldY := halfWall - pixelSize*float64(y)
		for x := 0; x < *canvasPixels; x++ {
			worldX := -halfWall + pixelSize*float64(x)

			position := tuple.NewPoint(float64(worldX), float64(worldY), *wallZ)

			r, err := ray.New(eye, position.Subtract(eye).Normalize())
			if err != nil {
				panic(err)
			}
			xs := r.Intersect(shape)
			if h, ok := ray.Hit(xs...); ok {

				p := r.Position(h.T)
				n := h.Shape.NormalAt(p)
				eyev := r.Direction.Mult(-1)

				c.SetPixel(uint32(x), uint32(y), h.Shape.GetMaterial().Lighting(l, p, eyev, n))
			}

		}
	}

	if *frame {
		borderColor := tuple.NewColor(1, 1, 1)
		for x := 0; x < *canvasPixels; x++ {
			c.SetPixel(uint32(0), uint32(x), borderColor)
			c.SetPixel(uint32(*canvasPixels-1), uint32(x), borderColor)

			c.SetPixel(uint32(x), uint32(0), borderColor)
			c.SetPixel(uint32(x), uint32(*canvasPixels-1), borderColor)
		}
	}

	file, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("Failed to open %s for output: %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	c.WritePPM(file)
}
