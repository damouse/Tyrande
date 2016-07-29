package main

var mods = [...]struct {
	x, y int
}{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	{-1, -1}, {1, 1}, {1, -1}, {-1, 1},
}

// Perform a heavy fill
func chunkFill(pix *Pix, mat *PixMatrix, thresh float64) {
	q := []*Pix{pix}

	// Maintain a visited list for this fill operation
	visited := NewPixMatrix(mat.w, mat.h)

	for len(q) > 0 {
		// pop off the queue
		pix := q[0]
		q = q[1:]

		// Mark visted
		visited.set(pix)

		// Get neighbors that are of a similar color
		adj := mat.adjacentSimilarColor(pix, 1, thresh)

		// Check chunkiness

		// If chunky
		if chunky(pix, mat) {
			pix.ptype = PIX_CHUNK

			for _, p := range mat.adjacent(pix, 2) {
				p.ptype = PIX_CHUNK
			}
		} else {
			pix.ptype = PIX_LINE
		}

		// Queue all non-visited neighbors
		for _, p := range adj {
			if isVisited := visited.get(p.x, p.y); isVisited == nil {
				q = append(q, p)
			}
		}

		// this has been visited
		// if op.ptype == PIX_CHUNK || op.ptype == PIX_LINE {
		// 	continue
		// }

		// if chunky(op, mat) {
		// 	op.ptype = PIX_CHUNK
		// 	nearby := neighbors(op.x, op.y, 1, mat)
		// 	markChunk(op, mat, nearby, 1, thresh)
		// } else {
		// 	op.ptype = PIX_LINE
		// }

		// for _, mod := range mods {
		// 	newx := op.x + mod.x
		// 	newy := op.y + mod.y

		// 	if 0 <= newy && newy < mat.h && 0 <= newx && newx < mat.w {
		// 		next := mat.get(newx, newy)

		// 		if chunky(next, mat) {
		// 			next.ptype = PIX_CHUNK
		// 			continue
		// 		}

		// 		// Only append the next pixel if its NOT chunky
		// 		if dist := colorDistance(pix, next); dist <= thresh {
		// 			q = append(q, next)
		// 		}
		// 	}
		// }
	}
}

// Returns true if the surounding pixels are of a similar color
func chunky(p *Pix, mat *PixMatrix) bool {
	// Get adjacent pixels
	surrouding := neighbors(p.x, p.y, 1, mat)

	// Determine if this is the start of a chunk or a line. This thresh seems to do better higher
	matches := scanChunk(surrouding, p, 0.4)
	return len(matches) == len(surrouding)
}
