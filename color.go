package main

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
)

type ptype int

const (
	PIX_NOTHING ptype = iota
	PIX_CHUNK
	PIX_LINE
	PIX_VISITED
)

type Pix struct {
	color.Color

	x, y    int     // coordinates of this pixel
	r, g, b float64 // these are also l, u, v

	lazyInit bool // true if luv has been calculated, else false

	ptype       // A marker that vision may set as needed
	filltag int // used by the fill algo
}

func index(c color.Color) uint32 {
	if cast, ok := c.(color.NRGBA); ok {
		return (uint32(cast.R) << 16) | (uint32(cast.G) << 8) | uint32(cast.B)
	} else if cast, ok := c.(color.RGBA); ok {
		return (uint32(cast.R) << 16) | (uint32(cast.G) << 8) | uint32(cast.B)
	} else {
		panic("Unknown color type!")
	}
}

//
// Pix
func NewPix(x, y int, c color.Color) *Pix {
	return &Pix{
		Color:    c,
		x:        x,
		y:        y,
		r:        0.0,
		g:        0.0,
		b:        0.0,
		lazyInit: false,
		ptype:    PIX_NOTHING,
		filltag:  0,
	}
}

// Trigger the lazy initializer for this pixels luv color
func (p *Pix) initLuv() {
	if p.lazyInit {
		return
	}

	l, u, v := convertToColorful(p.Color).Luv()

	p.r = l
	p.g = u
	p.b = v

	p.lazyInit = true
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
	// a.initLuv()
	// b.initLuv()

	// return math.Sqrt(sq(a.r-b.r) + sq(a.g-b.g) + sq(a.b-b.b))

	if CACHE_LUV {
		c := luvCacheList[index(a.Color)]
		d := luvCacheList[index(b.Color)]

		return math.Sqrt(sq(c.R-d.R) + sq(c.G-d.G) + sq(c.B-d.B))
	} else {
		a.initLuv()
		b.initLuv()

		return math.Sqrt(sq(a.r-b.r) + sq(a.g-b.g) + sq(a.b-b.b))
	}
}

func buildLuvCache() {
	cacher := make([]colorful.Color, 16777216)

	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g++ {
			for b := 0; b < 256; b++ {
				key := uint32((r << 16) | (g << 8) | b)
				cul := colorful.Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
				l, u, v := cul.Luv()
				cacher[key] = colorful.Color{l, u, v}
			}
		}
	}

	dat, err := os.OpenFile("D:\\tyrande.dat", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	checkError(err)

	enc := gob.NewEncoder(dat)

	err = enc.Encode(&cacher)
	checkError(err)

	dat.Close()
	fmt.Println("Done writing")
}

func loadLuvCache() {
	s := time.Now()
	dat, err := os.Open("D:\\tyrande.dat")
	dec := gob.NewDecoder(dat)

	err = dec.Decode(&luvCacheList)
	checkError(err)

	fmt.Printf("LUV Cache loaded in: \t%s\n", time.Since(s))
}
