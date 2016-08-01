package main

import (
	"fmt"
	"sort"
	"time"
)

type Char struct {
	center     Vector
	offset     Vector // vector from this char to the center of the screen
	offsetDist float64
}

// var lastUpdate time.Time = time.Now()

// Gets updates from vision and updates Chars
func modeling() {
	op := <-visionChan

	op.lines = filterLines(op.lines)

	for _, l := range op.lines {
		l.process()
	}

	x, y := op.mat.center()
	op.chars = buildChars(op.lines, Vector{x, y})

	// for _, l := range op.lines {
	// 	op.chars = append(op.chars, &Char{l, Vector{}})
	// }

	op.model = time.Now()
	op.bench()

	if DEBUG_SAVE_LINES {
		op.save("huntress.png")
		stop()
	}

	if DEBUG_WINDOW {
		// window.show(op.mat.toImage())
		window.queueShow(&op)
	}

	if len(op.chars) != len(Chars) {
		debug("Targets: %d", len(op.chars))
	}

	// Update shared store
	CharLock.Lock()
	Chars = op.chars
	CharLock.Unlock()

	// log("Update: %s", time.Since(lastUpdate))
	// lastUpdate = time.Now()
}

// Build a list of characters from a list of lines, sets Char data, and orders return based on closest
func buildChars(lines []*Line, center Vector) (ret []*Char) {
	for _, l := range lines {
		c := &Char{}

		c.center = Vector{l.centerX, l.centerY}
		c.offset = Vector{center.x - c.center.x, center.y - c.center.y}
		c.offsetDist = euclideanDistanceVec(center, c.center)

		ret = append(ret, c)
	}

	sort.Sort(ByDistance(ret))
	return
}

// Targeting
// Return the line whose center is closest to the screen center. If no lines passed, returns nil
func closestCenter(chars []*Char, center Vector) (ret *Char) {
	return chars[0]
}

// func (c *Char) centerOffset()

func printChars(chars []*Char) {
	for _, c := range chars {
		fmt.Printf("\tcenter: %d, %d\n", c.offset.x, c.offset.y)
	}
}

// Allows sorting of characters by distance
type ByDistance []*Char

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].offsetDist < a[j].offsetDist }
