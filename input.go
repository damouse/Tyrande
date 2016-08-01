package main

import (
	"time"

	"github.com/AllenDang/w32"
)

// Returns true if the left alt key is pressed
func input() {
	if k := w32.GetAsyncKeyState(w32.VK_F1); k != 0 {
		stop()
	}

	shouldTarget := w32.GetAsyncKeyState(w32.VK_LMENU) != 0

	if targeting != shouldTarget {
		targeting = shouldTarget
		debug("Targeting %v", targeting)
	}

	time.Sleep(100 * time.Millisecond)
}
