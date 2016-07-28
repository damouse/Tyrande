package main

// import . "github.com/stretchr/testify/assert"

// // Subclass
// type IntMatrix struct{ *Matrix }

// func NewTestingMatrix(w, h int) *IntMatrix {
// 	return &IntMatrix{NewMatrix(w, h)}
// }

// func (m *IntMatrix) get(x, y int) int {
// 	return m.Matrix._get(x, y).(int)
// }

// func (m *IntMatrix) iter(fn func(x int, y int, obj int)) {
// 	for y := 0; y < m.Matrix.h; y++ {
// 		for x := 0; x < m.Matrix.w; x++ {
// 			fn(x, y, m.get(x, y))
// 		}
// 	}
// }

// func TestMatrixSize(t *testing.T) {
// 	m := NewTestingMatrix(10, 10)
// 	Equal(t, len(m.Matrix.arr), 100)
// }

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
