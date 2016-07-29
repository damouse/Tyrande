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
	i, _ := CaptureRect(image.Rect(0, 0, 2100, 1440))
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

	// ret := NewPixMatrix(x, y)

	// Um... this cant possibly work. How is it working? Has this version been tested?
	// Have to do w/h calculation
	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
		// p := &ret.arr[i/4]
		// p.Color = color.RGBA{uint8(slice[i+2]), uint8(slice[i+1]), uint8(slice[i]), uint8(slice[i+3])}
	}

	img := &image.NRGBA{imageBytes, 4 * x, image.Rect(0, 0, x, y)}
	// return ret, nil
	return img, nil
}

// Golang dx9 wrapper: https://github.com/gonutz/d3d9
// This will most likely work, but may be slow.
// The dx9 library may help?
// func overlay() {
// 	hDC := win.GetDC(0)
// 	if hDC == 0 {
// 		fmt.Printf("Could not Get primary display err:%d.\n", win.GetLastError())
// 	}
// 	defer win.ReleaseDC(0, hDC)

// 	m_hDC := win.CreateCompatibleDC(hDC)
// 	if m_hDC == 0 {
// 		fmt.Printf("Could not Create Compatible DC err:%d.\n", win.GetLastError())
// 	}
// 	defer win.DeleteDC(m_hDC)

// 	rect, e := ScreenRect()
// 	checkError(e)

// 	hwnd := win.GetConsoleWindow()
// 	// rect := &win.RECT{}
// 	// win.GetClientRect(hwnd, rect)

// 	x, y := rect.Dx(), rect.Dy()

// 	bt := win.BITMAPINFO{}
// 	bt.BmiHeader.BiSize = uint32(reflect.TypeOf(bt.BmiHeader).Size())
// 	bt.BmiHeader.BiWidth = int32(x)
// 	bt.BmiHeader.BiHeight = int32(-y)
// 	bt.BmiHeader.BiPlanes = 1
// 	bt.BmiHeader.BiBitCount = 32
// 	bt.BmiHeader.BiCompression = win.BI_RGB

// 	ptr := unsafe.Pointer(uintptr(0))

// 	m_hBmp := win.CreateDIBSection(m_hDC, &bt.BmiHeader, win.DIB_RGB_COLORS, &ptr, 0, 0)
// 	if m_hBmp == 0 {
// 		fmt.Printf("Could not Create DIB Section err:%d.\n", win.GetLastError())
// 	}

// 	defer win.DeleteObject(win.HGDIOBJ(m_hBmp))

// 	obj := win.SelectObject(m_hDC, win.HGDIOBJ(m_hBmp))
// 	if obj == 0 {
// 		fmt.Printf("error occurred and the selected object is not a region err:%d.\n", win.GetLastError())
// 	}
// 	if obj == 0xffffffff { //GDI_ERROR
// 		fmt.Printf("GDI_ERROR while calling SelectObject err:%d.\n", win.GetLastError())
// 	}
// 	defer win.DeleteObject(obj)

// 	win.CreateBrushIndirect(&win.LOGBRUSH{win.BS_SOLID, win.DIB_RGB_COLORS, win.HS_BDIAGONAL})
// 	win.fill
// }

// void DrawRectangleOnTransparent(HWND hWnd, const RECT& rc)
// {
//     HDC hDC = GetDC(hWnd);
//     if (hDC)
//     {
//         RECT rcClient;
//         GetClientRect(hWnd, &rcClient);

//         BITMAPINFO bmi = { 0 };
//         bmi.bmiHeader.biSize = sizeof(bmi.bmiHeader);
//         bmi.bmiHeader.biBitCount = 32;
//         bmi.bmiHeader.biWidth = rcClient.right;
//         bmi.bmiHeader.biHeight = -rcClient.bottom;

//         LPVOID pBits;
//         HBITMAP hBmpSource = CreateDIBSection(hDC, &bmi, DIB_RGB_COLORS, &pBits, 0, 0);
//         if (hBmpSource)
//         {
//             HDC hDCSource = CreateCompatibleDC(hDC);
//             if (hDCSource)
//             {
//                 // fill the background in red
//                 HGDIOBJ hOldBmp = SelectObject(hDCSource, hBmpSource);
//                 HBRUSH hBsh = CreateSolidBrush(RGB(0,0,255));
//                 FillRect(hDCSource, &rcClient, hBsh);
//                 DeleteObject(hBsh);

//                 // draw the rectangle in black
//                 HGDIOBJ hOldBsh = SelectObject(hDCSource, GetStockObject(NULL_BRUSH));
//                 HGDIOBJ hOldPen = SelectObject(hDCSource, CreatePen(PS_SOLID, 2, RGB(0,0,0)));
//                 Rectangle(hDCSource, rc.left, rc.top, rc.right, rc.bottom);
//                 DeleteObject(SelectObject(hDCSource, hOldPen));
//                 SelectObject(hDCSource, hOldBsh);

//                 GdiFlush();

//                 // fix up the alpha channel
//                 DWORD* pPixel = reinterpret_cast<DWORD*>(pBits);
//                 for (int y = 0; y < rcClient.bottom; y++)
//                 {
//                     for (int x = 0; x < rcClient.right; x++, pPixel++)
//                     {
//                         if ((*pPixel & 0x00ff0000) == 0x00ff0000)
//                             *pPixel |= 0x01000000; // transparent
//                         else
//                             *pPixel |= 0xff000000; // solid
//                     }
//                 }

//                 // Update the layered window
//                 POINT pt = { 0 };
//                 BLENDFUNCTION bf = { AC_SRC_OVER, 0, 255, AC_SRC_ALPHA };
//                 UpdateLayeredWindow(hWnd, hDC, NULL, NULL, hDCSource, &pt, 0, &bf, ULW_ALPHA);

//                 SelectObject(hDCSource, hOldBmp);
//                 DeleteDC(hDCSource);
//             }
//             DeleteObject(hBmpSource);
//         }
//         ReleaseDC(hWnd, hDC);
//     }
// }
