package main

import (
	"image"
	"image/color"

	"github.com/disintegration/gift"
	"github.com/lucasb-eyer/go-colorful"
)

// Modeling and detecting on-screen players
type Line struct {
	pixels []*Pix
	id     int
	cX, cY int // center
}

func (l *Line) add(p *Pix) {
	l.pixels = append(l.pixels, p)
	p.line = l
}

func NewLine(id int) *Line {
	return &Line{
		[]*Pix{},
		id,
		0,
		0,
	}
}

type Pix struct {
	color.Color
	x, y int
	line *Line
}

// Here's some more info: http://stackoverflow.com/questions/29156091/opencv-edge-border-detection-based-on-color
// Another silhouette detection: http://stackoverflow.com/questions/13586686/extract-external-contour-or-silhouette-of-image-in-python
// Could be very useful: http://homepages.inf.ed.ac.uk/rbf/HIPR2/canny.htm

func iter(i image.Image, fn func(x int, y int, pixel color.Color)) {
	b := i.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fn(x, y, i.At(x, y))
		}
	}
}

// transform TO rgb
func transformRGB(i image.Image, fn func(int, int, color.Color) color.Color) image.Image {
	b := i.Bounds()
	n := image.NewNRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			n.SetNRGBA(x, y, fn(x, y, i.At(x, y)).(color.NRGBA))
		}
	}

	return n
}

// tranform TO grey
func transformGrey(i image.Image, fn func(int, int, color.Color) color.Color) image.Image {
	b := i.Bounds()
	n := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			n.SetGray(x, y, fn(x, y, i.At(x, y)).(color.Gray))
		}
	}

	return n
}

//
// Operations
//
func colorDistance(a color.Color, b color.Color) float64 {
	return convertToColorful(a).DistanceLuv(convertToColorful(b))
}

func photoshop(i image.Image) image.Image {
	g := gift.New(
		gift.UnsharpMask(12.0, 30.0, 20.0),
		gift.Contrast(30),
		// gift.Hue(45),
		// gift.Gamma(0.1),
		// gift.Saturation(10),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewNRGBA(g.Bounds(i.Bounds()))

	g.Draw(dst, i)
	return dst
}

func localmax(i image.Image) image.Image {
	g := gift.New(
		gift.Maximum(5, true),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewGray(g.Bounds(i.Bounds()))

	g.Draw(dst, i)
	return dst
}

func sobel(i image.Image) image.Image {
	g := gift.New(
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
		gift.Sobel(),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewGray(g.Bounds(i.Bounds()))

	g.Draw(dst, i)
	return dst
}

func seperateHue(i image.Image) image.Image {
	// Saturation looks very useful
	// Hue... does not

	return transformRGB(i, func(x int, y int, p color.Color) color.Color {
		pix := p.(color.NRGBA)
		c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}

		_, h, _ := c.Hsv()
		// h = h
		return color.NRGBA{R: uint8(255 * h), G: uint8(255 * h), B: uint8(255 * h), A: 255}
	})
}

func accentColorDifference(i image.Image) image.Image {
	return transformRGB(i, func(x int, y int, c color.Color) color.Color {
		distance := colorDistance(c, SEPERATION_TARGETCOLOR1)
		return color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255}
	})
}

func accentColorDiffereenceGreyscale(i image.Image, checkAgainst color.NRGBA, threshold float64) image.Image {
	return transformGrey(i, func(x int, y int, c color.Color) color.Color {
		if h := colorDistance(c, checkAgainst); h > threshold {
			// fmt.Println(h)
			return color.Gray{0}
		} else {
			// fmt.Println(h)
			return color.Gray{uint8(255 - h*255)}
		}
	})
}

// This works, but its going to take some time to tweak the knobs
func accentColorDiffereenceGreyscaleAggregate(i image.Image) image.Image {
	return transformGrey(i, func(x int, y int, c color.Color) color.Color {
		rgb := c.(color.NRGBA)

		d1 := colorDistance(rgb, SEPERATION_TARGETCOLOR1)
		d2 := colorDistance(rgb, SEPERATION_TARGETCOLOR2)
		d3 := colorDistance(rgb, SEPERATION_TARGETCOLOR3)

		if h := (d1/3 + d2/3 + d3/3); h > SEPERATION_THRESHOLD {
			return color.Gray{0}
		} else {
			return color.Gray{uint8(255 - h*255)}
		}
	})
}
