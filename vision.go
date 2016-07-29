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

// Returns a slice of lines from a provided image
func lineify(p *PixMatrix, colors []*Pix, thresh float64, width int) []*Line {
	var pixLines []*Pix

	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x += 2 {
			// Get the pixel from the matrix
			pix := p.get(x, y)

			// Check if the pixel is close to any of the target colors
			if !isClose(pix, colors, thresh) {
				continue
			}

			// Get adjacent pixels
			neighbors := neighbors(x, y, width, p)

			// Determine if this is the start of a chunk or a line. This thresh seems to do better higher
			matches := scanChunk(neighbors, pix, 0.4)

			// If all neighbors are matches then mark this as a chunk
			if len(matches) == len(neighbors) {
				markChunk(pix, p, neighbors, width, thresh)

			} else {
				// Else mark it as a line
				pixLines = append(pixLines, pix)
				pix.ptype = PIX_LINE
			}
		}
	}

	// Cluster line pixels into lines
	rawLines := cluster(p)

	// Scrub out lines that are most likely not lines
	outlines := filterLines(rawLines)

	return outlines
}

func lineifyExperimental(p *PixMatrix, colors []*Pix, thresh float64, width int) []*Line {
	for y := 0; y < p.h; y += 3 {
		for x := 0; x < p.w; x += 3 {
			// Get the pixel from the matrix
			pix := p.get(x, y)

			// This pixel has been visited by something or another
			// if pix.ptype == PIX_NOTHING {
			// 	continue
			// }

			// Colors dont match
			if !isClose(pix, colors, thresh) {
				continue
			}

			chunkFill(pix, p, 0.4)
		}
	}

	// Cluster line pixels into lines
	// rawLines := trace(pixLines, p)
	// fmt.Println("Returning lines: ", len(rawLines))
	return nil

	// Scrub out lines that are most likely not lines
	// outlines := filterLines(rawLines)
	// return outlines
}

func isClose(c *Pix, targets []*Pix, thresh float64) bool {
	for _, t := range targets {
		distance := colorDistance(c, t)

		if distance <= thresh {
			return true

			// Check to see if the color is even *close*. This saves all kinds of time when it comes to performance
		} else if distance > 1.2 {
			return false
		}
	}

	return false
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
func neighbors(tX, tY, distance int, m *PixMatrix) (ret []*Pix) {
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

// Like neighbors but only returns line pixels. (ptype == PIX_LINE)
func neighborsCluster(tX, tY, distance int, i *PixMatrix) (ret []*Pix) {
	for _, p := range neighbors(tX, tY, distance, i) {
		if p.ptype == PIX_LINE {
			ret = append(ret, p)
		}
	}

	return
}

// Return matching pixels from chunk that are within thresh of the given color
func scanChunk(chunk []*Pix, target *Pix, thresh float64) (ret []*Pix) {
	for _, p := range chunk {
		if d := colorDistance(p, target); d < thresh {
			ret = append(ret, p)
		}
	}

	return
}

// Mark chunk pixels in the matrix
func markChunk(c *Pix, m *PixMatrix, neighbor []*Pix, width int, thresh float64) (chunkPixels []*Pix) {
	for _, p := range neighbor {
		chunkPixels = append(chunkPixels, p)
		p.ptype = PIX_CHUNK
	}

	// extend the chunking to each of the neighboring pixels
	for _, p := range neighbor {
		for _, nearby := range neighbors(p.x, p.y, width, m) {
			if distance := colorDistance(nearby, c); distance > thresh {
				return
			}
		}
	}

	return
}

// Creates lines from pixels
func cluster(mat *PixMatrix) (ret []*Line) {
	mat.iter(func(x, y int, pix *Pix) {
		// Ignore non-line pixels or already added pixels
		if pix == nil || pix.ptype != PIX_LINE || pix.line != nil {
			return
		}

		// Create a new line and a queue
		q := []*Pix{pix}
		line := NewLine(0)
		ret = append(ret, line)

		for len(q) > 0 {
			// Pop the next pixel off the queue
			next := q[0]
			q = q[1:]

			// continue if next already belongs to a line
			if next.line != nil {
				continue
			}

			// Add this pixel to the line
			line.add(next)

			// Queue this pixels neighbors
			q = append(q, neighborsCluster(next.x, next.y, 1, mat)...)
		}
	})

	return
}
