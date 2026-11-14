package wx

import "unsafe"

// Size is a wxSize value: a width and a height, in pixels.
type Size struct {
	Width  int
	Height int
}

// Point is a wxPoint value: a position in device coordinates.
type Point struct {
	X int
	Y int
}

// Rect is a wxRect value: the position and size of a rectangle.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// sizeFromPtr converts a heap-allocated wxSize returned by the C layer into a
// Size value and releases the native object.
func sizeFromPtr(ptr unsafe.Pointer) Size {
	if ptr == nil {
		return Size{}
	}
	size := Size{Width: cgo_wxSize_GetWidth(ptr), Height: cgo_wxSize_GetHeight(ptr)}
	cgo_wxSize_Delete(ptr)
	return size
}

// pointFromPtr converts a heap-allocated wxPoint returned by the C layer into a
// Point value and releases the native object.
func pointFromPtr(ptr unsafe.Pointer) Point {
	if ptr == nil {
		return Point{}
	}
	point := Point{X: cgo_wxPoint_GetX(ptr), Y: cgo_wxPoint_GetY(ptr)}
	cgo_wxPoint_Delete(ptr)
	return point
}

// rectFromPtr converts a heap-allocated wxRect returned by the C layer into a
// Rect value and releases the native object.
func rectFromPtr(ptr unsafe.Pointer) Rect {
	if ptr == nil {
		return Rect{}
	}
	rect := Rect{X: cgo_wxRect_GetX(ptr), Y: cgo_wxRect_GetY(ptr),
		Width: cgo_wxRect_GetWidth(ptr), Height: cgo_wxRect_GetHeight(ptr)}
	cgo_wxRect_Delete(ptr)
	return rect
}
