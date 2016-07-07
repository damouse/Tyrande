package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"runtime"

	"github.com/lazywei/go-opencv/opencv"

	"github.com/disintegration/gift"
)

func runpipe() {
	p := NewPipeline()
	p.run(open("sample.png"))
	p.save()
}

func adjustments() {
	p := open("sample.png")
	// d := imaging.AdjustContrast(p, 50)
	// d = imaging.Sharpen(d, 50)

	// 1. Create a new GIFT filter list and add some filters:
	g := gift.New(
		// gift.UnsharpMask(1.0, 10.0, 10.0),
		// gift.Contrast(50),
		// gift.Saturation(50),
		//gift.ColorBalance(-50, 50, 50),
		gift.Convolution( // emboss
			[]float32{
				-1, -1, 0,
				-1, 1, 1,
				0, 1, 1,
			},
			false, false, false, 0.0,
		),
		// gift.Convolution( // edge detection
		// 	[]float32{
		// 		-1, -1, -1,
		// 		-1, 8, -1,
		// 		-1, -1, -1,
		// 	},
		// 	false, false, false, 0.0,
		// ),
		// gift.Sobel(),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image
	dst := image.NewRGBA(g.Bounds(p.Bounds()))

	// 3. Use Draw func to apply the filters to src and store the result in dst:
	g.Draw(dst, p)

	f, err := os.Create("./assets/adjusted.png")
	defer f.Close()
	err = png.Encode(f, dst)
	checkError(err)
}

func sandbox() {
	p := open("adjusted.png")

	b := p.Bounds()
	n := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pix := p.NRGBAAt(x, y)

			h := colorDistance(pix, targetColor1)
			n.Set(x, y, color.Gray{uint8(225 - h*255)})

			// c := colorful.Color{float64(pix.R) / 255.0, float64(pix.G) / 255.0, float64(pix.B) / 255.0}
			// _, h, _ := c.Hsv()
			// n.Set(x, y, color.Gray{uint8(h * 255)})
		}
	}

	f, err := os.Create("./assets/hue.png")
	defer f.Close()
	err = png.Encode(f, n)
	checkError(err)
}

func sand2() {
	_, currentfile, _, _ := runtime.Caller(0)
	filename := path.Join(path.Dir(currentfile), "./assets/hue.png")
	img := opencv.LoadImage(filename)
	defer img.Release()

	// Create the output image
	edge := opencv.CreateImage(img.Width(), img.Height(), opencv.IPL_DEPTH_8U, 1)

	opencv.Canny(img, edge, 200, 400, 6)

	ret := edge.ToImage()
	f, err := os.Create("./assets/sand.png")
	defer f.Close()
	err = png.Encode(f, ret)
	checkError(err)
}

func main() {
	adjustments()
	sandbox()
	// sand2()

	// stripImage()
	// edgy("hue.png")
	// runpipe()
}
