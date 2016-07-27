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
	"image"
	"image/color"

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
	mat := newTrackingMat(img.Bounds().Max.X+2, img.Bounds().Max.Y+1)

	for _, p := range lines {
		mat.set(p.x, p.y, &p)
	}

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

func cluster(points []Pix, mat *TrackingMat) (ret []Line) {
	counter := 0

	for len(points) > 0 {
		// Pop the first pixel, add it to a line
		pix := points[0]
		points = points[0:]

		line := Line{[]Pix{pix}, counter, 0, 0}
		ret = append(ret, line)
		counter += 1

		// Fetch neighbors
		neighbors := neighborPixelsMat(pix.x, pix.y, 1, mat)

		// Mark all neighbors
		// Recursively extend
	}

External:
	for _, p := range points {
		// Check neighbors of this pixel for line membership
		// If a neighbor match is found add this pixel to that line
		// Else create new line and add this pixel

		// Skip if this point has already been added
		for _, line := range ret {
			for _, pix := range line.pixels {
				if pix == p {
					continue External
				}
			}
		}

		// This pixel has not been added to a line. Create a new line now
		ret = append(ret, Line{[]Pix{p}, len(ret), 0, 0})

		// Search neighbors for
	}

	return

	// Old kmeans clustering
	// linePoints := []gokmeans.Node{}
	// for _, p := range points {
	// 	linePoints = append(linePoints, gokmeans.Node{float64(p.x), float64(p.y)})
	// }

	// // Run kmeans on the lines
	// // Get a list of centroids and output the values
	// if success, centroids := gokmeans.Train(linePoints, 3, 25); success {
	// 	// Show the centroids
	// 	fmt.Println("The centroids are")

	// 	for i, centroid := range centroids {
	// 		ret = append(ret, Line{id: i, cX: int(centroid[0]), cY: int(centroid[1])})
	// 		fmt.Println(centroid)
	// 	}

	// 	for i, observation := range linePoints {
	// 		index := gokmeans.Nearest(observation, centroids)
	// 		// fmt.Println(observation, "belongs in cluster", index+1, ".")

	// 		ret[index].pixels = append(ret[index].pixels, points[i])
	// 	}
	// }

	// fmt.Printf("Clustering completing with %d clusters\n", len(ret))
	// return
}

// output an image for testing purposes
func output(bounds image.Rectangle, chunks []Pix, lines []Line) image.Image {
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

// returns all pixels within range d of target coordinates x, y
// The pixel in the middle is included in results
func neighborPixelsMat(tX, tY, distance int, i *TrackingMat) (ret []Pix) {
	for x := tX - distance; x <= tX+distance; x++ {
		if x < 0 || x > i.w {
			continue
		}

		for y := tY - distance; y <= tY+distance; y++ {
			if y < 0 || y > i.h {
				continue
			}

			ret = append(ret, Pix{i.get(x, y), x, y, nil})
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
