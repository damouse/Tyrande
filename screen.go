package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"reflect"
	"unsafe"

	"github.com/lxn/win"
)

func capture() *PixMatrix {
	var p image.Image

	if DEBUG_STATIC {
		p = imageStatic
	} else if targeting {
		p = CaptureLeftNarrow(0.3, 0.3)
	} else {
		// Testing direct mat
		return CaptureMat(LEFT_SCREEN_DIM)
		// p = CaptureLeft()
	}

	return convertImage(p)
}

func ScreenRect() (image.Rectangle, error) {
	hDC := win.GetDC(0)
	if hDC == 0 {
		return image.Rectangle{}, fmt.Errorf("Could not Get primary display err:%d\n", win.GetLastError())
	}
	defer win.ReleaseDC(0, hDC)
	x := win.GetDeviceCaps(hDC, win.HORZRES)
	y := win.GetDeviceCaps(hDC, win.VERTRES)

	fmt.Println(image.Rect(0, 0, int(x), int(y)))

	return image.Rect(0, 0, int(x), int(y)), nil
}

func CaptureScreen() (image.Image, error) {
	r, e := ScreenRect()
	if e != nil {
		return nil, e
	}
	return CaptureRect(r)
}

func CaptureLeft() image.Image {
	i, _ := CaptureRect(LEFT_SCREEN_DIM)
	return i
}

// Same as above, but captures only a fraction of the screen
func CaptureLeftNarrow(fracX, fracY float64) image.Image {
	dx := LEFT_SCREEN_DIM.Max.X - LEFT_SCREEN_DIM.Min.X
	dy := LEFT_SCREEN_DIM.Max.Y - LEFT_SCREEN_DIM.Min.Y

	offx := int(float64(dx) * fracX * 0.5)
	offy := int(float64(dy) * fracY * 0.5)

	newRect := image.Rect(
		LEFT_SCREEN_DIM.Min.X+offx,
		LEFT_SCREEN_DIM.Min.Y+offy,
		LEFT_SCREEN_DIM.Max.X-offx,
		LEFT_SCREEN_DIM.Max.Y-offy,
	)

	// fmt.Println("Rect: ", newRect)

	i, _ := CaptureRect(newRect)
	return i
}

func CaptureRect(rect image.Rectangle) (image.Image, error) {
	hDC := win.GetDC(0)
	if hDC == 0 {
		return nil, fmt.Errorf("Could not Get primary display err:%d.\n", win.GetLastError())
	}
	defer win.ReleaseDC(0, hDC)

	m_hDC := win.CreateCompatibleDC(hDC)
	if m_hDC == 0 {
		return nil, fmt.Errorf("Could not Create Compatible DC err:%d.\n", win.GetLastError())
	}
	defer win.DeleteDC(m_hDC)

	x, y := rect.Dx(), rect.Dy()

	bt := win.BITMAPINFO{}
	bt.BmiHeader.BiSize = uint32(reflect.TypeOf(bt.BmiHeader).Size())
	bt.BmiHeader.BiWidth = int32(x)
	bt.BmiHeader.BiHeight = int32(-y)
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = win.BI_RGB

	ptr := unsafe.Pointer(uintptr(0))

	m_hBmp := win.CreateDIBSection(m_hDC, &bt.BmiHeader, win.DIB_RGB_COLORS, &ptr, 0, 0)
	if m_hBmp == 0 {
		return nil, fmt.Errorf("Could not Create DIB Section err:%d.\n", win.GetLastError())
	}

	// if m_hBmp == win.InvalidParameter {
	// 	return nil, fmt.Errorf("One or more of the input parameters is invalid while calling CreateDIBSection.\n")
	// }

	defer win.DeleteObject(win.HGDIOBJ(m_hBmp))

	obj := win.SelectObject(m_hDC, win.HGDIOBJ(m_hBmp))

	if obj == 0 {
		return nil, fmt.Errorf("error occurred and the selected object is not a region err:%d.\n", win.GetLastError())
	}

	if obj == 0xffffffff { //GDI_ERROR
		return nil, fmt.Errorf("GDI_ERROR while calling SelectObject err:%d.\n", win.GetLastError())
	}

	defer win.DeleteObject(obj)

	win.BitBlt(m_hDC, 0, 0, int32(x), int32(y), hDC, int32(rect.Min.X), int32(rect.Min.Y), win.SRCCOPY)

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = x * y * 4
	hdrp.Cap = x * y * 4

	imageBytes := make([]byte, len(slice))

	// Switching this to a matrix will save some time, but not a whole lot

	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	img := &image.NRGBA{imageBytes, 4 * x, image.Rect(0, 0, x, y)}
	return img, nil
}

// Capture a matrix directly instead of an image
func CaptureMat(rect image.Rectangle) *PixMatrix {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)

	m_hDC := win.CreateCompatibleDC(hDC)
	defer win.DeleteDC(m_hDC)

	x, y := rect.Dx(), rect.Dy()

	bt := win.BITMAPINFO{}
	bt.BmiHeader.BiSize = uint32(reflect.TypeOf(bt.BmiHeader).Size())
	bt.BmiHeader.BiWidth = int32(x)
	bt.BmiHeader.BiHeight = int32(-y)
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = win.BI_RGB

	ptr := unsafe.Pointer(uintptr(0))

	m_hBmp := win.CreateDIBSection(m_hDC, &bt.BmiHeader, win.DIB_RGB_COLORS, &ptr, 0, 0)
	defer win.DeleteObject(win.HGDIOBJ(m_hBmp))

	obj := win.SelectObject(m_hDC, win.HGDIOBJ(m_hBmp))
	defer win.DeleteObject(obj)

	win.BitBlt(m_hDC, 0, 0, int32(x), int32(y), hDC, int32(rect.Min.X), int32(rect.Min.Y), win.SRCCOPY)

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = x * y * 4
	hdrp.Cap = x * y * 4

	imageBytes := make([]byte, len(slice))

	ret := NewPixMatrix(x, y)
	// Switching this to a matrix will save some time, but not a whole lot

	//(x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	// i = (y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4
	// i = y * Stride + x * 4

	for i := 0; i < len(imageBytes); i += 4 {
		pixIndex := i / 4
		// imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
		p := &ret.arr[pixIndex]
		p.Color = color.NRGBA{uint8(slice[i+2]), uint8(slice[i+1]), uint8(slice[i]), 255}

		p.x = int(math.Mod(float64(pixIndex), float64(ret.w)))
		p.y = pixIndex / ret.w
	}

	// img := &image.NRGBA{imageBytes, 4 * x, image.Rect(0, 0, x, y)}
	return ret
}
