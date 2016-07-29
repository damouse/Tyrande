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
	COLOR_THRESHOLD = 0.15
	LINE_WIDTH      = 1
	SWATCH          []*Pix
	POLL_TIME       = 100 * time.Millisecond

	CONVERTING_GOROUTINES = 8    // Number of concurrent workers for converting rgb -> LUV
	CACHE_LUV             = true // Cache luv processing implementation not correct

	DEBUG_DRAW_CHUNKS   = false
	DEBUG_SAVE_LINES    = false
	DEBUG_STATIC        = false // if true sources image from DEBUG_SOURCE below instead of screen
	DEBUG_WINDOW        = true  // display a window of the running capture
	DEBUG_BENCH         = false
	DEBUG_LOG           = true
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
	running   bool
	targeting bool
	target    *Character

	closingChan = make(chan bool, 0)
	visionChan  = make(chan []*Line, 10)
	outputChan  = make(chan Vector, 10)

	characters    []*Character
	characterLock = &sync.RWMutex{}

	centerVector Vector
	targetVector Vector
	outputVector Vector
)

// Main loop
func hunt() {
	// start := time.Now()

	// Returns true if left alt is pressed, signifying we should track
	altPressed := input()

	// Update targeting state
	if targeting != altPressed {
		targeting = altPressed
		// debug("Targeting %v", targeting)
	}

	// Track to the closest char
	if targeting {
		characterLock.RLock()
		target = closestCenter(characters, centerVector)
		characterLock.RUnlock()

		// This is "tracking"
		if target != nil {
			outputVector = target.offset
			moveTo(outputVector)
		}
	}

	// bench("TYR", start)

	if !running {
		return
	}

	time.Sleep(POLL_TIME)
}

func start() {
	fmt.Println("TYR Starting")
	running = true

	startRoutine(vision)
	startRoutine(modeling)
	// startRoutine(output)
	startRoutine(hunt)

	if DEBUG_WINDOW {
		window.wait()
	}

	<-closingChan
}

func stop() {
	fmt.Println("TYR Stopped")
	running = false
	closingChan <- true
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
