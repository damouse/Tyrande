package main

import (
	"fmt"
	"math"
	"time"

	"github.com/AllenDang/w32"
)

// Input and output manipulation
// https://godoc.org/github.com/lxn/win
// https://play.golang.org/p/kwfYDhhiqk

// Listens for output commands or makes progress on the last output command
func output() {
	// out := <-outputChan
}

func moveNow(t Vec) {
	moveRelative(-t.x, -t.y)
}

// Move the cursor a certain distance
func moveRelative(x, y int) {
	// The integer values for dx and dy are deltas if MOUSEEVENTF_ABSOLUTE is not set,
	// else its where the mouse ends up
	// debug("Output: %d %d", x, y)

	inputs := []w32.INPUT{w32.INPUT{
		Type: w32.INPUT_MOUSE,
		Mi: w32.MOUSEINPUT{
			Dx:          int32(x),
			Dy:          int32(y),
			MouseData:   0,
			DwFlags:     w32.MOUSEEVENTF_MOVE,
			Time:        0,
			DwExtraInfo: 0,
		},
	}}

	w32.SendInput(inputs)
}

// Moves to the given coordinates
func moveTo(t Vec) {
	if tracking {
		return
	}

	tracking = true

	// Progress towards the destination. Updated with the contents of each cycle
	px, py := 0.0, 0.0

	// Target coordinates
	tx, ty := -float64(t.x), -float64(t.y)

	// The distance we're going to travel
	dist := euclideanDistanceFloat(0.0, 0.0, tx, ty)

	fmt.Printf("Target: %.0f %0.f\n", tx, ty)

	count := 0.0

	for {
		// remaining distance in x and y
		rx, ry := tx-px, ty-py

		// Divide up the remaining space based on the number of cycles we've run
		mult := 1 / (OUT_CYCLES - count)

		// Multimple cycle fraction by remaining distance and get update
		ux, uy := math.Ceil(mult*rx), math.Ceil(mult*ry)

		px, py = px+ux, py+uy

		moveRelative(int(ux), int(uy))

		dist = euclideanDistanceFloat(px, py, tx, ty)

		// fmt.Printf("Mult: %0.3f\tProg: %0.f,%0.f\tRemaining: %0.f,%0.f\tUpdate: %0.f,%0.f\tDist: %0.f\n", mult, px, py, rx, ry, ux, uy, dist)

		if dist == 0 {
			fmt.Println("Soft breat")
			break
		}

		dist = euclideanDistanceFloat(px, py, tx, ty)
		count += 1

		time.Sleep(OUT_TIME)
	}

	// Mostly works
	// for i := 0.0; i <= OUT_CYCLES; i++ {
	// 	// Where do we expect to be this cycle
	// 	ex, ey := i*dx, i*dy
	// 	// ex, ey := tx-px, ty-py

	// 	// Update to x, y
	// 	ux, uy := ex-px, ey-py

	// 	// Update our progress
	// 	px, py = px+float64(int(ux)), py+float64(int(uy))

	// 	dist = euclideanDistanceFloat(px, py, tx, ty)
	// 	fmt.Printf("Prog: %0.f,%0.f\tDist: %0.f\n", px, py, dist)
	// 	moveRelative(int(ux), int(uy))

	// 	time.Sleep(OUT_TIME)
	// }

	fmt.Println("Tracking completed")

	tracking = false
}
