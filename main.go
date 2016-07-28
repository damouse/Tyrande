package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Global running settings
var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = false // draw the rejected color matches on the resulting debug image
	CACHE_LUV         = true  // Cache luv processing

	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	CONVERTING_GOROUTINES = 8 // Number of concurrent workers for converting rgb -> LUV
)

func runOnce(colors []*Pix) {
	// Load the image
	// p := open("lowsett.png")

	// Benchmark
	start := time.Now()

	// Grab the screen
	p := CaptureLeft()
	save(p, "cap.png")

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
			// p := open("lowsett.png")

			// Benchmark
			start := time.Now()
			p := CaptureLeft()

			mat := convertImage(p)

			lineify(mat, colors, COLOR_THRESHOLD, LINE_WIDTH)

			fmt.Printf("Hunt completed in: %s\n", time.Since(start))

			w.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func main() {
	fmt.Println("Tyrande starting")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Until the perfromance issues are handled within getLines we cant handle all the swatch colors
	// swatch := loadSwatch()

	swatch := convertSwatches()

	runContinuously(swatch)
	// runOnce(swatch)

	// sandbox()
}

func sandbox() {
	fmt.Println("Hello")

	a := 3

	fmt.Println(a / 4)
}
