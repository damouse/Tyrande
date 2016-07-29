package main

import (
	"fmt"
	"time"

	"github.com/AllenDang/w32"
	"github.com/lxn/win"
)

// Input and output manipulation
// This whole API comes from here: https://github.com/lxn/win
// docs: https://godoc.org/github.com/lxn/win

// If we cant get this working directly check out this lib:
// https://play.golang.org/p/kwfYDhhiqk

func windowsAPI() {
	time.Sleep(1000 * time.Millisecond)

	for i := 0; i < 100; i++ {
		p := win.POINT{}
		win.GetCursorPos(&p)
		fmt.Printf("Current position: %v\n", p)

		// Example in C
		// INPUT input;
		// input.type = INPUT_MOUSE;

		// input.mi.mouseData=0;
		// input.mi.dx =  x*(65536/GetSystemMetrics(SM_CXSCREEN));//x being coord in pixels
		// input.mi.dy =  y*(65536/GetSystemMetrics(SM_CYSCREEN));//y being coord in pixels
		// input.mi.dwFlags = MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;

		// SendInput(1,&input,sizeof(input));

		// Attempt 1 in Go
		// input := win.MOUSE_INPUT{}
		// input.Type = win.INPUT_MOUSE

		// input.Mi.MouseData = 0
		// input.Mi.Dx = p.X + int32(i)
		// input.Mi.Dy = p.Y + int32(i)

		// input.Mi.DwFlags = win.MOUSEEVENTF_ABSOLUTE | win.MOUSEEVENTF_MOVE

		// s := []win.MOUSE_INPUT{input}

		// ptr := unsafe.Pointer(&s)
		// sz := int32(unsafe.Sizeof(s))

		// win.SendInput(1, ptr, sz)

		// Attempt 2 in Go
		var inputs []w32.INPUT
		inputs = append(inputs, w32.INPUT{
			Type: w32.INPUT_MOUSE,
			Mi: w32.MOUSEINPUT{
				Dx:          1, //int32
				Dy:          1, //int32
				MouseData:   0, //uint32
				DwFlags:     1, //uint32
				Time:        0, //uint32
				DwExtraInfo: 0, //uintptr
			},
		})

		w32.SendInput(inputs)

		// win.SetCursorPos(p.X+1, p.Y+1)
		time.Sleep(10 * time.Millisecond)
	}

	// p := win.POINT{}
	// win.GetCursorPos(&p)
	// fmt.Printf("Current position: %v\n", p)
}

// Move the cursor a certain distance
// func moveTo(x, y, int) {

// }
