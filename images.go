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
func convertImage(i image.Image) *PixMatrix {
	b := i.Bounds()
	mat := NewPixMatrix(b.Max.X, b.Max.Y)

	// This works, but is wildly not concurrent.
	// Maybe move this to the screen processing area?

	// In the end it may just be easier to write this in C. Here's an example
	// http://stackoverflow.com/questions/6629798/whats-wrong-with-this-rgb-to-xyz-color-space-conversion-algorithm

	// sliceX := b.Max.X / CONVERTING_GOROUTINES
	// sliceY := b.Max.Y / CONVERTING_GOROUTINES

	// fmt.Printf("Bounds: %d sliceX: %d, sliceY: %d\n", b, sliceX, sliceY)

	// wg := &sync.WaitGroup{}
	// wg.Add(CONVERTING_GOROUTINES)

	// mtx := &sync.Mutex{}

	// for worker := 0; worker < CONVERTING_GOROUTINES; worker++ {
	// 	go func(n int) {

	// 		fmt.Printf("N: %d, yMin: %d yMax: %d xMin: %d xMax: %d\n", n, n*sliceY, (n+1)*sliceY, n*sliceX, (n+1)*sliceX)

	// 		for y := n * sliceY; y < (n+1)*sliceY; y++ {
	// 			for x := b.Min.X; x < b.Max.X; x++ {
	// 				mtx.Lock()
	// 				mat.set(x, y, NewPix(x, y, i.At(x, y)))
	// 				mtx.Unlock()
	// 			}
	// 		}

	// 		wg.Done()
	// 	}(worker)
	// }

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// mat.set(NewPix(x, y, i.At(x, y)))
			p := mat.get(x, y)
			p.Color = i.At(x, y)
			p.x = x
			p.y = y
		}
	}

	// wg.Wait()
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
	SWATCH = convertSwatches()
	return

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
	// save(ps, "edittedswatch.png")
	return
}

var TARGET_SWATCH = []color.Color{
	color.NRGBA{219, 18, 29, 255},
	color.NRGBA{140, 31, 59, 255},
	color.NRGBA{182, 40, 59, 255},
	color.NRGBA{212, 128, 151, 255},
}

func convertSwatches() (ret []*Pix) {
	for _, c := range TARGET_SWATCH {
		ret = append(ret, NewPix(0, 0, c))
	}

	return
}
