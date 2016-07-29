package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Global running settings
var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = false // draw the rejected color matches on the resulting debug image
	CACHE_LUV         = true  // Cache luv processing (NOTE: this fucks with the colors!) implementation is not correct

	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	CONVERTING_GOROUTINES = 8 // Number of concurrent workers for converting rgb -> LUV

	SWATCH []*Pix
)

func runStaticOnce() {
	p := open("retry.png")

	start := time.Now()
	mat := convertImage(p)
	lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))

	mat.save("huntress.png")
}

func runScreencapOnce() {
	start := time.Now()

	p := CaptureLeft()
	mat := convertImage(p)
	lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))
	mat.save("huntress.png")
}

func runContinuously() {
	w := NewWindow()

	go func(win *Window) {
		for {
			// p := open("lowsett.png")

			// Benchmark
			start := time.Now()
			p := CaptureLeft()
			mat := convertImage(p)

			lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

			fmt.Printf("Hunt completed in: %s\n", time.Since(start))

			w.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func main() {
	fmt.Println("Tyrande starting")
	loadSwatch()

	runContinuously()
	// runScreencapOnce()
	// runStaticOnce()

	// sandbox()
}

type Alpha struct {
	a int
}

func sandbox() {
	fmt.Println("Hello")

	slicer := make([]Alpha, 3)
	fmt.Println(slicer[2])
}
