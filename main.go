package main

import (
	"fmt"
	"image/color"
	"runtime"
	"time"
)

var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = false
)

// former manual checks on "lowsett.png"
// allColors := []color.Color{
// 	color.NRGBA{219, 18, 29, 255},
// 	color.NRGBA{140, 31, 59, 255},
// 	color.NRGBA{182, 40, 59, 255},
// 	color.NRGBA{212, 128, 151, 255},
// }

func runpipe() {
	p := NewPipeline()
	p.run(open("0.png"))
	p.save()
}

func runOnce(colors []color.Color) {
	p := open("lowsett.png")

	// Start benchmark
	start := time.Now()

	chunks, lines := hunt(p, colors, COLOR_THRESHOLD, LINE_WIDTH)

	// End benchmark
	fmt.Printf("Bench: %s\n", time.Since(start))

	p = output(p.Bounds(), chunks, lines)
	save(p, "huntress.png")

}

func saveShop() {
	p := open("lowsett.png")
	p = photoshop(p)
	save(p, "1.png")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Until the perfromance issues are handled within getLines we cant handle all the swatch colors
	// swatch := loadSwatch()

	swatch := []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}

	runOnce(swatch)

	// saveShop()
}
