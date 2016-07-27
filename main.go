package main

import (
	"image"
	"image/color"
	"runtime"
)

func runpipe() {
	p := NewPipeline()
	p.run(open("0.png"))
	p.save()
}

func testshop() {
	// p := open("lowsett.png")

	// g := gift.New(
	// 	gift.Contrast(1),
	// 	gift.Gamma(.75),
	// 	gift.ColorFunc(
	// 		func(r0, g0, b0, a0 float32) (r, g, b, a float32) {
	// 			r = r0 - (100.0 / 255.0) // invert the red channel
	// 			g = g0 - (60.0 / 255.0)  // shift the green channel by 0.1
	// 			b = b0                   // set the blue channel to 0
	// 			a = a0                   // preserve the alpha channel
	// 			return
	// 		},
	// 	),
	// )

	// dst := image.NewNRGBA(g.Bounds(p.Bounds()))
	// g.Draw(dst, p)
	// go save(dst, "11.png")

	// p = accentColorDiffereenceGreyscale(dst, color.NRGBA{188, 5, 18, 255}, 0.7)

	dst := open("11.png")

	b := dst.Bounds()
	n := image.NewGray(b)

	// For all pixels
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			c := dst.At(x, y)

			// get similiarity between this and target colors
			h := colorDistance(c, color.NRGBA{188, 5, 18, 255})

			// If the color is not similar set to black and continue
			if h > 0.8 {
				n.SetGray(x, y, color.Gray{0})
				continue
			}

			// For all neighboring pixels
			// n.SetGray(x, y, color.Gray{uint8(255 - h*255)})
		}
	}

	save(n, "12.png")
}

func testhunter() {
	// fmt.Printf("NumCPU: %d\n", runtime.NumCPU())

	allColors := []color.Color{
		color.NRGBA{219, 18, 29, 255},
		color.NRGBA{140, 31, 59, 255},
		color.NRGBA{182, 40, 59, 255},
		color.NRGBA{212, 128, 151, 255},
	}

	p := open("small.png")

	// w := NewWindow()
	// w.show(p)

	hunt(p, allColors, 0.2, 1)

	// w.show(i)
	// w.wait()

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// runpipe()
	// testshop()

	testhunter()
}
