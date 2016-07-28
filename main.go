package main

import (
	"fmt"
	"image/color"
	"time"
)

var (
	COLOR_THRESHOLD float64 = 0.2
	LINE_WIDTH      int     = 1

	DEBUG_DRAW_CHUNKS = true

	TARGET_SWATCH = []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}
)

func convertSwatches() (ret []*Pix) {
	for _, c := range TARGET_SWATCH {
		ret = append(ret, NewPix(0, 0, c))
	}
	return
}

func runOnce(colors []*Pix) {
	// Load the image
	p := open("lowsett.png")

	// Line detection
	start := time.Now()
	chunks, lines := hunt(p, colors, COLOR_THRESHOLD, LINE_WIDTH)
	fmt.Printf("Hunt completed in: %s\n", time.Since(start))

	// Model detection

	// Update movement logic

	p = output(p.Bounds(), chunks, lines)
	save(p, "huntress.png")
}

func runContinuously(colors []*Pix) {
	w := NewWindow()

	go func(win *Window) {
		for {
			p := open("lowsett.png")
			// Start benchmark
			start := time.Now()

			chunks, lines := hunt(p, colors, COLOR_THRESHOLD, LINE_WIDTH)

			// End benchmark
			fmt.Printf("Hunt completed in: %s\n", time.Since(start))

			p = output(p.Bounds(), chunks, lines)

			win.show(p)
		}
	}(w)

	w.wait()
}

func saveShop() {
	p := open("lowsett.png")
	p = photoshop(p)
	save(p, "1.png")
}

func main() {
	// s := NewSentinal()
	// s.hunt()
	// s.wait()

	// Until the perfromance issues are handled within getLines we cant handle all the swatch colors
	// swatch := loadSwatch()

	swatch := convertSwatches()

	// runContinuously(swatch)
	runOnce(swatch)
}
