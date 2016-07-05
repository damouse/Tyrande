package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/lxn/win"
	"github.com/vova616/screenshot"
)

// The border seems to be a burnt orange-ish color
// It looks like its always about one pixel wide
// Pure: (233, 88, 61)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func capAndSave() {
	img, err := screenshot.CaptureScreen()

	f, err := os.Create("./ss.png")
	checkError(err)

	err = png.Encode(f, img)
	checkError(err)

	f.Close()
}

func openAndProcess() {
	f, err := os.Open("./sample.png")
	checkError(err)
	defer f.Close()

	img, _, err := image.Decode(bufio.NewReader(f))
	checkError(err)

	fmt.Println("Image opened")

	f, err = os.Create("./out.png")
	checkError(err)

	err = png.Encode(f, img)
	checkError(err)

	f.Close()
	fmt.Println("Results saved")
}

func oldMain() {
	start := time.Now()
	iterations := 20

	for i := 0; i < iterations; i++ {
		_, err := screenshot.CaptureScreen()
		checkError(err)
	}

	fmt.Printf("%d shots took %s", iterations, time.Since(start))

}

func windowsAPI() {
	p := win.POINT{}
	win.GetCursorPos(&p)

	fmt.Printf("Current position: %v\n", p)
}

func main() {
	// oldMain()
	// capAndSave()
	// openAndProcess()
	windowsAPI()
}
