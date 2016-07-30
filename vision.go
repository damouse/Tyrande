package main

import (
	"image"
	"time"
)

// Main event loop. Continuously captures the screen and places results into visionQ for the main loop to pick up
func vision() {
	var p image.Image
	op := Cycle{}
	op.start = time.Now()

	if DEBUG_STATIC {
		p = imageStatic
	} else {
		p = CaptureLeft()
	}

	op.mat = convertImage(p)
	op.lines = lineify(op.mat, SWATCH, COLOR_THRESHOLD, LINE_WIDTH)
	op.vision = time.Now()

	updateCenter(&op)
	visionChan <- op
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

// Update center of screen coordinates
func updateCenter(op *Cycle) {
	x, y := op.mat.center()

	x += CENTER_OFFSET.x
	y += CENTER_OFFSET.y

	centerVector.x = x
	centerVector.y = y

	op.center.x = x
	op.center.y = y
}

// We can use this to bound the search distance for the sake of performance
// halfX := img.Bounds().Max.X / 3
// halfY := img.Bounds().Max.Y / 3
// if x < halfX || x > halfX*2 || y < halfY || y > halfY*2 {
// 	return
// }
