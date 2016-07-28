package main

// type OrderedPair struct {
// 	x, y int
// }

var mods = [...]struct {
	x, y int
}{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	{-1, -1}, {1, 1}, {1, -1}, {-1, 1},
}

// func FloodFill(graph [][]int, origin OrderedPair) []OrderedPair {
// 	val, _ := graph[origin.y][origin.x]

// 	// Create a visited list
// 	seen := make([][]bool, len(graph))

// 	for i, row := range graph {
// 		seen[i] = make([]bool, len(row))
// 	}

// 	// let go sort out the appended size.
// 	fill := []OrderedPair{}

// 	// go will shuffle memory too when adding/removing items from q
// 	q := []OrderedPair{origin}

// 	for len(q) > 0 {

// 		// shift the q
// 		op := q[0]
// 		q = q[1:]

// 		if seen[op.y][op.x] {
// 			continue
// 		}

// 		seen[op.y][op.x] = true
// 		fill = append(fill, op)

// 		for _, mod := range mods {
// 			newx := op.x + mod.x
// 			newy := op.y + mod.y

// 			if 0 <= newy && newy < len(graph) && 0 <= newx && newx < len(graph[newy]) {
// 				if graph[newy][newx] == val {
// 					q = append(q, OrderedPair{newx, newy})
// 				}
// 			}
// 		}
// 	}

// 	return fill
// }
