package main

import "fmt"

// The border seems to be a burnt orange-ish color
// It looks like its always about one pixel wide
// Pure: (233, 88, 61)

// Remove everything in the image except outlines
func stripImage() {
	i := open("sample.png")

	fmt.Println("Bounds: ", i.data.Bounds())

	i.save("out.png")
}

func main() {
	stripImage()
}
