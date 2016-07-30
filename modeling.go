package main

import "time"

type Line struct {
	pixels []*Pix

	centerX, centerY       int // center
	maxX, maxY, minX, minY int
}

type Character struct {
	*Line
	offset Vector // vector from this char to the center of the screen
}

var lastUpdate time.Time = time.Now()

// Gets updates from vision and updates characters
func modeling() {
	lines := <-visionChan
	start := time.Now()

	lines = filterLines(lines)

	for _, l := range lines {
		l.process()
	}

	// Update the list of characters
	var chars []*Character

	for _, l := range lines {
		chars = append(chars, &Character{l, Vector{}})
	}

	bench("MOD", start)

	if len(characters) != len(chars) {
		// log("MOD %d chars", len(chars))
	}

	// Update shared store
	characterLock.Lock()
	characters = chars
	characterLock.Unlock()

	log("Update: %s", time.Since(lastUpdate))
	lastUpdate = time.Now()
}

// Calculate statistics about this line
func (l *Line) process() {
	var sumX, sumY int

	for _, p := range l.pixels {
		if p.x < l.minX {
			l.minX = p.x
		}

		if p.y < l.minY {
			l.minY = p.y
		}

		if p.x > l.maxX {
			l.maxX = p.x
		}

		if p.y > l.maxY {
			l.maxY = p.y
		}

		sumX += p.x
		sumY += p.y
	}

	l.centerX = sumX / len(l.pixels)
	l.centerY = sumY / len(l.pixels)
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

//
// Targeting
// Return the line whose center is closest to the screen center. If no lines passed, returns nil
func closestCenter(chars []*Character, center Vector) (ret *Character) {
	closest := 10000.0

	for _, char := range chars {
		l := char.Line
		// fmt.Println(center, ret)

		char.offset = Vector{center.x - char.centerX, center.y - char.centerY}
		dist := euclideanDistance(l.centerX, l.centerY, center.x, center.y)

		if dist < closest {
			closest = dist
			ret = char
		}
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
