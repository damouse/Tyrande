package main

import (
	"os"
	"path"
	"runtime"

	"github.com/lazywei/go-opencv/opencv"
)

// // A set of operations makes up the pipeline
// type Window struct {
// 	cvWindow *opencv.Window
// 	*Pipeline
// 	selectedImage int
// }

// func NewWindow() *Window {
// 	win := opencv.NewWindow("Tyrande")

// 	return &Window{
// 		cvWindow:      win,
// 		selectedImage: 0,
// 	}
// }

// func (w *Window) build(pipe *Pipeline) {
// 	w.Pipeline = pipe

// 	// Current image
// 	w.cvWindow.CreateTrackbar("Operation", 0, 1, func(pos int, param ...interface{}) {
// 		w.selectedImage = pos
// 		w.refresh()
// 	})

// 	// Edge threshold
// 	w.cvWindow.CreateTrackbar("edgeThreshold", 1, 100, func(pos int, param ...interface{}) {
// 		w.Pipeline.edgeThreshold = pos
// 		w.run(*w.get(0))
// 		w.refresh()
// 	})
// }

// func (w *Window) wait() {
// 	for {
// 		key := opencv.WaitKey(20)
// 		if key == 27 {
// 			os.Exit(0)
// 		}
// 	}
// }

// func (w *Window) show(i Image) {
// 	w.cvWindow.ShowImage(opencv.FromImage(i))
// }

// func (w *Window) run(i Image) {
// 	w.Pipeline.run(i)
// 	w.refresh()
// }

// func (w *Window) refresh() {
// 	img := w.get(w.selectedImage)

// 	if img != nil {
// 		w.cvWindow.ShowImage(opencv.FromImage(img))
// 	}
// }

// Used for testing stuff with opencv
func edgy(f string) {
	_, currentfile, _, _ := runtime.Caller(0)
	filename := path.Join(path.Dir(currentfile), "./assets/"+f)
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

		opencv.Not(gray, edge)
		opencv.Canny(gray, edge, float64(edge_thresh), float64(edge_thresh*5), 3)
		opencv.Zero(cedge)
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
