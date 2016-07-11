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
	b := i.Bounds()

	for x := tX - distance; x <= tX+distance; x++ {
		if x < b.Min.X || x > b.Max.X {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < b.Min.Y || y > b.Max.Y {
				continue
			}

			ret = append(ret, Pix{i.At(x, y), x, y})
		}
	}

	return
}

// Return matching pixels from chunk that are within thresh of the given color
func scanChunk(chunk []Pix, color color.Color, thresh float64) (ret []Pix) {
	for _, p := range chunk {
		if d := colorDistance(p, color); d < thresh {
			ret = append(ret, p)
		}
	}

	return
}

// Check for a match in one of the colors. Return true if its within thresh of a color in targets
func similarColor(c color.Color, targets []color.Color, thresh float64) bool {
	for _, t := range targets {
		if distance := colorDistance(c, t); distance < thresh {
			return true
		}
	}

	return false
}

// Reject the given pixels and return adjacent chunks. Ignore adjacent chunks that
// have already been rejected.
// Return adjacent chunks
func reject(chunk []Pix, mat *TrackingMat, width int, img image.Image, thresh float64, c color.Color) {
	// Reject these pixels
	for _, p := range chunk {
		mat.set(p.x, p.y, 1)
	}

	// Iterate over adjacent pixels
	for _, p := range chunk {
		neighbors := neighborPixels(p.x, p.y, width, img)

		for _, n := range neighbors {

			// ignore already visited pixels
			if mat.get(n.x, n.y) != 0 {
				continue
			}

			// apply recursively to neighbors
			matches := scanChunk(neighbors, c, 0.1)

			if len(matches) == (width*2 + 1) {
				reject(neighbors, mat, width, img, thresh, c)
			}
		}
	}

	// queue := make([]int, 0)
	// // Push
	// queue := append(queue, 1)
	// // Top (just get next element, don't remove it)
	// x = queue[0]
	// // Discard top element
	// queue = queue[1:]
	// // Is empty ?
	// if len(queue) == 0 {
	// 	fmt.Println("Queue is empty !")
	// }
}

// Identifies lines in a picture that have a color within thresh distance of a color in col
func huntLines(img image.Image, colors []color.Color, thresh float64, width int) (image.Image, []Line) {
	// We already have information about visited pixels from iter coordinates
	// This stops the algo from re-processing a rejected shape for every pixel
	visited := newTrackingMat(img.Bounds().Max.X+2, img.Bounds().Max.Y+1)

	iter(img, func(x, y int, c color.Color) {
		// If this pixel has been visited do not process it
		if visited.get(x, y) != 0 {
			return
		}

		// -- Hunting
		// Measure color distance between this pixel and target colors
		if !similarColor(c, colors, thresh) {
			return
		}

		// Get adjacent pixels
		neighbors := neighborPixels(x, y, width, img)

		// Determine if this is the start of a chunk or a line. Reject if chunk, Trace if line
		matches := scanChunk(neighbors, c, 0.1)

		// NOTE: only mark pixels in Chunking or Tracing, not when visiting. Otherwise visiting the edges
		// of a line will reject the line... I think?

		// -- Chunking
		// If the chunk all has the same color begin rejection: chunk and reject whole area
		if len(matches) == (width*2 + 1) {
			reject(neighbors, visited, width, img, thresh, c)
			return
		}

		visited.set(x, y, 2)

		// Note: a line that connects to a chunk should not be rejected

		// -- Tracing
		// Create a line object
		// Add matching pixels
		// Reject bad pixels

		// For all neighbors in group
		// If pixel matches and has not been visited repeat chunking
		// If no more matches found add group to return and resume hunt

		// -- Modeling
		// Merge close lines
		// Determine a center
		// Note: this isnt going into this function, put it somewhere else
	})

	ret := image.NewNRGBA(img.Bounds())

	for x := 0; x < visited.w; x++ {
		for y := 0; y < visited.h; y++ {

			if v := visited.get(x, y); v == 0 {
				ret.Set(x, y, color.Black)

			} else if v == 1 {
				ret.Set(x, y, color.NRGBA{225, 0, 0, 225})

			} else {
				ret.Set(x, y, color.White)
			}
		}
	}

	return ret, nil
}

// Tracks the results of a GroupLines operation
// 0 is univisited, 1 is rejected, 2 is line
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
	return m.arr[y*m.w+x]
}

func (m *TrackingMat) set(x, y int, v uint8) {
	m.arr[y*m.w+x] = v
}
