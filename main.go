package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Global settings
var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = false // draw the rejected color matches on the resulting debug image
	CACHE_LUV         = true  // Cache luv processing (NOTE: this fucks with the colors!) implementation is not correct

	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	CONVERTING_GOROUTINES = 8 // Number of concurrent workers for converting rgb -> LUV

	SWATCH []*Pix // colors to check against
)

//
// Main loop
func hunt(mat *PixMatrix) {
	start := time.Now()

	lines := lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	cX, cY := mat.center()
	closestCenter(lines, cX, cY)

	//

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))
}

//
// Tasks
func staticOnce() {
	p := open("retry.png")
	mat := convertImage(p)

	hunt(mat)
	mat.save("huntress.png")
}

func screencapOnce() {
	p := CaptureLeft()
	mat := convertImage(p)

	hunt(mat)
	mat.save("huntress.png")
}

func continuouslyWindowed() {
	w := NewWindow()

	go func(win *Window) {
		for {
			p := CaptureLeft()
			mat := convertImage(p)

			hunt(mat)

			w.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func main() {
	fmt.Println("Tyrande starting")
	loadSwatch()

	// continuouslyWindowed()
	// screencapOnce()
	// staticOnce()

	windowsAPI()
}
