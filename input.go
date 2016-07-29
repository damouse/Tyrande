package main

import "github.com/AllenDang/w32"

// Returns true if the left alt key is pressed
func input() bool {
	if k := w32.GetAsyncKeyState(w32.VK_F1); k != 0 {
		stop()
	}

	if k := w32.GetAsyncKeyState(w32.VK_LMENU); k != 0 {
		return true
	}

	return false
}
