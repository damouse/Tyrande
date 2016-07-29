package main

// Gets updates from vision and updates characters
func modeling() {
	if !running {
		return
	}
}

type Character struct {
	*Line
}

// Modeling and detecting on-screen players
type Line struct {
	pixels []*Pix

	centerX, centerY       int // center
	maxX, maxY, minX, minY int
}

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

//
// Targeting
// Return the line whose center is closest to the screen center. If no lines passed, returns nil
func closestCenter(lines []*Line, centerX, centerY int) (ret *Line) {
	closest := 10000.0

	for _, l := range lines {
		dist := euclideanDistance(l.centerX, l.centerY, centerX, centerY)

		if dist < closest {
			closest = dist
			ret = l
		}

		// fmt.Println("Center: ", l.centerX, l.centerY, dist)
	}

	return
}
