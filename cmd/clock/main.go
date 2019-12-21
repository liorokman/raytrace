package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/liorokman/raytrace/pkg/canvas"
	"github.com/liorokman/raytrace/pkg/matrix"
	"github.com/liorokman/raytrace/pkg/tuple"
)

func main() {

	var filename = flag.String("filename", "", "The file to write output to")
	var width = flag.Int("width", 150, "Width of the output")
	var height = flag.Int("height", 150, "Height of the output")
	flag.Parse()
	if *filename == "" {
		fmt.Printf("Must provide a filename.\n")
		flag.Usage()
		os.Exit(1)
	}

	c := canvas.New(uint32(*width), uint32(*height))

	radius := 3.0 * float64(c.Width()) / 8.0

	point := tuple.NewPoint(0, 0, 1)

	color := tuple.NewColor(1, 1, 1)

	for i := 0; i < 12; i++ {
		rotate := matrix.NewRotateY(math.Pi / 6.0 * float64(i))
		currHour := rotate.MultiplyTuple(point).Mult(radius)

		fmt.Printf("Hour %d, currHour: %#v \n", i, currHour)

		c.SetPixel(uint32(currHour.X())+c.Width()/2, uint32(currHour.Z())+c.Height()/2, color)
	}

	file, err := os.Create(*filename)
	if err != nil {
		fmt.Printf("Failed to open %s for output: %s\n", *filename, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	c.WritePPM(file)
}
