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

	CONVERTING_GOROUTINES = 8
	CACHE_LUV             = false

	// Debugging Settings
	DEBUG_DRAW_CHUNKS = false
	DEBUG_SAVE_LINES  = true
	DEBUG_WINDOW      = false

	DEBUG_BENCH = true
	DEBUG_LOG   = true

	DEBUG_STATIC        = true
	DEBUG_SOURCE_STATIC = "lowsett.png"

	// Utility Globals
	luvCache    = map[uint32]colorful.Color{}
	linearMutex = &sync.RWMutex{}

	lastPixMat *PixMatrix

	window      *Window
	imageStatic image.Image

	// Main Logic globals
	running   bool
	targeting bool
	target    *Character

	closingChan = make(chan bool, 0)
	visionChan  = make(chan Cycle, 10)
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

	// startRoutineTime(vision)

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

func (op *Cycle) bench() {
	if DEBUG_BENCH {
		fmt.Printf("Cycle: \t\t%s\n\tVIS: \t%s\n  MOD: \t%s\n",
			op.model.Sub(op.start), op.vision.Sub(op.start), op.model.Sub(op.vision))
	}
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
	// sandbox()
}

func sandbox() {

}
