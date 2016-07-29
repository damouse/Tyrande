package main

import (
	"image/color"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestMatrixSize(t *testing.T) {
	m := NewPixMatrix(10, 10)
	Equal(t, len(m.arr), 100)
}

func TestAdjacent(t *testing.T) {
	m := createMat(3, 3)
	n := m.adjacent(m.get(1, 1), 1)

	Equal(t, 8, len(n))
}

func TestAdjacentColor(t *testing.T) {
	m := createMat(3, 3)

	m.set(NewPix(0, 0, color.RGBA{255, 0, 0, 255}))
	m.set(NewPix(1, 1, color.RGBA{255, 0, 0, 255}))
	m.set(NewPix(1, 0, color.RGBA{255, 0, 0, 255}))

	n := m.adjacentSimilarColor(m.get(1, 1), m.get(1, 1), 1, 0.2)

	Equal(t, 2, len(n))
}

func TestAdjacentLarge(t *testing.T) {
	m := createMat(6, 6)
	n := m.adjacent(m.get(2, 2), 2)

	Equal(t, 24, len(n))
}

func TestAdjacentEdge(t *testing.T) {
	m := createMat(3, 3)
	n := m.adjacent(m.get(2, 2), 1)

	Equal(t, 3, len(n))
}

// Utility method
func createMat(w, h int) *PixMatrix {
	m := NewPixMatrix(w, h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			m.set(NewPix(x, y, color.Gray{}))
		}
	}

	return m
}
