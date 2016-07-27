package main

import (
	"image/color"
	"testing"
)

func TestHunt(t *testing.T) {

}

func BenchmarkHunt(b *testing.B) {
	swatch := []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}

	p := open("lowsett.png")

	hunt(p, swatch, COLOR_THRESHOLD, LINE_WIDTH)
}
