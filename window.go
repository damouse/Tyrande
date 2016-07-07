package main

import (
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

// A set of operations makes up the pipeline
type Window struct {
	cvWindow *opencv.Window
	*Pipeline
	selectedImage int
}

func NewWindow() *Window {
	win := opencv.NewWindow("Tyrande")

	return &Window{
		cvWindow:      win,
		selectedImage: 0,
	}
}

func (w *Window) build(pipe *Pipeline) {
	w.Pipeline = pipe

	// Current image
	w.cvWindow.CreateTrackbar("Operation", 0, 1, func(pos int, param ...interface{}) {
		w.selectedImage = pos
		w.refresh()
	})

	// Edge threshold
	w.cvWindow.CreateTrackbar("edgeThreshold", 1, 100, func(pos int, param ...interface{}) {
		w.Pipeline.edgeThreshold = pos
		w.run(*w.get(0))
		w.refresh()
	})
}

func (w *Window) wait() {
	for {
		key := opencv.WaitKey(20)
		if key == 27 {
			os.Exit(0)
		}
	}
}

func (w *Window) show(i Image) {
	w.cvWindow.ShowImage(opencv.FromImage(i))
}

func (w *Window) run(i Image) {
	w.Pipeline.run(i)
	w.refresh()
}

func (w *Window) refresh() {
	img := w.get(w.selectedImage)

	if img != nil {
		w.cvWindow.ShowImage(opencv.FromImage(img))
	}
}
