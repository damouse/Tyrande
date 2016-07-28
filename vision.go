package main

// -- Modeling
// Merge close lines
// Determine a center
// Determine a top

// We can use this to bound the search distance for the sake of performance
// halfX := img.Bounds().Max.X / 3
// halfY := img.Bounds().Max.Y / 3

// if x < halfX || x > halfX*2 || y < halfY || y > halfY*2 {
// 	return
// }
/*
TODO:
	Clean pipeline
	Test live
	Hunting performance
*/

import (
	"image"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// Returns a slice of lines from a provided image
func hunt(img image.Image, colors []*Pix, thresh float64, width int) ([]Pix, []*Line) {
	// Convert the image to a more usable format
	p := convertImage(img)

	var pixChunks, pixLines []*Pix

	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			// Get the pixel from the matrix
			pix := p.get(x, y)

			// Check if the pixel is close to any of the target colors
			if !isClose(pix, colors, thresh) {
				continue
			}

			// Get adjacent pixels
			neighbors := neighborPix(x, y, width, p)

			// Determine if this is the start of a chunk or a line. This thresh seems to do better higher
			matches := scanChunkPix(neighbors, pix, 0.4)

			// If all neighbors are matches then mark this as a chunk
			if len(matches) == len(neighbors) {
				pixChunks = append(pixChunks, markChunk(pix, p, neighbors, width, thresh)...)

			} else {
				// Else mark it as a line
				pixLines = append(pixLines, pix)
				pix.ptype = PIX_LINE
			}
		}
	}

	rawLines := clusterPix(p)

	finalLines := filterLines(rawLines)

	// Color lines for debug output
	for _, l := range finalLines {
		rcolor := colorful.FastHappyColor()

		for _, a := range l.pixels {
			a.Color = rcolor
		}
	}

	p.save("matrix.png")

	// Extract lines and chunks
	chunks, lines := extract(img, p, colors, thresh, width)

	// Create a new tracking matrix containing all the points
	mat := newTrackingMat(img.Bounds().Max.X, img.Bounds().Max.Y)
	for _, p := range lines {
		mat.setPix(p)
	}

	// Mask the lines by nilling out chunks in the matrix
	for _, p := range chunks {
		mat.set(p.x, p.y, nil)
	}

	// Do we want to trace lines seperately?
	raw := cluster(lines, mat)

	// Filter out obvious non-outlines
	filtered := filterLines(raw)

	return chunks, filtered
}

// Identifies lines in a picture that have a color within thresh distance of a color in col
// Returns lines and chunks
func extract(img image.Image, mat *TrackingMat, colors []*Pix, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {

	iter(img, func(x, y int, c color.Color) {
		if !isClose(mat.get(x, y), colors, thresh) {
			return
		}

		// Get adjacent pixels
		neighbors := neighborPixels(x, y, width, img)

		// Determine if this is the start of a chunk or a line. Reject if chunk, Trace if line
		// This thresh seems to do better higher
		matches := scanChunk(neighbors, c, 0.4)

		// -- Chunking
		if len(matches) == len(neighbors) {
			for _, p := range neighbors {
				chunkPixels = append(chunkPixels, p)
			}

			// extend the chunking to each of the neighboring pixels
			for _, p := range neighbors {
				for _, nearby := range neighborPixels(p.x, p.y, width, img) {
					if distance := colorDistance(nearby.Color, c); distance > thresh {
						return
					}
				}
			}

		} else {
			linePixels = append(linePixels, *NewPix(x, y, c))
		}
	})

	return
}

func isClose(c *Pix, targets []*Pix, thresh float64) bool {
	for _, t := range targets {
		distance := tyrDistance(c, t)

		if distance <= thresh {
			return true

			// Check to see if the color is even *close*. This saves all kinds of time when it comes to performance
		} else if distance > 1.2 {
			return false
		}
	}

	return false
}

// Creates lines from pixels
func cluster(points []Pix, mat *TrackingMat) (ret []*Line) {
	mat.iter(func(x, y int, pix *Pix) {
		// Ignore non-line pixels or already added pixels
		if pix == nil || pix.line != nil {
			return
		}

		q := []*Pix{pix}
		line := NewLine(0)
		ret = append(ret, line)

		for len(q) > 0 {
			next := q[0]
			q = q[1:]

			// continue if next is marked
			if next.line != nil {
				continue
			}

			// Add next
			line.add(next)

			// Queue neighbors
			q = append(q, neighborsCluster(next.x, next.y, 1, mat)...)
		}
	})

	return
}

// Filter lines that dont look like actual lines
func filterLines(lines []*Line) (ret []*Line) {
	for _, l := range lines {
		if len(l.pixels) < 100 {
			for _, p := range l.pixels {
				p.ptype = PIX_CHUNK
			}

			continue
		}

		ret = append(ret, l)
	}

	return
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

			ret = append(ret, *NewPix(x, y, i.At(x, y)))
		}
	}

	return
}

