package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
)

// I have serious doubts about the above working.
// Here's some more info: http://stackoverflow.com/questions/29156091/opencv-edge-border-detection-based-on-color

// Pure go image filtering library: https://github.com/disintegration/gift

// Another silhouette detection: http://stackoverflow.com/questions/13586686/extract-external-contour-or-silhouette-of-image-in-python

// The border seems to be a burnt orange-ish color
// It looks like its always about one pixel wide
// Pure: (233, 88, 61)

var targetColor1 = color.NRGBA{224, 84, 64, 255} // works for second to left
var targetColor2 = color.NRGBA{166, 64, 71, 255} // works for leftmost

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

func convertCv(i *opencv.IplImage) Image {
	img := i.ToImage()
	return Image{img.(*image.NRGBA)}
}

func (i *Image) iter(fn func(x int, y int, pixel color.NRGBA)) {
	b := i.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fn(x, y, i.NRGBAAt(x, y))
		}
	}
}

func (i *Image) transform(fn func(int, int, color.NRGBA) color.NRGBA) Image {
	b := i.Bounds()
	n := Image{image.NewNRGBA(b)}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			n.SetNRGBA(x, y, fn(x, y, i.NRGBAAt(x, y)))
		}
	}

	return n
}

//
// Operations
//
func colorDistance(a color.NRGBA, c color.NRGBA) float64 {
	c1 := colorful.Color{float64(a.R) / 255.0, float64(a.G) / 255.0, float64(a.B) / 255.0}
	c2 := colorful.Color{float64(c.R) / 255.0, float64(c.G) / 255.0, float64(c.B) / 255.0}

	// Luv seems quite good
	return c1.DistanceCIE76(c2)
}

func seperateHue(i Image) Image {
	// Saturation looks very useful
	// Hue... does not

	return i.transform(func(x int, y int, pix color.NRGBA) color.NRGBA {
		c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}

		h, _, _ := c.Hsv()
		h = h / 360
		return color.NRGBA{R: uint8(255 * h), G: uint8(255 * h), B: uint8(255 * h), A: 255}
	})
}

func accentColorDifference(i Image) Image {
	n := Image{image.NewNRGBA(i.Bounds())}
	b := i.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {

			pix := i.NRGBAAt(x, y)
			distance := colorDistance(pix, targetColor1)
			newColor := color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255}
			n.SetNRGBA(x, y, newColor)
		}
	}

	return n
}
