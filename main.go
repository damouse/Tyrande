package main

import (
	"fmt"
	"image"
	"os"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	// Settings
	COLOR_THRESHOLD = 0.3
	LINE_WIDTH      = 1
	SWATCH          []*Pix
	POLL_TIME       = 100 * time.Millisecond

	PARALELIZE   = false // kick off multiple vision goroutines
	NUM_PARALLEL = 2
	CACHE_LUV    = true

	LEFT_SCREEN_DIM = image.Rect(0, 32, 2180, 1380)
	CENTER_OFFSET   = Vector{5, 9} // Where the retircle is wrt the screencap /2

	// Debugging Settings
	DEBUG_DRAW_CHUNKS = false
	DEBUG_SAVE_LINES  = false
	DEBUG_DARKEN      = true

	DEBUG_WINDOW   = false
	DEBUG_RUN_ONCE = false

	DEBUG_BENCH = true
	DEBUG_LOG   = true

	DEBUG_STATIC        = false
	DEBUG_SOURCE_STATIC = "cap.png"

	// Utility Globals
	luvCache     = map[uint32]colorful.Color{}
	luvCacheList = make([]colorful.Color, 16777216)
	linearMutex  = &sync.RWMutex{}

	window      *Window
	imageStatic image.Image

	sumVisions, sumModles, totalCycles float64

	// Main Logic globals
	running   bool
	targeting bool
	target    *Char

	closingChan = make(chan bool, 0)
	visionChan  = make(chan Cycle, 100)
	outputChan  = make(chan Vector, 10)

	Chars    []*Char
	CharLock = &sync.RWMutex{}

	centerVector, targetVector, outputVector Vector
)

// Main loop
func hunt() {
	// Returns true if left alt is pressed, signifying we should track
	altPressed := false //input()

	// Update targeting state
	if targeting != altPressed {
		targeting = altPressed
		// debug("Targeting %v", targeting)
	}

	// Track to the closest char
	if targeting {
		CharLock.RLock()
		target = closestCenter(Chars, centerVector)
		CharLock.RUnlock()

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

func linearStart() {
	running = true

	go startRoutine(input)

	for {
		linearHunt()

		if !running || DEBUG_RUN_ONCE {
			fmt.Println("Linear stopping")
			break
		}
	}

	vis := sumVisions / totalCycles
	fmt.Printf("Cycles: \t%1.0f\nAvg cycle: \t%1.0f ms\n", totalCycles, vis)
}

// Like the hunt method, but linearlized
func linearHunt() {
	start := time.Now()

	// Capture
	var mat *PixMatrix

	if DEBUG_STATIC {
		mat = convertImage(imageStatic)
	} else if targeting {
		mat = convertImage(CaptureLeftNarrow(0.3, 0.3))
	} else {
		mat = convertImage(CaptureLeft())
	}

	// Vision
	lines := lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	cx, cy := mat.center()
	center := Vector{cx, cy}

	// Modeling
	lines = filterLines(lines)

	for _, l := range lines {
		l.process()
	}

	chars := buildChars(lines, center)
	Chars = chars

	if DEBUG_WINDOW {
		go window.show(mat.toImage())
	}

	// Input
	// altPressed := input()
	// altPressed = true

	// Update targeting state
	// if targeting != altPressed {
	// 	targeting = altPressed
	// 	// debug("Targeting %v", targeting)
	// }

	// Track to the closest char
	if targeting && len(chars) != 0 {
		target = chars[0]

		// This is "tracking"
		if target != nil {
			outputVector = target.offset
			moveNow(outputVector)
		}
	}

	// fmt.Printf("Cycle: %s\n", time.Since(start))

	sumVisions += time.Since(start).Seconds() * 1000
	totalCycles += 1
}

func stop() {
	fmt.Println("TYR Stopped")
	running = false
	closingChan <- true
}

func init() {
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
}

func main() {
	// sandbox()

	// start()

	if DEBUG_WINDOW {
		go linearStart()
		window.wait()
	} else {
		linearStart()
	}
}

func sandbox() {
	for i := 0; i < 10; i++ {
		fmt.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
		fmt.Println("hi")
		fmt.Println(i)
		// fmt.Printf("\r\rOn %d/10\non\non", i)
		time.Sleep(100 * time.Millisecond)
	}

	os.Exit(0)
}

/*
TOOD:
	- Vision accuracy improvement
	- Vision performance
	- Merge lines
	- Narrowing: smaller capture arear when Targeting
		- Will this fail on closer targets?
	- Seperate parallelization and implementation methods

Questions
	What effect does mouse sensetivity have on output? Check different sensetivities

	Does move-delta make sense at different ranges?
		Dont do move-delta right now. It is harder to implement and cannot account for vertical movement well.
		In other words: can take player mouse movement into account, but not jumping or falling

	Is there a simple multiplier we can use to calculate output? Does this relate to sensetivity testing?
*/
