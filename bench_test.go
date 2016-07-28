package main

import "testing"

func TestHunt(t *testing.T) {

}

func BenchmarkHunt(b *testing.B) {
	swatch := convertSwatches()
	p := open("lowsett.png")
	mat := convertImage(p)
	lineify(mat, swatch, COLOR_THRESHOLD, LINE_WIDTH)
}
