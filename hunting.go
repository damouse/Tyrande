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

// Identifies lines in a picture that have a color within thresh distance of a color in col
func huntLines(img image.Image, colors []color.Color, thresh float64, width int) (image.Image, []Line) {
	chunks := newTrackingMat(img.Bounds().Max.X+2, img.Bounds().Max.Y+1)
	lines := newTrackingMat(img.Bounds().Max.X+2, img.Bounds().Max.Y+1)

	// linePixels := []Pix{}
	// linePoints := []gokmeans.Node{}

	iter(img, func(x, y int, c color.Color) {
		// If this is marked as chunk do not process
		// if chunks.get(x, y) != 0 {
		// 	return
		// }

		// Measure color distance between this pixel and target colors
		if !similarColor(c, colors, thresh) {
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
				chunks.set(p.x, p.y, 1)
			}

			// extend the chunking to each of the neighboring pixels
			for _, p := range neighbors {
				for _, nearby := range neighborPixels(p.x, p.y, width, img) {
					if similarColor(nearby.Color, colors, thresh) {
						chunks.set(nearby.x, nearby.y, 1)
					}
				}
			}

		} else {
			lines.set(x, y, 1)

			// linePixels = append(linePixels, Pix{c, x, y})
			// linePoints = append(linePoints, gokmeans.Node{float64(x), float64(y)})
		}

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

	// Run kmeans on the lines
	// Get a list of centroids and output the values
	// if success, centroids := gokmeans.Train(linePoints, 2, 25); success {
	// 	// Show the centroids
	// 	fmt.Println("The centroids are")

	// 	for _, centroid := range centroids {
	// 		fmt.Println(centroid)
	// 	}

	// 	// Output the clusters
	// 	fmt.Println("...")
	// 	for i, observation := range linePoints {
	// 		index := gokmeans.Nearest(observation, centroids)
	// 		// fmt.Println(observation, "belongs in cluster", index+1, ".")

	// 		pix := linePixels[i]
	// 		lines.set(pix.x, pix.y, uint8(index+1))
	// 	}
	// }

	ret := image.NewNRGBA(img.Bounds())

	for x := 0; x < chunks.w; x++ {
		for y := 0; y < chunks.h; y++ {

			// Turns chunks red
			if v := chunks.get(x, y); v == 1 {
				ret.Set(x, y, color.NRGBA{255, 0, 0, 255})

				// if v := chunks.get(x, y); v == 1 {
				// 	ret.Set(x, y, color.Black)

			} else if v := lines.get(x, y); v == 1 {
				ret.Set(x, y, color.White)

			} else if v := lines.get(x, y); v == 2 {
				ret.Set(x, y, color.NRGBA{255, 0, 0, 255})

			} else if v := lines.get(x, y); v == 3 {
				ret.Set(x, y, color.NRGBA{0, 255, 0, 255})

			} else if v := lines.get(x, y); v == 4 {
				ret.Set(x, y, color.NRGBA{255, 0, 255, 255})

			} else {
				ret.Set(x, y, color.Black)
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
