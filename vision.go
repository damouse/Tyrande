package main

type Line struct {
	pixels []*Pix

	avg, max, min Vec // average, max, and min of all points
}

func lineify(p *PixMatrix, colors []*Pix, thresh float64, width int) (lines []*Line) {
	filltag := 0
	fillThresh := 0.4

	for y := 0; y < p.h; y += 5 {
		for x := 0; x < p.w; x += 5 {
			current := p.get(x, y)

			// This pixel has been filled
			if current.ptype != PIX_NOTHING {
				continue
			}

			// Colors dont match
			if !isClose(current, colors, thresh) {
				continue
			}

			// Get adjacents
			adj := p.adjacentSimilarColor(current, current, 1, fillThresh)

			// If this is a chunk continue
			if len(adj) > 7 {
				continue
			}

			// Create a new line object
			line := &Line{}
			lines = append(lines, line)

			// Begin fill loop
			q := []*Pix{current}

			// Pixels visited as part of this fill are tagged with this int
			filltag += 1

			for len(q) > 0 {
				pix := q[0]
				q = q[1:]

				// Get neighbors that are of a similar color
				adj := p.adjacentSimilarColor(pix, current, 1, fillThresh)

				// Stop immediately if a chunk is found
				if len(adj) > 7 {
					// Shouldn't this be "continue?" Yes, yes it should.
					// For a reason I cant possibly explain it works with break and not with continue.
					// I dont remember any other time a bug accidentally made my algo work.
					break
				}

				// Mark as line
				pix.ptype = PIX_LINE
				line.add(pix)

				// Queue all non-visited neighbors (not visited means not queued)
				for _, p := range adj {
					if p.ptype == PIX_NOTHING && p.filltag != filltag {
						p.filltag = filltag
						q = append(q, p)
					}
				}
			}
		}
	}

	return
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

// Update center of screen coordinates. This method seems bulky and not well designed. Consider removing it
func getCenter(mat *PixMatrix) Vec {
	x, y := mat.center()
	return Vec{x + CENTER_OFFSET.x, y + CENTER_OFFSET.y}
}

// Calculate statistics about this line
func (l *Line) process() {
	var sumX, sumY int

	for _, p := range l.pixels {
		if p.x < l.min.x {
			l.min.x = p.x
		}

		if p.y < l.min.y {
			l.min.y = p.y
		}

		if p.x > l.max.x {
			l.max.x = p.x
		}

		if p.y > l.max.y {
			l.max.y = p.y
		}

		sumX += p.x
		sumY += p.y
	}

	l.avg.x = sumX / len(l.pixels)
	l.avg.y = sumY / len(l.pixels)
}

// Filter lines that dont look like actual lines
// Note: density may also be a good measure of "lineiness"
func filterLines(lines []*Line) (ret []*Line) {
	for _, l := range lines {
		if len(l.pixels) < 150 {
			for _, p := range l.pixels {
				p.ptype = PIX_CHUNK
			}

			continue
		}

		ret = append(ret, l)
	}

	return
}

// Lines
func (l *Line) add(p *Pix) {
	l.pixels = append(l.pixels, p)
}

func (l *Line) addAll(p []*Pix) {
	for _, a := range p {
		l.add(a)
	}
}

func (l *Line) merge(o *Line) {
	for _, p := range o.pixels {
		l.add(p)
	}
}

//

// We can use this to bound the search distance for the sake of performance
// halfX := img.Bounds().Max.X / 3
// halfY := img.Bounds().Max.Y / 3
// if x < halfX || x > halfX*2 || y < halfY || y > halfY*2 {
// 	return
// }
