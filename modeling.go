package main

import (
	"image"
	"image/color"
)

// Modeling and detecting on-screen players
type Line struct {
	pixels []Pix
}

type Pix struct {
	color.Color
	x, y int
}

// returns all pixels within range d of target coordinates x, y
// The pixel in the middle is included in results
func neighborPixels(tX, tY, distance int, i image.Image) (ret []Pix) {
	d := distance*2 + 1
	b := i.Bounds()

	for x := tX - distance; x <= tX+d; x++ {
		if x < b.Min.X || x > b.Max.X {
			continue
		}

		for y := tY - distance; y <= tY+d; y++ {
			if y < b.Min.Y || y > b.Max.Y {
				continue
			}

			ret = append(ret, Pix{i.At(x, y), x, y})
		}
	}

	return
}

// Return matching pixels from chunk
func scanChunk(chunk []Pix, color color.Color, thresh float64) (ret []Pix) {
	for _, p := range chunk {
		if d := colorDistance(p, color); d < thresh {
			ret = append(ret, p)
		}
	}

	return
}

// Identifies lines in a picture that have a color within thresh distance of a color in col
func huntLines(img image.Image, colors []color.Color, thresh float64, width int) (ret []Line) {
	// We already have information about visited pixels from iter coordinates
	// This stops the algo from re-processing a rejected shape for every pixel
	visited := newTrackingMat(img.Bounds().Max.X+1, img.Bounds().Max.Y+1)

	// For each pixel
	iter(img, func(x, y int, c color.Color) {
		// If this pixel has been visited do not process it
		if visited.get(x, y) != 0 {
			return
		}

		// Measure color distance between this pixel and target colors
		var reference *color.Color

		for _, t := range colors {
			if distance := colorDistance(c, t); distance < thresh {
				reference = &c
				break
			}
		}

		// If no close match found continue
		if reference == nil {
			return
		}

		//
		// -- Hunting
		// Get adjacent pixels
		neighbors := neighborPixels(x, y, width, img)

		//
		// -- Chunking
		matches := scanChunk(neighbors, *reference, 0.1)

		// If the chunk all has the same color begin rejection: chunk and reject whole area
		if len(matches) == (width*2 + 1) {
			// TODO: Grow the chunk, then reject it
			return
		}

		// Note: a line that connects to a chunk should not be rejected

		//
		// -- Tracing
		// Create a line object
		// Add matching pixels
		// Reject bad pixels

		// For all neighbors in group
		// If pixel matches and has not been visited repeat chunking
		// If no more matches found add group to return and resume hunt

		//
		// -- Modeling
		// Merge close lines
		// Determine a center
		// Note: this isnt going into this function, put it somewhere else
	})

	return
}

// Tracks the results of a GroupLines operation
type TrackingMat struct {
	arr  []uint8
	w, h int
}

func newTrackingMat(width, height int) *TrackingMat {
	return &TrackingMat{
		arr: make([]uint8, width*height),
		w:   width,
		h:   height,
	}
}

func (m *TrackingMat) get(x, y int) uint8 {
	return m.arr[x*m.w+y]
}

func (m *TrackingMat) set(x, y int, v uint8) {
	m.arr[x*m.w+y] = v
}
