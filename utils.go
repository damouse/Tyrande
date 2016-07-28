package main

import (
	"fmt"
	"math"
	"time"
)

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
func euclideanDistance(a Pix, b Pix) float64 {
	dx := float64(a.x) - float64(b.x)
	dy := float64(a.y) - float64(b.y)

	return math.Sqrt(dx*dx + dy*dy)
}

func sq(v float64) float64 {
	return v * v
}
