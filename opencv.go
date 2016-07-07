package main

import (
	"image"

	"github.com/lazywei/go-opencv/opencv"
)

func edgeCV(i Image, threshold int) Image {
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

	ret := cedge.ToImage()
	return Image{ret.(*image.NRGBA)}
}
