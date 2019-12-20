package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/liorokman/raytrace/pkg/camera"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/world"
)

func hexagonCorner() shapes.Shape {
	mb := material.NewBuilder(material.Glass()).WithColor(tuple.Red).WithAmbient(1.0)
	return shapes.NewSphere().WithTransform(matrix.NewTranslation(0, 0, -1).Scale(0.25, 0.25, 0.25)).WithMaterial(mb.Build())
}

func hexagonEdge() shapes.Shape {
	mb := material.NewDefaultBuilder().WithColor(tuple.Green).WithAmbient(1.0)
	return shapes.NewConstrainedCylinder(0, 1, true).WithTransform(matrix.NewTranslation(0, 0, -1).RotateY(-math.Pi/6).RotateZ(-math.Pi/2.0).Scale(0.25, 1, 0.25)).WithMaterial(mb.Build())
}

func hexagonSide() shapes.Shape {
	side := shapes.NewGroup()

	shapes.Connect(side, hexagonCorner())
	shapes.Connect(side, hexagonEdge())

	return side
}

func hexagon() shapes.Shape {
	hex := shapes.NewGroup()

	for i := 0; i < 6; i++ {
		side := hexagonSide().WithTransform(matrix.NewRotateY((math.Pi / 3.0) * float64(i)))
		shapes.Connect(hex, side)
	}

	return hex
}

func main() {

	var filename = flag.String("filename", "", "The file to write output to")
	var frame = flag.Bool("frame", false, "Frame the canvas")

	flag.Parse()
	if *filename == "" {
		fmt.Printf("Must provide an output filename.\n")
		flag.Usage()
		os.Exit(1)
	}

	w := world.New()
	w.AddShapes(shapes.NewPlane().WithTransform(matrix.NewTranslation(0, -1, 0).RotateX(math.Pi / 6)))
	w.AddShapes(hexagon().WithTransform(matrix.NewScale(1.5, 1.5, 1.5)))

	cam := camera.NewCamera(400, 400, math.Pi/3.0).
		WithTransform(camera.ViewTransformation(
			tuple.NewPoint(0, 5, 5),
			tuple.NewPoint(0, 0, 0),
			tuple.NewVector(0, 1, 0)))
	fmt.Printf("Pixelsize: %v\n", cam.PixelSize())
	image := cam.Render(w)

	if *frame {
		borderColor := tuple.Red
		for x := uint32(0); x < cam.HSize(); x++ {
			for y := uint32(0); y < cam.VSize(); y++ {
				image.SetPixel(uint32(0), uint32(x), borderColor)
				image.SetPixel(uint32(cam.VSize()-1), uint32(x), borderColor)

				image.SetPixel(uint32(x), uint32(0), borderColor)
				image.SetPixel(uint32(x), uint32(cam.HSize()-1), borderColor)
			}
		}
	}

	file, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("Failed to open %s for output: %s\n", *filename, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	err = image.WritePPM(file)
	if err != nil {
		fmt.Printf("Failed to generate the output file: %s\n", err)
		os.Exit(1)
	}
}
