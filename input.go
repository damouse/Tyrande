package main

import (
	"time"

	"github.com/lxn/win"
)

// Input and output manipulation
// This whole API comes from here: https://github.com/lxn/win
// docs: https://godoc.org/github.com/lxn/win
func windowsAPI() {

	for i := 0; i < 100; i++ {
		p := win.POINT{}
		win.GetCursorPos(&p)
		win.SetCursorPos(p.X+1, p.Y+1)
		time.Sleep(10 * time.Millisecond)
	}

	// p := win.POINT{}
	// win.GetCursorPos(&p)
	// fmt.Printf("Current position: %v\n", p)
}
