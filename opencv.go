package main

import (
	"image"
	"image/png"
	"os"
	"path"
	"runtime"

	"github.com/lazywei/go-opencv/opencv"
)

func edgeCV(i image.Image, threshold int) image.Image {
	img := opencv.FromImage(i)

	w := img.Width()
	h := img.Height()

	// Create the output image
	cedge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 4)
	defer cedge.Release()

	// Convert to grayscale
	gray := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	edge := opencv.CreateImage(w, h, opencv.IPL_DEPTH_8U, 1)
	defer gray.Release()
	defer edge.Release()

	opencv.CvtColor(img, gray, opencv.CV_BGR2GRAY)

	opencv.Not(gray, edge)
	opencv.Canny(gray, edge, float64(threshold), float64(threshold*5), 3)
	opencv.Zero(cedge)
	opencv.Copy(img, cedge, edge)

	return cedge.ToImage()
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
