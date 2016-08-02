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

	// Progress towards the destination
	px, py := 0.0, 0.0
	// pDist := 0.0

	// target x,y
	tx, ty := -float64(t.x), -float64(t.y)

	// How much to update each cycle
	lx, ly := math.Ceil(tx/OUT_CYCLES), math.Ceil(ty/OUT_CYCLES)

	// The distance we're going to travel
	dist := euclideanDistanceFloat(0.0, 0.0, tx, ty)

	closeEnough := math.Ceil(euclideanDistanceFloat(0.0, 0.0, lx, ly))

	// fmt.Printf("Target: %.0f closeenough %0.f", dist, closeEnough)

	// for i := 0.0; i < OUT_CYCLES; i++ {
	for {
		moveRelative(int(lx), int(ly))

		px += lx
		py += ly

		dist = euclideanDistanceFloat(0.0, 0.0, px, py)

		fmt.Printf("Update: %d, %d Dist: %.0f Target: %0.f\n", int(lx), int(ly), dist, closeEnough)

		if dist <= closeEnough {
			fmt.Printf("Remaining: %0.f, %0.f\n", tx-px, ty-py)
			break
		}

		time.Sleep(OUT_TIME)
	}

	// dist := euclideanDistanceFloat(0.0, 0.0, tx, ty)
	// closeEnough := math.Ceil(euclideanDistanceFloat(0.0, 0.0, lx, ly))

	// fmt.Printf("Dist: %f Close enough: %f delta: %d %d\n", dist, closeEnough, lx, ly)

	// for {
	// 	moveRelative(lx, ly)

	// 	px += lx
	// 	py += ly

	// 	// Update the progress and see how close we are
	// 	dist = euclideanDistanceFloat(px, py, tx, ty)

	// 	// moveRelative(lx, ly)

	// 	if dist <= closeEnough {
	// 		break
	// 	}
	// }

	fmt.Println("Tracking completed")

	tracking = false
}
