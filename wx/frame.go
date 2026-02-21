// Package wx frame implementation
//
// This file implements the Frame type, representing a top-level window with
// decorations like title bar, menu bar, and status bar.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "kwx_def.h"
#include "kwx_types.h"

// Frame functions
extern void* wxFrame_Create(void* parent, int id, void* title,
                            int x, int y, int w, int h, int style);
extern void wxFrame_SetStatusText(void* self, void* text, int number);
extern void* wxFrame_CreateStatusBar(void* self, int number, int style);
extern void* wxFrame_CreateToolBar(void* self, long style);
extern void* wxFrame_GetMenuBar(void* self);
extern void wxFrame_SetMenuBar(void* self, void* menuBar);
*/
import "C"
import (
	"unsafe"
)

// Frame represents a top-level window with decorations.
// Frames can have menu bars, toolbars, and status bars.
type Frame struct {
	BaseWindow
}

// NewFrame creates a new frame window.
//
// Parameters:
//   - parent: Parent window (usually nil for top-level frames)
//   - id: Window identifier (use ID_ANY for automatic)
//   - title: Frame title
//   - pos: Frame position (use DefaultPosition for default)
//   - size: Frame size (use DefaultSize for default)
//   - style: Window style flags (use DEFAULT_FRAME_STYLE for default)
//
// Example:
//
//	frame := wx.NewFrame(nil, wx.ID_ANY, "My Application",
//	    wx.DefaultPosition, wx.Size{W: 800, H: 600},
//	    wx.DEFAULT_FRAME_STYLE)
func NewFrame(parent Window, id int, title string, pos Point, size Size, style int) *Frame {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxTitle := NewWxString(title)
	defer wxTitle.Free()

	ptr := C.wxFrame_Create(
		parentPtr,
		C.int(id),
		wxTitle.Ptr(),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.int(style),
	)

	frame := &Frame{}
	frame.SetPtr(ptr)
	return frame
}

// SetStatusText sets the text in the status bar.
// The number parameter selects which status bar pane to update (0-based).
// Most frames have only one pane (number = 0).
func (f *Frame) SetStatusText(text string, number int) {
	wxText := NewWxString(text)
	defer wxText.Free()
	C.wxFrame_SetStatusText(f.ptr, wxText.Ptr(), C.int(number))
}

// CreateStatusBar creates a status bar at the bottom of the frame.
//
// Parameters:
//   - number: Number of status bar panes (usually 1)
//   - style: Status bar style flags (use 0 for default)
//
// Returns the status bar window pointer.
//
// Example:
//
//	frame.CreateStatusBar(1, 0)
//	frame.SetStatusText("Ready", 0)
func (f *Frame) CreateStatusBar(number, style int) unsafe.Pointer {
	return C.wxFrame_CreateStatusBar(f.ptr, C.int(number), C.int(style))
}

// CreateToolBar creates a toolbar at the top of the frame.
//
// Parameters:
//   - style: Toolbar style flags (use 0 for default)
//
// Returns the toolbar window pointer.
func (f *Frame) CreateToolBar(style int) unsafe.Pointer {
	return C.wxFrame_CreateToolBar(f.ptr, C.long(style))
}

// GetMenuBar returns the frame's menu bar pointer, or nil if none exists.
func (f *Frame) GetMenuBar() unsafe.Pointer {
	return C.wxFrame_GetMenuBar(f.ptr)
}

// SetMenuBar sets the frame's menu bar.
// The frame takes ownership of the menu bar.
func (f *Frame) SetMenuBar(menuBar unsafe.Pointer) {
	C.wxFrame_SetMenuBar(f.ptr, menuBar)
}

// Centre centers the frame on screen.
// This is a convenience method that uses SetSize with centering logic.
func (f *Frame) Centre() {
	// Get display size
	dispW, dispH := GetDisplaySize()
	// Get frame size
	w, h := f.GetSize()
	// Calculate center position
	x := (dispW - w) / 2
	y := (dispH - h) / 2
	// Move to center
	f.Move(x, y)
}

// Center is an alias for Centre (American spelling).
func (f *Frame) Center() {
	f.Centre()
}
