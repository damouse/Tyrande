package main

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"time"

	"github.com/disintegration/gift"
)

func runpipe() {
	p := NewPipeline()
	p.run(open("0.png"))
	p.save()
}

func testshop(p image.Image) image.Image {
	g := gift.New(
		gift.Contrast(1),
		gift.Gamma(.75),
		gift.UnsharpMask(12.0, 30.0, 20.0),
	)

	dst := image.NewNRGBA(g.Bounds(p.Bounds()))
	g.Draw(dst, p)

	// save(dst, "1.png")

	return dst
}

func runOnce() {
	// for "lowsett.png"
	allColors := []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}

	// for 0.png
	// allColors := []color.Color{
	// 	color.NRGBA{244, 88, 54, 255},
	// 	color.NRGBA{177, 38, 48, 255},
	// }

	p := open("lowsett.png")

	// Start benchmark
	start := time.Now()

	// p = testshop(p)

	chunks, lines := hunt(p, allColors, 0.2, 1)

	// End benchmark
	fmt.Printf("Bench: %s\n", time.Since(start))

	p = output(p.Bounds(), chunks, lines)
	save(p, "huntress.png")

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runOnce()
}
