package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
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

// Save this image inside the assets folder with the given name
func save(img image.Image, name string) {
	f, err := os.Create("./assets/" + name)
	checkError(err)
	defer f.Close()

	err = png.Encode(f, img)
	checkError(err)
}

func open(path string) image.Image {
	f, err := os.Open("./assets/" + path)
	checkError(err)
	defer f.Close()

	img, err := png.Decode(f)
	checkError(err)

	return img
}

func convertCv(i *opencv.IplImage) image.Image {
	return i.ToImage()
}

func convertToColorful(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0}
}

// Tracks the results of a GroupLines operation
// 0 is univisited, 1 is rejected, 2 is line
type TrackingMat struct {
	arr  []*Pix
	w, h int
}

func newTrackingMat(width, height int) *TrackingMat {
	return &TrackingMat{
		arr: make([]*Pix, width*height),
		w:   width,
		h:   height,
	}
}

func (m *TrackingMat) get(x, y int) *Pix {
	return m.arr[y*m.w+x]
}

func (m *TrackingMat) set(x, y int, v *Pix) {
	m.arr[y*m.w+x] = v
}

func (m *TrackingMat) iter(fn func(x int, y int, pixel *Pix)) {
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			fn(x, y, m.get(x, y))
		}
	}
}

func euclideanDistance(a Pix, b Pix) float64 {
	dx := float64(a.x) - float64(b.x)
	dy := float64(a.y) - float64(b.y)

	return math.Sqrt(dx*dx + dy*dy)
}
