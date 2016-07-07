package main

func runpipe() {
	p := NewPipeline()
	p.run(open("0.png"))
	p.save()
}

func testshop() {
	p := open("1.png")

	// p = transformRGB(p, func(x int, y int, p color.Color) color.Color {
	// 	pix := p.(color.NRGBA)

	// 	c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}

	// 	// tweaking contrast and hue
	// 	h, s, v := c.Hsv()
	// 	return color.NRGBA{R: uint8(255 * h), G: uint8(255 * h), B: uint8(255 * h), A: 255}
	// })

	save(p, "11.png")
}

func main() {
	runpipe()

	// testshop()
}
