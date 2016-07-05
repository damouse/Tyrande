package main

import (
	"fmt"
	"image"
	"image/color"
)

// The border seems to be a burnt orange-ish color
// It looks like its always about one pixel wide
// Pure: (233, 88, 61)

// Remove everything in the image except outlines
func stripImage() {
	i := open("sample.png")

	fmt.Println("Bounds: ", i.Bounds())
	n := Image{image.NewNRGBA(i.Bounds())}

	// targetColor := color.NRGBA{224, 84, 64, 255} // works for second to left
	targetColor := color.NRGBA{166, 64, 71, 255} // works for leftmost

	// Very rough iterator for the images
	b := i.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pix := i.NRGBAAt(x, y)

			distance := colorDistance(pix, targetColor)
			if distance <= 0.2 {
				// n.SetNRGBA(x, y, color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255})
				n.SetNRGBA(x, y, pix)
			} else {
				n.SetNRGBA(x, y, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			}

			// fmt.Printf("Distance: %f\n", )

			if pix == targetColor {
				fmt.Println("Target found: ", distance)

				// n.SetNRGBA(x, y, pix)
			}
		}
	}

	n.save("out.png")
}

func main() {
	stripImage()
}
