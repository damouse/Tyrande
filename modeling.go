package main

import (
	"fmt"
	"sort"
)

type Char struct {
	center     Vec
	offset     Vec // Vec from this char to the center of the screen
	offsetDist float64
}

// Build a list of characters from a list of lines, sets Char data, and orders return based on closest
func buildChars(lines []*Line, center Vec) (ret []*Char) {
	for _, l := range lines {
		c := &Char{}

		c.center = Vec{l.avg.x, l.avg.y}
		c.offset = Vec{center.x - c.center.x, center.y - c.center.y}
		c.offsetDist = euclideanDistanceVec(center, c.center)

		ret = append(ret, c)
	}

	sort.Sort(ByDistance(ret))
	return
}

// Targeting
// Return the line whose center is closest to the screen center. If no lines passed, returns nil
func closestCenter(chars []*Char, center Vec) (ret *Char) {
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
