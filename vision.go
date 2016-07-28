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
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Returns a slice of lines from a provided image
func hunt(img image.Image, colors []color.Color, thresh float64, width int) ([]Pix, []*Line) {
	// Convert the image to a more usable format
	// pixMatrix := convertImage(img)
	convertImage(img)

	// Extract lines and chunks
	chunks, lines := extract(img, colors, thresh, width)

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
func extract(img image.Image, colors []color.Color, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {
	// Can we just convert image and colors to LUV here once and then not bother with it again?

	// output an image for testing purposes
	loveImage := NewPixMat(img.Bounds().Max.X, img.Bounds().Max.Y)

	iter(img, func(x, y int, c color.Color) {
		l1, u1, v1 := convertToColorful(c).Luv()
		loveImage.set(x, y, &colorful.Color{l1, u1, v1})
	})

	lovelyTargets := []*colorful.Color{}

	for _, c := range colors {
		l1, u1, v1 := convertToColorful(c).Luv()
		lovelyTargets = append(lovelyTargets, &colorful.Color{l1, u1, v1})
	}

	iter(img, func(x, y int, c color.Color) {
		// Manually recasting the luv color
		// luv := love.At(x, y).(colorful.Color)
		// luv := convertToColorful(c)

		isClose := false
		for i, _ := range colors {
			// distance := colorDistance(c, target)

			// l1, u1, v1 := convertToColorful(c).Luv()
			// l2, u2, v2 := convertToColorful(target).Luv()

			// man := luvSample.DistanceLuv(luvTarget)
			// man := math.Sqrt(sq(l1-l2) + sq(u1-u2) + sq(v1-v2))

			// Testing manual conversion
			c := loveImage.get(x, y)
			t := lovelyTargets[i]
			man := math.Sqrt(sq(c.R-t.R) + sq(c.G-t.G) + sq(c.B-t.B))

			// fmt.Println(distance, man)
			// end testing manual conversion

			distance := man

			if distance <= thresh {
				isClose = true
				break

				// Check to see if the color is even *close*. This saves all kinds of time when it comes to performance
			} else if distance > 1.2 {
				return
			}
		}

		if !isClose {
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

// Return matching pixels from chunk that are within thresh of the given color
func scanChunk(chunk []Pix, color color.Color, thresh float64) (ret []Pix) {
	for _, p := range chunk {
		if d := colorDistance(p, color); d < thresh {
			ret = append(ret, p)
		}
	}

	return
}
