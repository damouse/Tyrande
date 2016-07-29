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
	if input() {
		characterLock.RLock()
		target = closestCenter(characters, centerVector)
		characterLock.RUnlock()

		// If target vector is not set we're tracking but not targeting.
		// In the future the output vector will only be fired if we're aligning

		// Fire off the output vector
		outputVector = target.offset
		// outputChan <- outputVector
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
