package main

import (
	"fmt"
	"image"
	"image/color"
	"time"
)

var (
	EDGE_THRESHOLD float64 = 100 // Canny recognizer threshold

	SEPERATION_THRESHOLD    float64     = 0.6 // colors greater in distance than this are turned black
	SEPERATION_TARGETCOLOR1 color.NRGBA = color.NRGBA{255, 0, 0, 255}
	SEPERATION_TARGETCOLOR2 color.NRGBA = color.NRGBA{166, 64, 71, 255} // works for leftmost
	SEPERATION_TARGETCOLOR3 color.NRGBA = color.NRGBA{224, 84, 64, 255} // works for second to left
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
	i = p.add(accentColorDiffereenceGreyscale(i, SEPERATION_TARGETCOLOR1))

	// k := accentColorDiffereenceGreyscale(i, SEPERATION_TARGETCOLOR2)
	// p.results = append(p.results, k)

	// j := accentColorDiffereenceGreyscale(i, SEPERATION_TARGETCOLOR3)
	// p.results = append(p.results, j)

	// OpenCV Edge Detection
	// i = edgeCV(i)
	// p.results = append(p.results, i)

	// Gift Sobel
	// i = sobel(i)
	// p.results = append(p.results, i)

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
