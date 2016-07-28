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

func NewPix(x, y int, c color.Color) *Pix {
	l, u, v := convertToColorful(c).Luv()

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
