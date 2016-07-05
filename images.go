package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Image struct {
	*image.NRGBA
}

// Save this image inside the assets folder with the given name
func (i *Image) save(name string) {
	f, err := os.Create("./assets/" + name)
	checkError(err)
	defer f.Close()

	err = png.Encode(f, i)
	checkError(err)
}

func open(path string) Image {
	f, err := os.Open("./assets/sample.png")
	checkError(err)
	defer f.Close()

	img, err := png.Decode(f)
	checkError(err)

	return Image{img.(*image.NRGBA)}
}

// Color manipulation. Returns the "distance" between two colors
func colorDistance(a color.NRGBA, b color.NRGBA) float64 {
	// r := math.Abs(float64(a.R - b.R))
	// g := math.Abs(float64(a.G - b.G))
	// e := math.Abs(float64(a.B - b.B))

	d := math.Sqrt(float64((a.R - b.R) ^ 2 + (a.G - b.G) ^ 2 + (a.B - b.B) ^ 2))

	return d / math.Sqrt((255)^2+(255)^2+(255)^2)
}

// I have serious doubts about the above working.
// Here's some more info: http://stackoverflow.com/questions/29156091/opencv-edge-border-detection-based-on-color

// Pure go image filtering library: https://github.com/disintegration/gift

// Another silhouette detection: http://stackoverflow.com/questions/13586686/extract-external-contour-or-silhouette-of-image-in-python

// Or just grab opencv: https://github.com/lazywei/go-opencv#disclaimer
