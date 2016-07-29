package main

import (
	"fmt"
	"math"
	"time"
)

// Tasks
func staticOnce() {
	p := open("retry.png")

	start := time.Now()

	mat := convertImage(p)

	lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)
	mat.save("huntress.png")

	fmt.Printf("Hunt completed in: %s\n", time.Since(start))
}

func screencapOnce() {
	p := CaptureLeft()
	mat := convertImage(p)

	lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

	mat.save("huntress.png")
}

func continuouslyWindowed() {
	w := NewWindow()

	go func(win *Window) {
		for {
			start := time.Now()

			p := CaptureLeft()
			mat := convertImage(p)

			lineify(mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)

			fmt.Printf("Hunt completed in: %s\n", time.Since(start))
			w.show(mat.toImage())
		}
	}(w)

	w.wait()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// OLD SAMPLE CODE:
func benchCaptures() {
	start := time.Now()
	iterations := 20

	for i := 0; i < iterations; i++ {
		_, err := CaptureScreen()
		checkError(err)
	}

	fmt.Printf("%d shots took %s", iterations, time.Since(start))
}

// Math utils
// func euclideanDistance(a Pix, b Pix) float64 {
// 	dx := float64(a.x) - float64(b.x)
// 	dy := float64(a.y) - float64(b.y)

// 	return math.Sqrt(dx*dx + dy*dy)
// }

func euclideanDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1) - float64(x2)
	dy := float64(y1) - float64(y2)

	return math.Sqrt(dx*dx + dy*dy)
}

func sq(v float64) float64 {
	return v * v
}

func bench(name string, start time.Time) {
	if DEBUG_BENCH {
		fmt.Printf("%s \t%s\n", name, time.Since(start))
	}
}
