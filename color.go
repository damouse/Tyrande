package main

import (
	"image"
	"image/color"
	"math"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
)

type ptype int

const (
	PIX_NOTHING ptype = iota
	PIX_CHUNK
	PIX_LINE
)

// Modeling and detecting on-screen players
type Line struct {
	pixels []*Pix
	id     int
	cX, cY int // center
}

type Pix struct {
	color.Color
	x, y    int
	r, g, b float64 // these are also l, u, v
	line    *Line
	ptype   // Algos may mark this this pixel as needed
}

func slide(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return (r << 16) | (g << 8) | b
}

//
// Pix
func NewPix(x, y int, c color.Color) *Pix {
	var l, u, v float64

	if CACHE_LUV {
		i := slide(c)

		if r, ok := luvCache[i]; ok {
			l = r.R
			u = r.G
			v = r.B
		} else {
			l, u, v = convertToColorful(c).Luv()
			luvCache[i] = colorful.Color{l, u, v}
		}
	} else {
		l, u, v = convertToColorful(c).Luv()
	}

	return &Pix{
		c,
		x,
		y,
		l,
		u,
		v,
		nil,
		PIX_NOTHING,
	}
}

//
// Line
func NewLine(id int) *Line {
	return &Line{
		[]*Pix{},
		id,
		0,
		0,
	}
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

//
// Misc color
func convertCv(i *opencv.IplImage) image.Image {
	return i.ToImage()
}

func convertToColorful(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0}
}

func colorDistance(a, b *Pix) float64 {
	return math.Sqrt(sq(a.r-b.r) + sq(a.g-b.g) + sq(a.b-b.b))
}

//
// Manual LUV Lookup
// func Luv(l, u, v float64) Color {
// 	return Xyz(LuvToXyz(l, u, v))
// }

// func LuvToXyz(l, u, v float64) (x, y, z float64) {
// 	// D65 white (see above).
// 	return LuvToXyzWhiteRef(l, u, v, D65)
// }

// func Xyz(x, y, z float64) Color {
// 	return LinearRgb(XyzToLinearRgb(x, y, z))
// 	// return FastLinearRgb(XyzToLinearRgb(x, y, z))
// }

// func LinearRgb(r, g, b float64) Color {
// 	// return FastLinearRgb(r, g, b)
// 	return Color{delinearize(r), delinearize(g), delinearize(b)}
// }
