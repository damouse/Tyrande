package main

import "testing"

func TestHunt(t *testing.T) {

}

func BenchmarkHunt(b *testing.B) {
	swatch := convertSwatches()
	mat := CaptureLeft()
	lineify(mat, swatch, COLOR_THRESHOLD, LINE_WIDTH)
}
