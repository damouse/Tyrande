package main

var mods = [...]struct {
	x, y int
}{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	{-1, -1}, {1, 1}, {1, -1}, {-1, 1},
}

// Perform a heavy fill
func chunkFill(pix *Pix, mat *TrackingMat, thresh float64) {
	q := []*Pix{pix}

	for len(q) > 0 {
		// pop off the queue
		op := q[0]
		q = q[1:]

		// this has been visited
		if op.ptype == PIX_CHUNK || op.ptype == PIX_LINE {
			continue
		}

		if chunky(op, mat) {
			op.ptype = PIX_CHUNK
			nearby := neighbors(op.x, op.y, 1, mat)
			markChunk(op, mat, nearby, 1, thresh)
		} else {
			op.ptype = PIX_LINE
		}

		for _, mod := range mods {
			newx := op.x + mod.x
			newy := op.y + mod.y

			if 0 <= newy && newy < mat.h && 0 <= newx && newx < mat.w {
				next := mat.get(newx, newy)

				if chunky(next, mat) {
					next.ptype = PIX_CHUNK
					continue
				}

				// Only append the next pixel if its NOT chunky
				if dist := colorDistance(pix, next); dist <= thresh {
					q = append(q, next)
				}
			}
		}
	}
}

// Returns true if the surounding pixels are of a similar color
func chunky(p *Pix, mat *TrackingMat) bool {
	// Get adjacent pixels
	surrouding := neighbors(p.x, p.y, 1, mat)

	// Determine if this is the start of a chunk or a line. This thresh seems to do better higher
	matches := scanChunk(surrouding, p, 0.4)
	return len(matches) == len(surrouding)
}