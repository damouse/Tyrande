package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/mdesenfants/gokmeans"
)

// Modeling and detecting on-screen players
type Line struct {
	pixels []Pix
	id     int
	cX, cY int // center
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

// Note: a line that connects to a chunk should not be rejected

func hunt(img image.Image, colors []color.Color, thresh float64, width int) []Line {
	chunks := make(chan []Pix, 0)
	lines := make(chan []Pix, 0)

	allchunks := [][]Pix{}
	alllines := [][]Pix{}

	for _, c := range colors {
		go func(col color.Color) {
			ch, li := getLines(img, col, thresh, width)

			chunks <- ch
			lines <- li
		}(c)
	}

	done := 0

	for done != len(colors) {
		allchunks = append(allchunks, <-chunks)
		alllines = append(alllines, <-lines)

		done += 1
	}

	close(chunks)
	close(lines)

	aggChunks := aggregate(allchunks)
	aggLines := aggregate(alllines)

	trueLines := cluster(aggLines, 4, 50)

	// Do some coloring
	p := output(img.Bounds(), aggChunks, trueLines)
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

func cluster(points []Pix, groups int, iterations int) (ret []Line) {

	linePoints := []gokmeans.Node{}
	for _, p := range points {
		linePoints = append(linePoints, gokmeans.Node{float64(p.x), float64(p.y)})
	}

	// Run kmeans on the lines
	// Get a list of centroids and output the values
	if success, centroids := gokmeans.Train(linePoints, 2, 25); success {
		// Show the centroids
		fmt.Println("The centroids are")

		for i, centroid := range centroids {
			ret = append(ret, Line{id: i, cX: int(centroid[0]), cY: int(centroid[1])})
			fmt.Println(centroid)
		}

		for i, observation := range linePoints {
			index := gokmeans.Nearest(observation, centroids)
			// fmt.Println(observation, "belongs in cluster", index+1, ".")

			ret[index].pixels = append(ret[index].pixels, points[i])
		}
	}

	return
}

func colorify(img image.Image, pix []Pix, c color.Color) {
	// for _, p := range pix {
	// 	ret.Set(p.x, p.y, color.NRGBA{255, 0, 0, 255})
	// }
}

// output an image for testing purposes
func output(bounds image.Rectangle, chunks []Pix, lines []Line) image.Image {
	ret := image.NewNRGBA(bounds)

	iter(ret, func(x, y int, c color.Color) {
		ret.Set(x, y, color.Black)
	})

	for _, p := range chunks {
		ret.Set(p.x, p.y, color.NRGBA{255, 0, 0, 255})
	}

	for i, line := range lines {
		for _, pix := range line.pixels {
			if i == 0 {
				ret.Set(pix.x, pix.y, color.White)

			} else if i == 1 {
				ret.Set(pix.x, pix.y, color.NRGBA{255, 0, 0, 255})

			} else if i == 2 {
				ret.Set(pix.x, pix.y, color.NRGBA{0, 255, 0, 255})

			} else if i == 3 {
				ret.Set(pix.x, pix.y, color.NRGBA{255, 0, 255, 255})

			}
		}
	}

	return ret
}

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

// Identifies lines in a picture that have a color within thresh distance of a color in col
// Returns lines and chunks
func getLines(img image.Image, target color.Color, thresh float64, width int) (chunkPixels []Pix, linePixels []Pix) {
	iter(img, func(x, y int, c color.Color) {
		// Measure color distance between this pixel and target colors
		if distance := colorDistance(c, target); distance < thresh {
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
					if distance := colorDistance(nearby.Color, c); distance < thresh {
						return
					}
				}
			}

		} else {
			chunkPixels = append(chunkPixels, Pix{c, x, y})
		}
	})

	return
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
