package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/liorokman/raytrace/pkg/camera"
	"github.com/liorokman/raytrace/pkg/tuple"
	"github.com/liorokman/raytrace/pkg/world"
)

func main() {

	var filename = flag.String("filename", "", "The file to write output to")
	var scenefile = flag.String("scene", "", "The scene input file")
	var frame = flag.Bool("frame", false, "Frame the canvas")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var webpprof = flag.Bool("webpprof", false, "Launch a web-based pprof interface")

	flag.Parse()
	if *scenefile == "" {
		fmt.Printf("Must provide a scene filename.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *filename == "" {
		fmt.Printf("Must provide an output filename.\n")
		flag.Usage()
		os.Exit(1)
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	} else if *webpprof {
		go func() {
			fmt.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	w, camInput, err := world.NewWorld(*scenefile)
	if err != nil {
		fmt.Printf("Error parsing the scene file: %s", err)
		os.Exit(1)
	}

	fmt.Printf("World is %s\n", w)
	fmt.Printf("Cam input is %#v\n", camInput)

	cam := camera.NewCamera(camInput.Hsize, camInput.Vsize, camInput.FieldOfView).
		WithTransform(camera.ViewTransformation(
			camInput.From.ToPoint(),
			camInput.To.ToPoint(),
			camInput.Up.ToVector()))
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
