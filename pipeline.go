package main

import (
	"fmt"
	"time"
)

type Pipeline struct {
	results []Image

	edgeThreshold int
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		results:       []Image{},
		edgeThreshold: 100,
	}
}

func (p *Pipeline) run(img Image) {
	start := time.Now()

	p.results = []Image{}
	i := img
	p.results = append(p.results, i)

	// Initial image touchup
	i = accentColorDifference(i)
	p.results = append(p.results, i)

	// Difference contrasting
	// i = edgeCV(i, p.edgeThreshold)
	// p.results = append(p.results, i)

	// Edge
	i = edgeCV(i, p.edgeThreshold)
	p.results = append(p.results, i)

	fmt.Printf("Pipeline bench: %s", time.Since(start))
}

// Saves all the files with unique names for inspection
func (p *Pipeline) save() {
	for num, i := range p.results {
		i.save(fmt.Sprintf("%d.png", num))
	}
}

func (p *Pipeline) get(i int) *Image {
	if i >= len(p.results) {
		return nil
	}

	return &p.results[i]
}
