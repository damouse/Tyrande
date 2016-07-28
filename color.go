package main

import (
	"image"
	"image/color"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
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
	r, g, b float64
	line    *Line
}

func NewPix(x, y int, c color.Color) *Pix {
	return &Pix{
		c,
		x,
		y,
		0,
		0,
		0,
		nil,
	}
}

// Convert an image to a Pix matrix
func convertImage(i image.Image) *TrackingMat {
	b := i.Bounds()
	mat := newTrackingMat(b.Max.Y, b.Max.X)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			mat.set(x, y, NewPix(x, y, i.At(x, y)))
		}
	}

	return mat
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
