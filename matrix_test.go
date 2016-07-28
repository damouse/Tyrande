package main

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestMatrixSize(t *testing.T) {
	m := newTrackingMat(10, 10)
	Equal(t, len(m.arr), 100)
}

// func TestMatrixSet(t *testing.T) {
// 	m := NewTestingMatrix(3, 3)

// 	m.set(0, 0, 0)
// 	m.set(1, 1, 1)
// 	m.set(2, 2, 2)

// 	Equal(t, m.get(0, 0), 0)
// 	Equal(t, m.get(1, 1), 1)
// 	Equal(t, m.get(2, 2), 2)
// }

// func TestMatrixIter(t *testing.T) {
// 	m := NewTestingMatrix(3, 3)

// 	count := 0

// 	m.iter(func(x, y, z int) {
// 		count += 1
// 	})

// 	Equal(t, count, 9)
// }
