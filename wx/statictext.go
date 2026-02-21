// Package wx static text implementation
//
// This file implements the StaticText control, representing a non-editable text label.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "kwx_def.h"
#include "kwx_types.h"

// StaticText functions
extern void* wxStaticText_Create(void* parent, int id, void* label,
                                 int x, int y, int w, int h, int style);
*/
import "C"
import (
	"unsafe"
)

// StaticText represents a non-editable text label control.
type StaticText struct {
	BaseWindow
}

// NewStaticText creates a new static text (label) control.
//
// Parameters:
//   - parent: Parent window
//   - id: Window identifier (use ID_ANY for automatic)
//   - label: Label text
//   - pos: Label position (use DefaultPosition for default)
//   - size: Label size (use DefaultSize for default)
//   - style: Window style flags (use 0 for default)
//
// Example:
//
//	label := wx.NewStaticText(panel, wx.ID_ANY, "Enter your name:",
//	    wx.Point{X: 10, Y: 10}, wx.DefaultSize, 0)
func NewStaticText(parent Window, id int, label string, pos Point, size Size, style int) *StaticText {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxLabel := NewWxString(label)
	defer wxLabel.Free()

	ptr := C.wxStaticText_Create(
		parentPtr,
		C.int(id),
		wxLabel.Ptr(),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.int(style),
	)

	text := &StaticText{}
	text.SetPtr(ptr)
	return text
}
