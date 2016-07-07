package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/gift"
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

// Could be very useful: http://homepages.inf.ed.ac.uk/rbf/HIPR2/canny.htm

var targetColor1 = color.NRGBA{224, 84, 64, 255} // works for second to left
var targetColor2 = color.NRGBA{166, 64, 71, 255} // works for leftmost
var targetColor3 = color.NRGBA{255, 0, 0, 255}

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

func iter(i image.Image, fn func(x int, y int, pixel color.Color)) {
	b := i.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fn(x, y, i.At(x, y))
		}
	}
}

func transform(i image.Image, fn func(int, int, color.Color) color.Color) image.Image {
	b := i.Bounds()
	n := image.NewNRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			n.SetNRGBA(x, y, fn(x, y, i.At(x, y)).(color.NRGBA))
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

func photoshop(i image.Image) image.Image {
	g := gift.New(
		gift.UnsharpMask(1.0, 10.0, 10.0),
		gift.Contrast(50),
		gift.Saturation(50),
		//gift.ColorBalance(-50, 50, 50),
		// gift.Convolution( // emboss
		// 	[]float32{
		// 		-1, -1, 0,
		// 		-1, 1, 1,
		// 		0, 1, 1,
		// 	},
		// 	false, false, false, 0.0,
		// ),
		// gift.Convolution( // edge detection
		// 	[]float32{
		// 		-1, -1, -1,
		// 		-1, 8, -1,
		// 		-1, -1, -1,
		// 	},
		// 	false, false, false, 0.0,
		// ),
		// gift.Sobel(),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewNRGBA(g.Bounds(i.Bounds()))

	g.Draw(dst, i)
	return dst
}

func seperateHue(i image.Image) image.Image {
	// Saturation looks very useful
	// Hue... does not

	return transform(i, func(x int, y int, p color.Color) color.Color {
		pix := p.(color.NRGBA)
		c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}

		h, _, _ := c.Hsv()
		h = h / 360
		return color.NRGBA{R: uint8(255 * h), G: uint8(255 * h), B: uint8(255 * h), A: 255}
	})
}

func accentColorDifference(i image.Image) image.Image {
	return transform(i, func(x int, y int, c color.Color) color.Color {
		distance := colorDistance(c.(color.NRGBA), targetColor1)
		return color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255}
	})
}

func accentColorDiffereenceGreyscale(i image.Image) image.Image {
	return transform(i, func(x int, y int, c color.Color) color.Color {
		h := colorDistance(c.(color.NRGBA), targetColor1)

		if h > 0.40 {
			return color.Gray{0}
		} else {
			return color.Gray{uint8(225 - h*255)}
		}

		// c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}
		// _, h, _ := c.Hsv()
		// n.Set(x, y, color.Gray{uint8(h * 255)})
	})
}
