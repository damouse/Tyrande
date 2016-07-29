package main

import (
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Settings
var (
	COLOR_THRESHOLD = 0.2
	LINE_WIDTH      = 1
	SWATCH          []*Pix // colors to check against

	CONVERTING_GOROUTINES = 8    // Number of concurrent workers for converting rgb -> LUV
	CACHE_LUV             = true // Cache luv processing (NOTE: this fucks with the colors!) implementation is not correct

	DEBUG_DRAW_CHUNKS   = false
	DEBUG_SAVE_LINES    = false
	DEBUG_STATIC        = true  // if true sources image from DEBUG_SOURCE below instead of screen
	DEBUG_WINDOW        = false // display a window of the running capture
	DEBUG_SOURCE_STATIC = "lowsett.png"
)

// Utility Globals
var (
	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	window      *Window
	imageStatic image.Image
)

// Main Logic globals
var (
	running      bool
	target       *Character
	targetVector Vector

	visionChan = make(chan []*Line, 10)

	characters    []*Character
	characterLock = &sync.RWMutex{}

	outputVector Vector
)

// Main loop
func hunte() {
	for {
		// Check for input updates

		// Check for close

		// Targeting logic

		// Update output if needed
	}
}

// Kick off the four processing goroutines
func start() {
	fmt.Println("TYR Starting")

	running = true

	vision()
	// go modeling()
	// go output()
	// hunte()
}

func hunt(mat *PixMatrix) {
	lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	// cX, cY := mat.center()
	// closest := closestCenter(lines, cX, cY)

	// moveTo(closest.centerX, closest.centerY)
}

// Tasks
func staticOnce() {
	p := open("retry.png")

	start := time.Now()

	mat := convertImage(p)

	hunt(mat)
	mat.save("huntress.png")

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))
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
			start := time.Now()

			p := CaptureLeft()
			mat := convertImage(p)

			hunt(mat)

			fmt.Printf("Hunt completed in: %s\n", time.Since(start))
			w.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func main() {
	loadSwatch()

	if DEBUG_WINDOW {
		window = NewWindow()
	}

	if DEBUG_STATIC {
		imageStatic = open(DEBUG_SOURCE_STATIC)
	}

	start()

	// window.wait()

	// continuouslyWindowed()
	// screencapOnce()
	// staticOnce()
}
