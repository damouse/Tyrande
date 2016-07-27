package main

// -- Modeling
// Merge close lines
// Determine a center
// Determine a top

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

// // Note: a line that connects to a chunk should not be rejected
// func hunt(img image.Image, colors []color.Color, thresh float64, width int) ([]Pix, []*Line) {
// 	chunks, lines := getLines(img, colors, thresh, width)

// 	// Create a new tracking matrix containing all the points
// 	mat := newTrackingMat(img.Bounds().Max.X, img.Bounds().Max.Y)
// 	for _, p := range lines {
// 		mat.setPix(p)
// 	}

// 	// Try masking lines with chunks
// 	// for _, p := range chunks {
// 	// 	mat.set(p.x, p.y, nil)
// 	// }

// 	// Do we want to trace lines seperately?
// 	raw := cluster(lines, mat)

// 	// Filter out obvious non-outlines
// 	filtered := filterLines(raw)

// 	return chunks, filtered
// }

// Note: a line that connects to a chunk should not be rejected
func hunt(img image.Image, colors []color.Color, thresh float64, width int) ([]Pix, []*Line) {
	// chanChunks := make(chan []Pix, 0)
	// chanLines := make(chan []Pix, 0)

	// defer close(chanChunks)
	// defer close(chanLines)

	// chunks := []Pix{}
	// lines := []Pix{}

	// for _, c := range colors {
	// 	go func(col color.Color) {
	// ch, li := getLines(img, col, thresh, width)

	chunks, lines := getLines(img, colors, thresh, width)

	// 		chanChunks <- ch
	// 		chanLines <- li
	// 	}(c)
	// }

	// done := 0
	// for done != len(colors) {
	// 	chunks = append(chunks, <-chanChunks...)
	// 	lines = append(lines, <-chanLines...)

	// 	done += 1
	// }

	// Create a new tracking matrix containing all the points
	mat := newTrackingMat(img.Bounds().Max.X, img.Bounds().Max.Y)

	for _, p := range lines {
		mat.setPix(p)
	}

	// Do we want to trace lines seperately?
	raw := cluster(lines, mat)

	// Filter out obvious non-outlines
	filtered := filterLines(raw)

	//
	// Do some coloring and save the image for demo
	// p := output(img.Bounds(), chunks, filtered)
	// save(p, "huntress.png")

	return chunks, filtered
}

// Identifies lines in a picture that have a color within thresh distance of a color in col
// Returns lines and chunks
func getLines(img image.Image, colors []color.Color, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {

	iter(img, func(x, y int, c color.Color) {
		// Measure color distance between this pixel and target colors
		// if distance := colorDistance(c, target); distance > thresh {
		// 	return
		// }

		isClose := false
		for _, target := range colors {
			// if distance := colorDistance(c, target); distance > thresh {
			if distance := colorDistance(c, target); distance <= thresh {
				isClose = true
				break
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
			linePixels = append(linePixels, Pix{c, x, y, nil})
		}
	})

	return
}

// // Identifies lines in a picture that have a color within thresh distance of a color in col
// // Returns lines and chunks
// func getLines(img image.Image, colors []color.Color, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {

// 	iter(img, func(x, y int, c color.Color) {
// 		// Measure color distance between this pixel and target colors
// 		for _, target := range colors {
// 			// if distance := colorDistance(c, target); distance > thresh {
// 			if distance := colorDistance(c, target); distance < thresh {
// 				break
// 			}

// 			return
// 		}

// 		// Get adjacent pixels
// 		neighbors := neighborPixels(x, y, width, img)

// 		// Determine if this is the start of a chunk or a line. Reject if chunk, Trace if line
// 		// This thresh seems to do better higher
// 		matches := scanChunk(neighbors, c, 0.4)

// 		// -- Chunking
// 		if len(matches) == len(neighbors) {
// 			for _, p := range neighbors {
// 				chunkPixels = append(chunkPixels, p)
// 			}

// 			// extend the chunking to each of the neighboring pixels
// 			for _, p := range neighbors {
// 				for _, nearby := range neighborPixels(p.x, p.y, width, img) {
// 					if distance := colorDistance(nearby.Color, c); distance > thresh {
// 						return
// 					}
// 				}
// 			}

// 		} else {
// 			linePixels = append(linePixels, Pix{c, x, y, nil})
// 		}
// 	})

// 	return
// }

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

// output an image for testing purposes
func output(bounds image.Rectangle, chunks []Pix, lines []*Line) image.Image {
	ret := image.NewNRGBA(bounds)

	iter(ret, func(x, y int, c color.Color) {
		ret.Set(x, y, color.Black)
	})

	for _, line := range lines {
		rcolor := colorful.FastHappyColor()

		for _, pix := range line.pixels {
			ret.Set(pix.x, pix.y, rcolor)
		}
	}

	// Write out the chunks in white
	if DEBUG_DRAW_CHUNKS {
		for _, p := range chunks {
			ret.Set(p.x, p.y, color.White)
		}
	}

	return ret
}

func outputMat(bounds image.Rectangle, mat *TrackingMat) image.Image {
	ret := image.NewNRGBA(bounds)

	mat.iter(func(x, y int, c *Pix) {
		if c == nil {
			ret.Set(x, y, color.Black)
		} else {
			ret.Set(c.x, c.y, c.Color)
		}
	})

	return ret
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

			ret = append(ret, Pix{i.At(x, y), x, y, nil})
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
