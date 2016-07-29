package main

import (
	"time"

	"github.com/AllenDang/w32"
	"github.com/lxn/win"
)

// Input and output manipulation
// This whole API comes from here: https://github.com/lxn/win
// docs: https://godoc.org/github.com/lxn/win

// If we cant get this working directly check out this lib:
// https://play.golang.org/p/kwfYDhhiqk

func output() {
	if !running {
		return
	}
}

type Vector struct {
	x, y int
}

// Moves to the given coordinates
func moveTo(x, y int) {
	// time.Sleep(1000 * time.Millisecond)

	p := win.POINT{}
	win.GetCursorPos(&p)

	dx := x - int(p.X)
	dy := y - int(p.Y)

	// milliseconds?
	// duration := 1000
	cycles := 100

	lx := float64(dx) / float64(cycles)
	ly := float64(dy) / float64(cycles)

	for i := 0; i < cycles; i++ {
		moveRelative(int(lx), int(ly))
		time.Sleep(10 * time.Millisecond)
	}
}

// Move the cursor a certain distance
func moveRelative(x, y int) {
	// The integer values for dx and dy are deltas if MOUSEEVENTF_ABSOLUTE is not set,
	// else its where the mouse ends up

	var inputs []w32.INPUT

	inputs = append(inputs, w32.INPUT{
		Type: w32.INPUT_MOUSE,
		Mi: w32.MOUSEINPUT{
			Dx:          int32(x),
			Dy:          int32(y),
			MouseData:   0,
			DwFlags:     w32.MOUSEEVENTF_MOVE,
			Time:        0,
			DwExtraInfo: 0,
		},
	})

	w32.SendInput(inputs)
}

// func windowsAPI() {
// 	time.Sleep(1000 * time.Millisecond)

// 	for i := 0; i < 100; i++ {
// 		p := win.POINT{}
// 		win.GetCursorPos(&p)
// 		// fmt.Printf("Current position: %v\n", p)

// 		// Attempt 2 in Go

// 		// win.SetCursorPos(p.X+1, p.Y+1)
// 		time.Sleep(10 * time.Millisecond)
// 	}

// 	// p := win.POINT{}
// 	// win.GetCursorPos(&p)
// 	// fmt.Printf("Current position: %v\n", p)
// }

/*
Tyrande
	Main event loop

Vision
	Outline input loop. Reads the screen, calls Tyrande with Lines

Model
	Tracks characters in an abstract way

Input
	Handles input-over-time to system

Output
	Watches input, subtracts our input, determines what user is doing
*/
