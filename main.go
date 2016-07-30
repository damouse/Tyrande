package main

import (
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Cycle struct {
	mat   *PixMatrix
	lines []*Line
	chars []*Character

	start  time.Time
	vision time.Time
	model  time.Time
}

var (
	// Settings
	COLOR_THRESHOLD = 0.15
	LINE_WIDTH      = 1
	SWATCH          []*Pix
	POLL_TIME       = 100 * time.Millisecond

	PARALELIZE   = false // kick off multiple vision goroutines
	NUM_PARALLEL = 2
	CACHE_LUV    = true

	// Debugging Settings
	DEBUG_DRAW_CHUNKS = false
	DEBUG_SAVE_LINES  = false
	DEBUG_WINDOW      = false

	DEBUG_BENCH = false
	DEBUG_LOG   = true

	DEBUG_STATIC        = true
	DEBUG_SOURCE_STATIC = "lowsett.png"

	// Utility Globals
	luvCache     = map[uint32]colorful.Color{}
	luvCacheList = make([]colorful.Color, 16777216)
	linearMutex  = &sync.RWMutex{}

	lastPixMat *PixMatrix

	window      *Window
	imageStatic image.Image

	sumVisions, sumModles, totalCycles float64

	// Main Logic globals
	running   bool
	targeting bool
	target    *Character

	closingChan = make(chan bool, 0)
	visionChan  = make(chan Cycle, 100)
	outputChan  = make(chan Vector, 10)

	characters    []*Character
	characterLock = &sync.RWMutex{}

	centerVector, targetVector, outputVector Vector
)

// Main loop
func hunt() {
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

	if PARALELIZE {
		startRoutineTime(vision)
	} else {
		startRoutine(vision)
	}

	startRoutine(modeling)
	// startRoutine(output)
	startRoutine(hunt)

	if DEBUG_WINDOW {
		window.wait()
	}

	<-closingChan

	mod := sumModles / totalCycles
	vis := sumVisions / totalCycles
	avg := (sumVisions + sumModles) / totalCycles

	fmt.Printf("Cycles: \t%1.0f\nAvg Cycle: \t%1.0f ms\nAvg VIS: \t%1.0f ms\nAvg MOD: \t%1.0f ms\n", totalCycles, avg, vis, mod)
}

func stop() {
	fmt.Println("TYR Stopped")
	running = false
	closingChan <- true
}

func main() {
	loadSwatch()

	if CACHE_LUV {
		loadLuvCache()
	}

	if DEBUG_WINDOW {
		window = NewWindow()
	}

	if DEBUG_STATIC {
		imageStatic = open(DEBUG_SOURCE_STATIC)
	}

	start()
	// sandbox()
}

func sandbox() {

}
