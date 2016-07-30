package main

import (
	"image"
	"image/color"
)

// Tracks the results of a GroupLines operation
// 0 is univisited, 1 is rejected, 2 is line
type PixMatrix struct {
	arr  []Pix
	w, h int
}

func NewPixMatrix(width, height int) *PixMatrix {
	return &PixMatrix{
		arr: make([]Pix, width*height),
		w:   width,
		h:   height,
	}
}

func (m *PixMatrix) get(x, y int) *Pix {
	return &m.arr[y*m.w+x]
}

func (m *PixMatrix) set(v *Pix) {
	m.arr[v.y*m.w+v.x] = *v
}

func (m *PixMatrix) center() (int, int) {
	return m.w / 2, m.h / 2
}

// Get the adjacent pixels to the given pixel that are within distance in x and y
// The given pixel is not returned
func (m *PixMatrix) adjacent(p *Pix, distance int) (ret []*Pix) {
	for x := p.x - distance; x <= p.x+distance; x++ {
		if x < 0 || x >= m.w {
			continue
		}

		for y := p.y - distance; y <= p.y+distance; y++ {
			if y < 0 || y >= m.h {
				continue
			}

			// Dont return this pixel
			if !(x == p.x && y == p.y) {
				ret = append(ret, m.get(x, y))
			}
		}
	}

	return
}

// Returns adjacent pixels that are within thresh color of the given pixel
func (m *PixMatrix) adjacentSimilarColor(p *Pix, target *Pix, distance int, thresh float64) (ret []*Pix) {
	adj := m.adjacent(p, distance)

	for _, n := range adj {
		if colorDistance(target, n) < thresh {
			ret = append(ret, n)
		}
	}

	return
}

//
// File utils
func (m *PixMatrix) save(n string) {
	save(m.toImage(), n)
}

func (m *PixMatrix) toImage() *image.NRGBA {
	ret := image.NewNRGBA(image.Rect(0, 0, m.w, m.h))

	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			p := m.get(x, y)

			if p == nil {
				ret.Set(x, y, color.Black)

			} else if DEBUG_DRAW_CHUNKS && p.ptype == PIX_CHUNK {
				ret.Set(x, y, color.White)

			} else if p.ptype == PIX_LINE {
				ret.Set(x, y, p.Color)

			} else {
				if DEBUG_DARKEN {
					r, g, b, _ := p.Color.RGBA()
					ret.Set(x, y, color.RGBA{uint8(float64(r) / 65535.0 * 25), uint8(float64(g) / 65535.0 * 25), uint8(float64(b) / 65535.0 * 25), 255})
				} else {
					ret.Set(x, y, p.Color)
				}
			}
		}
	}

	return ret
}
