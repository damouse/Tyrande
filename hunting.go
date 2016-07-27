package main

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

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/lucasb-eyer/go-colorful"
)

// Note: a line that connects to a chunk should not be rejected
func hunt(img image.Image, colors []color.Color, thresh float64, width int) []Line {
	chanChunks := make(chan []Pix, 0)
	chanLines := make(chan []Pix, 0)

	defer close(chanChunks)
	defer close(chanLines)

	chunks := []Pix{}
	lines := []Pix{}

	for _, c := range colors {
		go func(col color.Color) {
			ch, li := getLines(img, col, thresh, width)

			chanChunks <- ch
			chanLines <- li
		}(c)
	}

	done := 0
	for done != len(colors) {
		chunks = append(chunks, <-chanChunks...)
		lines = append(lines, <-chanLines...)

		done += 1
	}

	// Create a new tracking matrix containing all the points
	mat := newTrackingMat(img.Bounds().Max.X, img.Bounds().Max.Y)

	for _, p := range lines {
		mat.setPix(p)
	}

	// Test to make sure the matrix we made is sane
	// tester := outputMat(img.Bounds(), mat)
	// go save(tester, "mat.png")

	// Do we want to trace lines seperately?
	trueLines := cluster(lines, mat)

	// Do some coloring
	p := output(img.Bounds(), chunks, trueLines)
	save(p, "huntress.png")

	return nil
}

// Not screening for duplicates
func aggregate(mat [][]Pix) (ret []Pix) {
	for _, list := range mat {
		ret = append(ret, list...)
	}

	return
}

// Identifies lines in a picture that have a color within thresh distance of a color in col
// Returns lines and chunks
func getLines(img image.Image, target color.Color, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {

	iter(img, func(x, y int, c color.Color) {
		// Measure color distance between this pixel and target colors
		if distance := colorDistance(c, target); distance > thresh {
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

func cluster(points []Pix, mat *TrackingMat) (ret []*Line) {
	mat.iter(func(x, y int, pix *Pix) {
		// Ignore non-line pixels or already added pixels
		if pix == nil || pix.line != nil {
			return
		}

		// Create a new line
		line := NewLine(0)
		line.add(pix)

		fmt.Println("First point: ", pix)

		ret = append(ret, line)

		// Fetch neighbors of this pixel
		temp := neighborsCluster(pix.x, pix.y, 1, mat)

		for _, z := range temp {
			if z == nil {
				fmt.Println(z)
			} else {
				fmt.Printf("%v, isnil: %v\n", z, z.line == nil)
			}
		}

		neighbors := filter(temp)
		fmt.Println()

		for _, z := range neighbors {
			if z == nil {
				fmt.Println(z)
			} else {
				fmt.Printf("%v, isnil: %v\n", z, z.line == nil)
			}
		}

		// fmt.Printf("Neighbors: %v\n", neighbors[0])
		os.Exit(0)

		for _, p := range temp {
			if p != nil && p.line != nil {
				neighbors = append(neighbors, p)
			}
		}

		for len(neighbors) != 0 {
			// fmt.Println("Size of neighbors ", len(neighbors))

			// Mark neighbors that are also line pixels. Remove them from points
			for _, neigh := range neighbors {
				line.add(neigh)
			}

			newNeighbors := []*Pix{}

			// Add neighbors of neighbors
			for _, n := range neighbors {
				temp := neighborsCluster(n.x, n.y, 1, mat)
				newNeighbors = append(newNeighbors, filter(temp)...)

			}

			neighbors = newNeighbors
		}
	})

	return
}

func filter(a []*Pix) (ret []*Pix) {
	for _, p := range a {
		if p != nil && p.line == nil {
			ret = append(ret, p)
		}
	}

	return
}

// Create a line from the given points
func makeLine(points []Pix) (ret Line) {
	return
}

// output an image for testing purposes
func output(bounds image.Rectangle, chunks []Pix, lines []*Line) image.Image {
	ret := image.NewNRGBA(bounds)

	iter(ret, func(x, y int, c color.Color) {
		ret.Set(x, y, color.Black)
	})

	// for _, p := range chunks {
	// 	ret.Set(p.x, p.y, color.White)
	// }

	for _, line := range lines {
		rcolor := colorful.FastHappyColor()

		for _, pix := range line.pixels {
			ret.Set(pix.x, pix.y, rcolor)
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
		if x < 0 || x > i.h {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < 0 || y > i.w {
				continue
			}

			ret = append(ret, i.get(x, y))
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
