package main

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

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

func (m *TrackingMat) setPix(p Pix) {
	m.set(p.x, p.y, &p)
}

func (m *TrackingMat) iter(fn func(x int, y int, pixel *Pix)) {
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			fn(x, y, m.get(x, y))
		}
	}
}

func (m *TrackingMat) save(n string) {
	ret := image.NewNRGBA(image.Rect(0, 0, m.w, m.h))

	m.iter(func(x, y int, p *Pix) {
		if p == nil {
			ret.Set(x, y, color.Black)

		} else if DEBUG_DRAW_CHUNKS && p.ptype == PIX_CHUNK {
			ret.Set(x, y, color.White)

		} else if p.ptype == PIX_LINE {
			ret.Set(x, y, p.Color)

		} else {
			r, g, b, _ := p.Color.RGBA()
			ret.Set(x, y, color.RGBA{uint8(float64(r) / 65535.0 * 25), uint8(float64(g) / 65535.0 * 25), uint8(float64(b) / 65535.0 * 25), 255})
		}
	})

	save(ret, n)
}

// Float64 matrix
type PixMat struct {
	arr  []*colorful.Color
	w, h int
}

func NewPixMat(width, height int) *PixMat {
	return &PixMat{
		arr: make([]*colorful.Color, width*height),
		w:   width,
		h:   height,
	}
}

func (m *PixMat) get(x, y int) *colorful.Color {
	return m.arr[y*m.w+x]
}

func (m *PixMat) set(x, y int, v *colorful.Color) {
	m.arr[y*m.w+x] = v
}
