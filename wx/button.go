// Package wx button implementation
//
// This file implements the Button control, representing a clickable button.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// Button functions
extern void* wxButton_Create(void* parent, int id, void* label,
                             int x, int y, int w, int h, int style);
extern void wxButton_SetDefault(void* self);
*/
import "C"
import (
	"unsafe"
)

// Button represents a clickable button control.
type Button struct {
	BaseWindow
}

// NewButton creates a new button control.
//
// Parameters:
//   - parent: Parent window
//   - id: Window identifier (use ID_ANY for automatic)
//   - label: Button label text
//   - pos: Button position (use DefaultPosition for default)
//   - size: Button size (use DefaultSize for default)
//   - style: Window style flags (use 0 for default)
//
// Example:
//
//	button := wx.NewButton(panel, wx.ID_ANY, "Click Me!",
//	    wx.Point{X: 10, Y: 10}, wx.DefaultSize, 0)
//	button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
//	    wx.MessageBox("Button clicked!", "Info", wx.OK, nil)
//	})
func NewButton(parent Window, id int, label string, pos Point, size Size, style int) *Button {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxLabel := NewWxString(label)
	defer wxLabel.Free()

	ptr := C.wxButton_Create(
		parentPtr,
		C.int(id),
		wxLabel.Ptr(),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.int(style),
	)

	button := &Button{}
	button.SetPtr(ptr)
	return button
}

// SetDefault sets this button as the default button for the parent window.
// The default button is activated when the user presses Enter.
//
// Example:
//
//	okButton := wx.NewButton(dialog, wx.ID_OK, "OK", ...)
//	okButton.SetDefault()
func (b *Button) SetDefault() {
	C.wxButton_SetDefault(b.ptr)
}
