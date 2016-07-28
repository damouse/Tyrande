package main

import "github.com/lucasb-eyer/go-colorful"

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

// Float64 matrix
type FloatMat struct {
	arr  []*colorful.Color
	w, h int
}

func newFloatMat(width, height int) *FloatMat {
	return &FloatMat{
		arr: make([]*colorful.Color, width*height),
		w:   width,
		h:   height,
	}
}

func (m *FloatMat) get(x, y int) *colorful.Color {
	return m.arr[y*m.w+x]
}

func (m *FloatMat) set(x, y int, v *colorful.Color) {
	m.arr[y*m.w+x] = v
}
