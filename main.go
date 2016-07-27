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
		gift.Contrast(1.5),
		gift.Saturation(100),
		gift.Gamma(0.75),
		// gift.UnsharpMask(12.0, 30.0, 20.0),
	)

	dst := image.NewNRGBA(g.Bounds(p.Bounds()))
	g.Draw(dst, p)

	return dst
}

func runOnce(colors []color.Color) {
	p := open("0.png")

	// Start benchmark
	start := time.Now()

	p = testshop(p)

	chunks, lines := hunt(p, colors, 0.5, 1)

	// End benchmark
	fmt.Printf("Bench: %s\n", time.Since(start))

	p = output(p.Bounds(), chunks, lines)
	save(p, "huntress.png")

}

func saveShop() {
	p := open("0.png")
	p = testshop(p)
	save(p, "1.png")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	swatch := loadSwatch()

	runOnce(swatch)
	// saveShop()
}