// Like neighbors, but does not return pixels that are nil or already marked
func neighborsCluster(tX, tY, distance int, i *TrackingMat) (ret []*Pix) {
	for x := tX - distance; x <= tX+distance; x++ {
		if x < 0 || x >= i.w {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < 0 || y >= i.h {
				continue
			}

			// fmt.Printf("%d, %d\t%d %d\n", x, y, i.w, i.h)
			if p := i.get(x, y); p != nil && p.line == nil {
				ret = append(ret, p)
			}
		}
	}

	return
}

// Like neighbors, but does not return pixels that are nil or already marked
func neighborsClusterPix(tX, tY, distance int, i *TrackingMat) (ret []*Pix) {
	for x := tX - distance; x <= tX+distance; x++ {
		if x < 0 || x >= i.w {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < 0 || y >= i.h {
				continue
			}

			// fmt.Printf("%d, %d\t%d %d\n", x, y, i.w, i.h)
			if p := i.get(x, y); p != nil && p.ptype == PIX_LINE && p.line == nil {
				ret = append(ret, p)
			}
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

// returns all pixels within range d of target coordinates x, y
// The pixel in the middle is included in results
func neighborPix(tX, tY, distance int, m *TrackingMat) (ret []*Pix) {
	for x := tX - distance; x <= tX+distance; x++ {
		if x < 0 || x > m.w {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < 0 || y > m.h {
				continue
			}

			ret = append(ret, m.get(x, y))
		}
	}

	return
}

// Return matching pixels from chunk that are within thresh of the given color
func scanChunkPix(chunk []*Pix, target *Pix, thresh float64) (ret []*Pix) {
	for _, p := range chunk {
		if d := tyrDistance(p, target); d < thresh {
			ret = append(ret, p)
		}
	}

	return
}

// Mark chunk pixels in the matrix
func markChunk(c *Pix, m *TrackingMat, neighbors []*Pix, width int, thresh float64) (chunkPixels []*Pix) {
	for _, p := range neighbors {
		chunkPixels = append(chunkPixels, p)
		p.ptype = PIX_CHUNK
	}

	// extend the chunking to each of the neighboring pixels
	for _, p := range neighbors {
		for _, nearby := range neighborPix(p.x, p.y, width, m) {
			if distance := tyrDistance(nearby, c); distance > thresh {
				return
			}
		}
	}

	return
}

// Creates lines from pixels
func clusterPix(mat *TrackingMat) (ret []*Line) {
	mat.iter(func(x, y int, pix *Pix) {
		// Ignore non-line pixels or already added pixels
		if pix == nil || pix.ptype != PIX_LINE || pix.line != nil {
			return
		}

		q := []*Pix{pix}
		line := NewLine(0)
		ret = append(ret, line)

		for len(q) > 0 {
			next := q[0]
			q = q[1:]

			// continue if next is marked
			if next.line != nil {
				continue
			}

			// Add next
			line.add(next)

			// Queue neighbors
			q = append(q, neighborsClusterPix(next.x, next.y, 1, mat)...)
		}
	})

	return
}
