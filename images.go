package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/gift"
)

// Modeling and detecting on-screen players
type Line struct {
	pixels []*Pix
	id     int
	cX, cY int // center
}

func (l *Line) add(p *Pix) {
	if p.line == nil {
		l.pixels = append(l.pixels, p)
		p.line = l
	}
}

func (l *Line) addAll(p []*Pix) {
	for _, a := range p {
		l.add(a)
	}
}

func (l *Line) merge(o *Line) {
	for _, p := range o.pixels {
		p.line = nil
		l.add(p)
	}
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
	// gift.Hue(45),
	// gift.Contrast(1),
	// gift.Saturation(2),
	// gift.Gamma(0.75),
	// gift.UnsharpMask(12.0, 30.0, 20.0),
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

func loadSwatch() (result []color.Color) {
	var ret []color.Color

	i := open("swatch.png")

	iter(i, func(x, y int, c color.Color) {
		r, g, b, a := c.RGBA()

		if a == 0 {
			return
		}

		for _, c := range ret {
			er, eg, eb, _ := c.RGBA()

			if er == r && eg == g && eb == b {
				return
			}
		}

		ret = append(ret, c)
	})

	img := image.NewNRGBA(image.Rect(0, 0, 1, len(ret)))

	for i, c := range ret {
		r, g, b, a := c.RGBA()
		img.SetNRGBA(0, i, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	ps := photoshop(img)

	iter(ps, func(x, y int, c color.Color) {
		result = append(result, c)
	})

	fmt.Printf("Loaded %d colors\n", len(ret))
	save(ps, "edittedswatch.png")
	return
}
