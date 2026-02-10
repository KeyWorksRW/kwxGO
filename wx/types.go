package wx

import "unsafe"

// Point represents a 2D point with integer coordinates.
// In wxFFI calls, Point is decomposed into separate x and y parameters.
type Point struct {
	X, Y int
}

// Size represents a 2D size with integer width and height.
// In wxFFI calls, Size is decomposed into separate w and h parameters.
type Size struct {
	W, H int
}

// Rect represents a rectangle with position and size.
// In wxFFI calls, Rect is decomposed into x, y, w, h parameters.
type Rect struct {
	X, Y, W, H int
}

// DefaultPosition is the default position for window creation.
// wxWidgets will choose an appropriate position.
var DefaultPosition = Point{X: -1, Y: -1}

// DefaultSize is the default size for window creation.
// wxWidgets will choose an appropriate size.
var DefaultSize = Size{W: -1, H: -1}

// BaseWindow provides the common implementation for all window types.
// It wraps a pointer to the underlying wxWindow object.
type BaseWindow struct {
	ptr unsafe.Pointer
}

// Ptr returns the underlying wxWindow pointer.
// This is used when passing windows to wxFFI functions.
func (w *BaseWindow) Ptr() unsafe.Pointer {
	return w.ptr
}

// IsNil returns true if the window pointer is nil.
func (w *BaseWindow) IsNil() bool {
	return w.ptr == nil
}

// SetPtr sets the underlying pointer. Used internally by constructors.
func (w *BaseWindow) SetPtr(ptr unsafe.Pointer) {
	w.ptr = ptr
}
