package main

import "C"
import (
	"image"
	"os"
	"path"
	"runtime"

	"github.com/disintegration/gift"
	"github.com/lazywei/go-opencv/opencv"
)

// A set of operations makes up the pipeline
type Window struct {
	cvWindow *opencv.Window
	// recv     chan image.Image
	working bool
}

func NewWindow() *Window {
	return &Window{
		opencv.NewWindow("Tyrande", 1),
		// make(chan image.Image, 0),
		false,
	}
}

func (w *Window) queueShow(op *Cycle) {
	if !w.working {
		go w.show(op.mat.toImage())
	}
}

func (w *Window) show(i image.Image) {
	w.working = true

	g := gift.New(
		gift.Resize(i.Bounds().Max.X/2, i.Bounds().Max.Y/2, gift.LinearResampling),
	)

	dst := image.NewNRGBA(g.Bounds(i.Bounds()))
	g.Draw(dst, i)
	w.cvWindow.ShowImage(opencv.FromImage(dst))

	w.working = false
}

func (w *Window) refresh(i image.Image) {
	w.cvWindow.ShowImage(opencv.FromImage(i))
}

func (w *Window) wait() {
	for {
		key := opencv.WaitKey(20)
		if key == 112 {
			return
		}
	}
}

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
