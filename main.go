package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"runtime"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/lucasb-eyer/go-colorful"
)

// The border seems to be a burnt orange-ish color
// It looks like its always about one pixel wide
// Pure: (233, 88, 61)

func colorfulDistance(a color.NRGBA, c color.NRGBA) float64 {
	c1 := colorful.Color{float64(a.R) / 255.0, float64(a.G) / 255.0, float64(a.B) / 255.0}
	c2 := colorful.Color{float64(c.R) / 255.0, float64(c.G) / 255.0, float64(c.B) / 255.0}

	// Luv seems quite good

	return c1.DistanceCIE76(c2)

	// c := colorful.Color{0.313725, 0.478431, 0.721569}
	// c, err := colorful.Hex("#517AB8")
	// if err != nil{
	//     log.Fatal(err)
	// }
	// c = colorful.Hsv(216.0, 0.56, 0.722)
	// c = colorful.Xyz(0.189165, 0.190837, 0.480248)
	// c = colorful.Xyy(0.219895, 0.221839, 0.190837)
	// c = colorful.Lab(0.507850, 0.040585,-0.370945)
	// c = colorful.Luv(0.507849,-0.194172,-0.567924)
	// c = colorful.Hcl(276.2440, 0.373160, 0.507849)
	// fmt.Printf("RGB values: %v, %v, %v", c.R, c.G, c.B)
}

// Does a pretty good job making the outlines stand out
func hsvConversion(a color.NRGBA) color.NRGBA {
	c := colorful.Color{float64(a.R) / 255.0, float64(a.G) / 255.0, float64(a.B) / 255.0}
	r, g, b := c.Luv()

	return color.NRGBA{R: uint8(225 - 255*r), G: uint8(225 - 255*g), B: uint8(225 - 255*b), A: 255}
}

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

			// a very good distance measurement
			distance := colorfulDistance(pix, targetColor)
			newColor := color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255}

			//newColor := hsvConversion(pix)
			// newColor := color.NRGBA{R: uint8(225 - 255*distance), G: uint8(225 - 255*distance), B: uint8(225 - 255*distance), A: 255}

			n.SetNRGBA(x, y, newColor)
		}
	}

	n.save("out.png")
}

func edgy() {
	_, currentfile, _, _ := runtime.Caller(0)
	filename := path.Join(path.Dir(currentfile), "./assets/out.png")
	image := opencv.LoadImage(filename)
	defer image.Release()

	w := image.Width()
	h := image.Height()

	// Create the output image
	cedge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 3)
	defer cedge.Release()

	// Convert to grayscale
	gray := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	edge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	defer gray.Release()
	defer edge.Release()

	opencv.CvtColor(image, gray, opencv.CV_BGR2GRAY)

	win := opencv.NewWindow("Edge")
	defer win.Destroy()

	win.CreateTrackbar("Thresh", 1, 100, func(pos int, param ...interface{}) {
		edge_thresh := pos

		opencv.Smooth(gray, edge, opencv.CV_BLUR, 3, 3, 0, 0)
		opencv.Not(gray, edge)

		// Run the edge detector on grayscale
		opencv.Canny(gray, edge, float64(edge_thresh), float64(edge_thresh*3), 3)

		opencv.Zero(cedge)

		// copy edge points
		opencv.Copy(image, cedge, edge)

		win.ShowImage(cedge)
	})

	win.ShowImage(image)

	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			os.Exit(0)
		}
	}

	os.Exit(0)
}

func main() {
	// stripImage()
	edgy()
}
