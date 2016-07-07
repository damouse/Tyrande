package main

import (
	"fmt"
	"image"
	"time"
)

var (
	EDGE_THRESHOLD int = 100
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
	i := img
	p.results = append(p.results, i)

	// Basic photo balancing and editing
	i = photoshop(i)
	p.results = append(p.results, i)

	// Pick out the right color
	i = accentColorDifference(i)
	p.results = append(p.results, i)

	// Difference contrasting
	// i = edgeCV(i, p.edgeThreshold)
	// p.results = append(p.results, i)

	// Edge
	i = edgeCV(i, EDGE_THRESHOLD)
	p.results = append(p.results, i)

	fmt.Printf("Pipeline bench: %s", time.Since(start))
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
