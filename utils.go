package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

func draw(mat *PixMatrix, lines []*Line, chars []*Char) *image.NRGBA {
	img := mat.toImage()

	colorPixFromVec(mat, img, getCenter(mat))

	for _, char := range chars {
		colorPixFromVec(mat, img, char.center)
	}

	return img
}

func colorPixFromVec(mat *PixMatrix, img *image.NRGBA, vec Vec) {
	m := &Pix{}
	m.x = vec.x
	m.y = vec.y

	// Draw the center of the image
	for _, p := range mat.adjacent(m, 1) {
		img.Set(p.x, p.y, color.NRGBA{0, 255, 255, 255})
	}
}

func debug(s string, args ...interface{}) {
	if LOG {
		fmt.Printf(s+"\n", args...)
	}
}

func log(s string, args ...interface{}) {
	fmt.Printf(s+"\n", args...)
}

// Tasks
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Math utils
func euclideanDistance(x1, y1, x2, y2 int) float64 {
	return euclideanDistanceFloat(float64(x1), float64(y1), float64(x2), float64(y2))
}

func euclideanDistanceFloat(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2

	return math.Sqrt(dx*dx + dy*dy)
}

func euclideanDistanceVec(a, b Vec) float64 {
	return euclideanDistance(a.x, a.y, b.x, b.y)
}

func sq(v float64) float64 {
	return v * v
}

func bench(name string, start time.Time) {
	if LOG_BENCH {
		fmt.Printf("%s \t%s\n", name, time.Since(start))
	}
}

func startRoutine(fn func()) {
	go func() {
		for {
			fn()

			if !running {
				break
			}
		}
	}()
}

//
// Old start and hunt methods
/*
// Main loop
func oldhunt() {
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
		target = closestCenter(Chars, centerVec)
		CharLock.RUnlock()

		// This is "tracking"
		if target != nil {
			outputVec = target.offset
			moveTo(outputVec)
		}
	}

	// bench("TYR", start)

	if !running {
		return
	}

	time.Sleep(POLL_TIME)
}

func oldstart() {
	fmt.Println("TYR Starting")
	running = true

	if PARALELIZE {
		startRoutineTime(vision)
	} else {
		startRoutine(vision)
	}

	startRoutine(modeling)
	// startRoutine(output)
	startRoutine(oldhunt)

	if DEBUG_WINDOW {
		window.wait()
	}

	<-closingChan

	mod := sumModles / totalCycles
	vis := sumVisions / totalCycles
	avg := (sumVisions + sumModles) / totalCycles

	fmt.Printf("Cycles: \t%1.0f\nAvg Cycle: \t%1.0f ms\nAvg VIS: \t%1.0f ms\nAvg MOD: \t%1.0f ms\n", totalCycles, avg, vis, mod)
}
*/
