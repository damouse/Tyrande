package main

import (
	"fmt"
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
