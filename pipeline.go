package main

import (
	"fmt"
	"image"
	"time"
)

type Pipeline struct {
	results []image.Image
}

func NewPipeline() *Pipeline {
	return &Pipeline{[]image.Image{}}
}

func (p *Pipeline) run(img image.Image) {
	start := time.Now()
	p.results = []image.Image{}
	i := p.add(img)

	// Basic photo balancing and editing
	i = p.add(photoshop(i))

	// Hue/Sat
	// i = p.add(seperateHue(i))

	// Pick out the right color
	// i = p.add(accentColorDifference(i))

	// Check against all seperation colors
	// i = p.add(accentColorDiffereenceGreyscaleAggregate(i))

	// Pick the right color and drop it in a greyscale
	// i = p.add(accentColorDiffereenceGreyscale(i, SEPERATION_TARGETCOLOR3, SEPERATION_THRESHOLD))

	// Local maxima
	// i = p.add(localmax(i))

	// OpenCV Edge Detection
	// i = p.add(edgeCV(i))

	// Gift Sobel
	// i = p.add(sobel(i))

	fmt.Printf("Pipeline bench: %s", time.Since(start))
}

func (p *Pipeline) add(i image.Image) image.Image {
	p.results = append(p.results, i)
	return i
}

// Saves all the files with unique names for inspection
func (p *Pipeline) save() {
	for num, i := range p.results {
		save(i, fmt.Sprintf("%d.png", num))
	}
}

func (p *Pipeline) get(i int) *image.Image {
	if i >= len(p.results) {
		return nil
	}

	return &p.results[i]
}
