package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/disintegration/gift"
)

// Save this image inside the assets folder with the given name
func save(img image.Image, name string) {
	f, err := os.Create("./assets/" + name)
	checkError(err)
	defer f.Close()

	err = png.Encode(f, img)
	checkError(err)
}

func open(path string) image.Image {
	f, err := os.Open("./assets/" + path)
	checkError(err)
	defer f.Close()

	img, err := png.Decode(f)
	checkError(err)

	return img
}

// Convert an image to a Pix matrix
func convertImage(i image.Image) *TrackingMat {
	b := i.Bounds()
	mat := newTrackingMat(b.Max.X, b.Max.Y)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			mat.set(x, y, NewPix(x, y, i.At(x, y)))
		}
	}

	return mat
}

func iter(i image.Image, fn func(x int, y int, pixel color.Color)) {
	b := i.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fn(x, y, i.At(x, y))
		}
	}
}

func photoshop(i image.Image) image.Image {
	g := gift.New(
	// gift.Hue(45),
	// gift.Contrast(1),
	// gift.Saturation(2),
	// gift.Gamma(0.75),
	// gift.UnsharpMask(12.0, 30.0, 20.0),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewNRGBA(g.Bounds(i.Bounds()))

	g.Draw(dst, i)
	return dst
}

func loadSwatch() (result []color.Color) {
	var ret []color.Color

	i := open("swatch.png")

	iter(i, func(x, y int, c color.Color) {
		r, g, b, a := c.RGBA()

		if a == 0 {
			return
		}

		for _, c := range ret {
			er, eg, eb, _ := c.RGBA()

			if er == r && eg == g && eb == b {
				return
			}
		}

		ret = append(ret, c)
	})

	img := image.NewNRGBA(image.Rect(0, 0, 1, len(ret)))

	for i, c := range ret {
		r, g, b, a := c.RGBA()
		img.SetNRGBA(0, i, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	ps := photoshop(img)

	iter(ps, func(x, y int, c color.Color) {
		result = append(result, c)
	})

	fmt.Printf("Loaded %d colors\n", len(ret))
	save(ps, "edittedswatch.png")
	return
}

func convertSwatches() (ret []*Pix) {
	for _, c := range TARGET_SWATCH {
		ret = append(ret, NewPix(0, 0, c))
	}

	return
}
