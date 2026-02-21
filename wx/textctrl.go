// Package wx text control implementation
//
// This file implements the TextCtrl control, representing a text input field.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "kwx_def.h"
#include "kwx_types.h"

// TextCtrl functions
extern void* wxTextCtrl_Create(void* parent, int id, void* value,
                               int x, int y, int w, int h, long style);
extern void* wxTextCtrl_GetValue(void* self);
extern void wxTextCtrl_SetValue(void* self, void* value);
extern void wxTextCtrl_AppendText(void* self, void* text);
extern void wxTextCtrl_Clear(void* self);
*/
import "C"
import (
	"unsafe"
)

// TextCtrl represents a text input control.
// It can be single-line or multi-line depending on style flags.
type TextCtrl struct {
	BaseWindow
}

// NewTextCtrl creates a new text control.
//
// Parameters:
//   - parent: Parent window
//   - id: Window identifier (use ID_ANY for automatic)
//   - value: Initial text value
//   - pos: Control position (use DefaultPosition for default)
//   - size: Control size (use DefaultSize for default)
//   - style: Window style flags (use 0 for default single-line, TE_MULTILINE for multi-line)
//
// Example:
//
//	// Single-line text input
//	input := wx.NewTextCtrl(panel, wx.ID_ANY, "",
//	    wx.Point{X: 10, Y: 40}, wx.Size{W: 200, H: -1}, 0)
//
//	// Multi-line text area
//	textArea := wx.NewTextCtrl(panel, wx.ID_ANY, "Initial text",
//	    wx.Point{X: 10, Y: 70}, wx.Size{W: 300, H: 100}, wx.TE_MULTILINE)
func NewTextCtrl(parent Window, id int, value string, pos Point, size Size, style int) *TextCtrl {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxValue := NewWxString(value)
	defer wxValue.Free()

	ptr := C.wxTextCtrl_Create(
		parentPtr,
		C.int(id),
		wxValue.Ptr(),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.long(style),
	)

	ctrl := &TextCtrl{}
	ctrl.SetPtr(ptr)
	return ctrl
}

// GetValue returns the current text in the control.
//
// Example:
//
//	text := textCtrl.GetValue()
//	fmt.Println("User entered:", text)
func (t *TextCtrl) GetValue() string {
	wxStr := C.wxTextCtrl_GetValue(t.ptr)
	return WxStringToGoAndFree(wxStr)
}

// SetValue sets the text in the control.
// This replaces any existing text.
//
// Example:
//
//	textCtrl.SetValue("New text")
func (t *TextCtrl) SetValue(value string) {
	wxValue := NewWxString(value)
	defer wxValue.Free()
	C.wxTextCtrl_SetValue(t.ptr, wxValue.Ptr())
}

// AppendText appends text to the end of the control.
// This is more efficient than GetValue + SetValue for adding text.
//
// Example:
//
//	logCtrl.AppendText("Log entry: Operation completed\n")
func (t *TextCtrl) AppendText(text string) {
	wxText := NewWxString(text)
	defer wxText.Free()
	C.wxTextCtrl_AppendText(t.ptr, wxText.Ptr())
}

// Clear removes all text from the control.
//
// Example:
//
//	textCtrl.Clear()
func (t *TextCtrl) Clear() {
	C.wxTextCtrl_Clear(t.ptr)
}
