package main

import (
	"fmt"
	"image"
	"os"
	"runtime/pprof"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	// Settings
	COLOR_THRESHOLD = 0.25
	LINE_WIDTH      = 1
	SWATCH          []*Pix

	OUT_TIME   = 1 * time.Millisecond // how long output should wait between cycles
	OUT_CYCLES = 20.0                 // how many update cycles to perform
	CACHE_LUV  = true

	LOG_BENCH = true
	LOG       = true

	LEFT_SCREEN_DIM = image.Rect(0, 32, 2180, 1380)
	CENTER_OFFSET   = Vec{5, 9} // Where the retircle is wrt the screencap /2

	// Debugging Settings
	DEBUG_DRAW_CHUNKS = false
	DEBUG_SAVE_LINES  = false
	DEBUG_DARKEN      = true
	DEBUG_CAP_MAT     = true

	DEBUG_WINDOW   = false
	DEBUG_RUN_ONCE = false

	DEBUG_STATIC        = false
	DEBUG_SOURCE_STATIC = "cap.png"

	window      *Window
	imageStatic image.Image

	// Utility Globals
	luvCacheList = make([]colorful.Color, 16777216)

	totalCycleTime, totalCycles float64

	// Main Logic globals
	running, targeting, tracking    bool
	centerVec, targetVec, outputVec Vec

	closingChan = make(chan bool, 0)
	outputChan  = make(chan Vec, 10)

	Chars []*Char
)

func hunt() *image.NRGBA {
	// Capture
	mat := capture()

	// Vision
	lines := lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	cx, cy := mat.center()
	center := Vec{cx, cy}

	// Modeling
	lines = filterLines(lines)

	for _, l := range lines {
		l.process()
	}

	chars := buildChars(lines, center)
	Chars = chars

	// Track to the closest char
	if targeting && len(chars) != 0 {
		moveTo(chars[0].offset)
	}

	if DEBUG_SAVE_LINES || DEBUG_WINDOW {
		return draw(mat, lines, chars)
	} else {
		return nil
	}
}

func start() {
	fmt.Println("TYR Starting")

	var start time.Time
	var i *image.NRGBA
	running = true

	go startRoutine(input)

	for {
		start = time.Now()
		i = hunt()

		if LOG_BENCH {
			fmt.Printf("Cycle: %s\n", time.Since(start))
		}

		totalCycleTime += time.Since(start).Seconds() * 1000
		totalCycles += 1

		if DEBUG_WINDOW {
			go window.show(i)
		}

		if DEBUG_SAVE_LINES {
			save(i, "huntmat.png")
		}

		if !running || DEBUG_RUN_ONCE || DEBUG_SAVE_LINES {
			break
		}
	}

	// Print out benchmarking information
	vis := totalCycleTime / totalCycles
	fmt.Printf("Cycles: \t%1.0f\nAvg cycle: \t%1.0f ms\n", totalCycles, vis)
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

func stop() {
	fmt.Println("TYR Stopped")
	running = false
	closingChan <- true
}

func profile() {
	f, err := os.Create("cpu.out")
	checkError(err)

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	start()
}

func main() {
	// profile()
	// sandbox()

	if DEBUG_WINDOW {
		go start()
		window.wait()
	} else {
		start()
	}
}

func sandbox() {

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
