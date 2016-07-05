package main

import (
	"fmt"
	"image"
	"reflect"
	"unsafe"

	"github.com/lxn/win"
)

func ScreenRect() (image.Rectangle, error) {
	hDC := win.GetDC(0)
	if hDC == 0 {
		return image.Rectangle{}, fmt.Errorf("Could not Get primary display err:%d\n", win.GetLastError())
	}
	defer win.ReleaseDC(0, hDC)
	x := win.GetDeviceCaps(hDC, win.HORZRES)
	y := win.GetDeviceCaps(hDC, win.VERTRES)
	return image.Rect(0, 0, int(x), int(y)), nil
}

func CaptureScreen() (*image.NRGBA, error) {
	r, e := ScreenRect()
	if e != nil {
		return nil, e
	}
	return CaptureRect(r)
}

func CaptureRect(rect image.Rectangle) (*image.NRGBA, error) {
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

	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	img := &image.NRGBA{imageBytes, 4 * x, image.Rect(0, 0, x, y)}
	return img, nil
}
