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
	SWATCH          []*Pix
	POLL_TIME       = 100 * time.Millisecond

	CONVERTING_GOROUTINES = 8    // Number of concurrent workers for converting rgb -> LUV
	CACHE_LUV             = true // Cache luv processing implementation not correct

	DEBUG_DRAW_CHUNKS   = false
	DEBUG_SAVE_LINES    = false
	DEBUG_STATIC        = true  // if true sources image from DEBUG_SOURCE below instead of screen
	DEBUG_WINDOW        = false // display a window of the running capture
	DEBUG_BENCH         = true
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
func hunt() {
	for {
		if !running {
			return
		}

		// start := time.Now()

		// Returns true if left alt is pressed, signifying we should track
		if input() {

		}

		time.Sleep(POLL_TIME)
		// bench("TYR", start)
	}
}

// Kick off the four processing goroutines
func start() {
	fmt.Println("TYR Starting")

	running = true

	go vision()
	go modeling()
	go output()
	hunt()
}

func stop() {
	fmt.Println("TYR Stopped")
	running = false
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

	// continuouslyWindowed()
	// screencapOnce()
	// staticOnce()
}
