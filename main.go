package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Global running settings
var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = false

	TARGET_SWATCH = []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}

	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	CACHE_LUV = false
)

func runOnce(colors []*Pix) {
	// Load the image
	p := open("lowsett.png")

	// Benchmark
	start := time.Now()

	mat := convertImage(p)

	lineify(mat, colors, COLOR_THRESHOLD, LINE_WIDTH)

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))

	// Model detection

	// Update movement logic

	mat.save("huntress.png")
}

func runContinuously(colors []*Pix) {
	w := NewWindow()

	go func(win *Window) {
		for {
			p := open("lowsett.png")

			// Benchmark
			start := time.Now()

			mat := convertImage(p)

			lineify(mat, colors, COLOR_THRESHOLD, LINE_WIDTH)

			fmt.Printf("Hunt completed in: %s\n", time.Since(start))

			win.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func main() {
	// Until the perfromance issues are handled within getLines we cant handle all the swatch colors
	// swatch := loadSwatch()

	swatch := convertSwatches()

	// runContinuously(swatch)
	runOnce(swatch)

	// sandbox()
}

func sandbox() {
	fmt.Println("Hello")

	red := color.RGBA{0, 0, 0, 255}

	r, g, b, _ := red.RGBA()

	// fmt.Printf("%#x %#x %#x\n", r, g+256, b)

	// fmt.Printf("%#x %#x %#x\n", red.R, red.G, red.B)

	// all := r + g + 256 + b + 65536
	// fmt.Printf("%#x\n", all)

	final := (r << 16) | (g << 8) | b

	fmt.Printf("%#x\n", final)
}
