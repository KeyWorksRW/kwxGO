// Package wx checkbox implementation
//
// This file implements the CheckBox control, representing a checkbox.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "kwx_def.h"
#include "kwx_types.h"

// CheckBox functions
extern void* wxCheckBox_Create(void* parent, int id, void* label,
                               int x, int y, int w, int h, int style);
extern int wxCheckBox_GetValue(void* self);
extern void wxCheckBox_SetValue(void* self, int state);
extern int wxCheckBox_Is3State(void* self);
extern int wxCheckBox_Get3StateValue(void* self);
extern void wxCheckBox_Set3StateValue(void* self, int state);
*/
import "C"
import (
	"unsafe"
)

// CheckBox represents a checkbox control.
// Checkboxes can be 2-state (checked/unchecked) or 3-state (checked/unchecked/undetermined).
type CheckBox struct {
	BaseWindow
}

// NewCheckBox creates a new checkbox control.
//
// Parameters:
//   - parent: Parent window
//   - id: Window identifier (use ID_ANY for automatic)
//   - label: Checkbox label text
//   - pos: Checkbox position (use DefaultPosition for default)
//   - size: Checkbox size (use DefaultSize for default)
//   - style: Window style flags (use 0 for default, CHK_3STATE for 3-state checkbox)
//
// Example:
//
//	// 2-state checkbox
//	checkbox := wx.NewCheckBox(panel, wx.ID_ANY, "Enable feature",
//	    wx.Point{X: 10, Y: 10}, wx.DefaultSize, 0)
//	checkbox.SetValue(true)
//
//	// 3-state checkbox
//	checkbox3 := wx.NewCheckBox(panel, wx.ID_ANY, "Optional setting",
//	    wx.Point{X: 10, Y: 40}, wx.DefaultSize, wx.CHK_3STATE)
func NewCheckBox(parent Window, id int, label string, pos Point, size Size, style int) *CheckBox {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxLabel := NewWxString(label)
	defer wxLabel.Free()

	ptr := C.wxCheckBox_Create(
		parentPtr,
		C.int(id),
		wxLabel.Ptr(),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.int(style),
	)

	checkbox := &CheckBox{}
	checkbox.SetPtr(ptr)
	return checkbox
}

// GetValue returns true if the checkbox is checked, false otherwise.
// For 3-state checkboxes, use Get3StateValue instead.
//
// Example:
//
//	if checkbox.GetValue() {
//	    fmt.Println("Feature enabled")
//	}
func (c *CheckBox) GetValue() bool {
	return C.wxCheckBox_GetValue(c.ptr) != 0
}

// SetValue sets the checkbox state.
// For 3-state checkboxes, use Set3StateValue instead.
//
// Example:
//
//	checkbox.SetValue(true)  // Check the box
func (c *CheckBox) SetValue(state bool) {
	val := C.int(0)
	if state {
		val = 1
	}
	C.wxCheckBox_SetValue(c.ptr, val)
}

// Is3State returns true if this is a 3-state checkbox.
func (c *CheckBox) Is3State() bool {
	return C.wxCheckBox_Is3State(c.ptr) != 0
}

// Get3StateValue returns the 3-state value:
//   - 0 = unchecked
//   - 1 = checked
//   - 2 = undetermined
func (c *CheckBox) Get3StateValue() int {
	return int(C.wxCheckBox_Get3StateValue(c.ptr))
}

// Set3StateValue sets the 3-state value:
//   - 0 = unchecked
//   - 1 = checked
//   - 2 = undetermined
//
// Example:
//
//	checkbox.Set3StateValue(2)  // Set to undetermined state
func (c *CheckBox) Set3StateValue(state int) {
	C.wxCheckBox_Set3StateValue(c.ptr, C.int(state))
}
