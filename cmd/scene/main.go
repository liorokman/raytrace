package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/liorokman/raytrace/pkg/camera"
	//"github.com/liorokman/raytrace/pkg/fixtures"
	"github.com/liorokman/raytrace/pkg/material"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/shapes"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/world"
)

func main() {

	var filename = flag.String("filename", "", "The file to write output to")
	var hsize = flag.Uint("hsize", 100, "Horizontal canvas size")
	var vsize = flag.Uint("vsize", 50, "Vertical canvas size")
	var frame = flag.Bool("frame", false, "Frame the canvas")

	flag.Parse()
	if *filename == "" {
		fmt.Printf("Must provide a filename.\n")
		flag.Usage()
		os.Exit(1)
	}

	w := world.New()
	//	w.Lights = append(w.Lights, fixtures.NewPointLight(tuple.NewPoint(10, 10, -15), tuple.NewColor(1, 0.5, 0.5)))
	//mb := material.NewDefaultBuilder().WithPattern(material.NewCheckerPattern(tuple.Red, tuple.NewColor(0.1, 0.1, 1))).WithSpecular(0).WithReflective(0.8)
	mb := material.NewDefaultBuilder().WithPattern(material.NewSolidPattern(tuple.White)).WithSpecular(0).WithReflective(1).WithDiffuse(1)
	/*
		w.AddShapes(shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewScale(10, 0.01, 10)),
			shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, 0, 5).RotateY(-math.Pi/4).RotateX(math.Pi/2).Scale(10, 0.01, 10)),
			shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(0, 0, 5).RotateY(math.Pi/4).RotateX(math.Pi/2).Scale(10, 0.01, 10)))
	*/
	w.AddShapes(shapes.NewPlane().WithTransform(matrix.NewTranslation(0, 0, -10)).WithMaterial(mb.Build()))

	mb.Reset().WithPattern(material.NewRingPattern(tuple.Green, tuple.Blue).WithTransform(matrix.NewTranslation(1.5, 1.5, 1.5).Scale(0.2, 0.2, 0.2))).WithDiffuse(0.7).WithSpecular(0.3)
	w.AddShapes(shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(-1.5, 1, 0.5)))

	mb.WithPattern(material.NewGradientPattern(tuple.NewColor(0.5, 1, 0.1), tuple.NewColor(0.15, 0, 0.9)).WithTransform(matrix.NewTranslation(1.5, 1.5, 1.5)))
	w.AddShapes(shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(1.5, 0.5, -0.5).Scale(0.5, 0.5, 0.5)))
	mb.WithColor(tuple.NewColor(1, 0.8, 0.1))
	w.AddShapes(shapes.NewSphere().WithMaterial(mb.Build()).WithTransform(matrix.NewTranslation(-1.5, 0.33, -0.75).Scale(0.33, 0.33, 0.33)))

	cam := camera.NewCamera(uint32(*hsize), uint32(*vsize), math.Pi/3.0).
		WithTransform(camera.ViewTransformation(
			tuple.NewPoint(-5, 1.5, -5),
			tuple.NewPoint(0, 1, 0),
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
		fmt.Printf("Failed to open %s for output: %s\n", filename, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	image.WritePPM(file)
}
