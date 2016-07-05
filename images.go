package main

import (
	"image"
	"image/png"
	"os"
)

type Image struct {
	data *image.NRGBA
}

// Save this image inside the assets folder with the given name
func (i *Image) save(name string) {
	f, err := os.Create("./assets/" + name)
	checkError(err)
	defer f.Close()

	err = png.Encode(f, i.data)
	checkError(err)
}

// func open(path string) Image {
// 	f, err := os.Open("./assets/sample.png")
// 	checkError(err)
// 	defer f.Close()

// 	img, t, err := image.Decode(bufio.NewReader(f))
// 	i := img.(*image.NRGBA)

// 	fmt.Println("type of image:", t)
// 	checkError(err)

// 	ret := Image{data: i}
// 	return ret
// }

func open(path string) Image {
	f, err := os.Open("./assets/sample.png")
	checkError(err)
	defer f.Close()

	img, err := png.Decode(f)
	checkError(err)

	return Image{data: img.(*image.NRGBA)}
}
