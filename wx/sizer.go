// Package wx sizer implementation
//
// This file implements the sizer system for automatic layout management.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// Sizer functions
extern void* wxBoxSizer_Create(int orient);
extern void* wxFlexGridSizer_Create(int rows, int cols, int vgap, int hgap);
extern void wxFlexGridSizer_AddGrowableRow(void* self, size_t idx);
extern void wxFlexGridSizer_AddGrowableCol(void* self, size_t idx);
extern void* wxGridBagSizer_Create(int vgap, int hgap);
extern void wxGridBagSizer_AddWindow(void* self, void* window, int row, int col, int rowspan, int colspan, int flag, int border, void* userData);
extern void* wxStaticBox_Create(void* parent, int id, void* text, int x, int y, int w, int h, int style);
extern void* wxStaticBoxSizer_Create(void* box, int orient);
extern void* wxStaticBoxSizer_GetStaticBox(void* self);
extern void wxSizer_AddWindow(void* self, void* window, int proportion, int flag, int border, void* userData);
extern void wxSizer_AddSizer(void* self, void* sizer, int proportion, int flag, int border, void* userData);
extern void wxSizer_AddSpacer(void* self, int size);
extern void wxSizer_AddStretchSpacer(void* self, int prop);
extern void wxSizer_Layout(void* self);
extern void wxSizer_SetMinSize(void* self, int width, int height);
extern void wxSizer_Fit(void* self, void* window);
*/
import "C"
import (
	"unsafe"
)

// Sizer interface represents any wxSizer.
type Sizer interface {
	// Ptr returns the underlying wxSizer pointer
	Ptr() unsafe.Pointer

	// Add adds a window to the sizer
	Add(window Window, proportion, flag, border int)

	// AddSizer adds a nested sizer
	AddSizer(sizer Sizer, proportion, flag, border int)

	// AddSpacer adds fixed-size empty space
	AddSpacer(size int)

	// AddStretchSpacer adds stretchable empty space
	AddStretchSpacer(proportion int)

	// Layout recalculates and applies the layout
	Layout()

	// SetMinSize sets the minimum size for the sizer
	SetMinSize(width, height int)

	// Fit resizes window to fit sizer contents
	Fit(window Window)
}

// BaseSizer provides common sizer implementation.
type BaseSizer struct {
	ptr unsafe.Pointer
}

// Ptr returns the underlying wxSizer pointer.
func (s *BaseSizer) Ptr() unsafe.Pointer {
	return s.ptr
}

// Add adds a window to the sizer.
//
// Parameters:
//   - window: Window to add
//   - proportion: Proportion for resizing (0 = fixed size, >0 = grows proportionally)
//   - flag: Layout flags (EXPAND, ALIGN_*, etc.)
//   - border: Border width in pixels
//
// Example:
//
//	sizer.Add(button, 0, wx.ALIGN_CENTER|wx.ALL, 5)
func (s *BaseSizer) Add(window Window, proportion, flag, border int) {
	C.wxSizer_AddWindow(s.ptr, window.Ptr(), C.int(proportion), C.int(flag), C.int(border), nil)
}

// AddSizer adds a nested sizer to this sizer.
//
// Example:
//
//	mainSizer.AddSizer(subSizer, 1, wx.EXPAND|wx.ALL, 10)
func (s *BaseSizer) AddSizer(sizer Sizer, proportion, flag, border int) {
	C.wxSizer_AddSizer(s.ptr, sizer.Ptr(), C.int(proportion), C.int(flag), C.int(border), nil)
}

// AddSpacer adds fixed-size empty space to the sizer.
//
// Example:
//
//	sizer.AddSpacer(10)  // Add 10 pixels of space
func (s *BaseSizer) AddSpacer(size int) {
	C.wxSizer_AddSpacer(s.ptr, C.int(size))
}

// AddStretchSpacer adds stretchable empty space to the sizer.
//
// Parameters:
//   - proportion: How much the spacer grows relative to other stretchable items
//
// Example:
//
//	sizer.AddStretchSpacer(1)  // Add space that takes remaining room
func (s *BaseSizer) AddStretchSpacer(proportion int) {
	C.wxSizer_AddStretchSpacer(s.ptr, C.int(proportion))
}

// Layout recalculates and applies the layout.
// Usually not needed as it's called automatically.
func (s *BaseSizer) Layout() {
	C.wxSizer_Layout(s.ptr)
}

// SetMinSize sets the minimum size for the sizer.
//
// Example:
//
//	sizer.SetMinSize(200, 100)
func (s *BaseSizer) SetMinSize(width, height int) {
	C.wxSizer_SetMinSize(s.ptr, C.int(width), C.int(height))
}

// Fit resizes the window to fit the sizer's contents.
//
// Example:
//
//	sizer.Fit(frame)
func (s *BaseSizer) Fit(window Window) {
	C.wxSizer_Fit(s.ptr, window.Ptr())
}

// =============================================================================
// BoxSizer
// =============================================================================

// BoxSizer arranges items in a single row or column.
type BoxSizer struct {
	BaseSizer
}

// NewBoxSizer creates a new box sizer.
//
// Parameters:
//   - orient: Orientation (HORIZONTAL or VERTICAL)
//
// Example:
//
//	// Horizontal layout
//	hSizer := wx.NewBoxSizer(wx.HORIZONTAL)
//	hSizer.Add(button1, 0, wx.ALL, 5)
//	hSizer.Add(button2, 0, wx.ALL, 5)
//
//	// Vertical layout
//	vSizer := wx.NewBoxSizer(wx.VERTICAL)
//	vSizer.Add(label, 0, wx.ALL, 5)
//	vSizer.Add(textCtrl, 0, wx.EXPAND|wx.ALL, 5)
func NewBoxSizer(orient int) *BoxSizer {
	ptr := C.wxBoxSizer_Create(C.int(orient))
	return &BoxSizer{BaseSizer{ptr: ptr}}
}

// =============================================================================
// FlexGridSizer
// =============================================================================

// FlexGridSizer arranges items in a grid with flexible rows/columns.
type FlexGridSizer struct {
	BaseSizer
}

// NewFlexGridSizer creates a new flexible grid sizer.
//
// Parameters:
//   - rows: Number of rows (0 = calculated automatically)
//   - cols: Number of columns (0 = calculated automatically)
//   - vgap: Vertical gap between cells in pixels
//   - hgap: Horizontal gap between cells in pixels
//
// At least one of rows or cols must be non-zero.
//
// Example:
//
//	// 2 columns, automatic rows, 5px gaps
//	sizer := wx.NewFlexGridSizer(0, 2, 5, 5)
//	sizer.Add(label1, 0, wx.ALIGN_RIGHT|wx.ALIGN_CENTER_VERTICAL, 0)
//	sizer.Add(textCtrl1, 0, wx.EXPAND, 0)
//	sizer.Add(label2, 0, wx.ALIGN_RIGHT|wx.ALIGN_CENTER_VERTICAL, 0)
//	sizer.Add(textCtrl2, 0, wx.EXPAND, 0)
//	sizer.AddGrowableCol(1, 1)  // Make second column stretchable
func NewFlexGridSizer(rows, cols, vgap, hgap int) *FlexGridSizer {
	ptr := C.wxFlexGridSizer_Create(C.int(rows), C.int(cols), C.int(vgap), C.int(hgap))
	return &FlexGridSizer{BaseSizer{ptr: ptr}}
}

// AddGrowableRow makes a row growable (stretchable).
//
// Parameters:
//   - idx: Row index (0-based)
func (s *FlexGridSizer) AddGrowableRow(idx int) {
	C.wxFlexGridSizer_AddGrowableRow(s.ptr, C.size_t(idx))
}

// AddGrowableCol makes a column growable (stretchable).
//
// Parameters:
//   - idx: Column index (0-based)
//
// Example:
//
//	sizer.AddGrowableCol(1)  // Make second column take remaining space
func (s *FlexGridSizer) AddGrowableCol(idx int) {
	C.wxFlexGridSizer_AddGrowableCol(s.ptr, C.size_t(idx))
}

// =============================================================================
// GridBagSizer
// =============================================================================

// GridBagSizer is the most flexible sizer, allowing items to span multiple cells.
type GridBagSizer struct {
	BaseSizer
}

// NewGridBagSizer creates a new grid bag sizer.
//
// Parameters:
//   - vgap: Vertical gap between cells in pixels
//   - hgap: Horizontal gap between cells in pixels
//
// Example:
//
//	sizer := wx.NewGridBagSizer(5, 5)
//	sizer.AddAt(label, 0, 0, 1, 1, wx.ALIGN_CENTER_VERTICAL, 0)
//	sizer.AddAt(textCtrl, 0, 1, 1, 2, wx.EXPAND, 0)  // Spans 2 columns
func NewGridBagSizer(vgap, hgap int) *GridBagSizer {
	ptr := C.wxGridBagSizer_Create(C.int(vgap), C.int(hgap))
	return &GridBagSizer{BaseSizer{ptr: ptr}}
}

// AddAt adds a window at a specific grid position with optional spanning.
//
// Parameters:
//   - window: Window to add
//   - row: Row position (0-based)
//   - col: Column position (0-based)
//   - rowspan: Number of rows to span (1 = single row)
//   - colspan: Number of columns to span (1 = single column)
//   - flag: Layout flags
//   - border: Border width in pixels
//
// Example:
//
//	sizer.AddAt(bigWidget, 0, 0, 2, 2, wx.EXPAND, 5)  // Spans 2x2 cells
func (s *GridBagSizer) AddAt(window Window, row, col, rowspan, colspan, flag, border int) {
	C.wxGridBagSizer_AddWindow(s.ptr, window.Ptr(), C.int(row), C.int(col),
		C.int(rowspan), C.int(colspan), C.int(flag), C.int(border), nil)
}

// =============================================================================
// StaticBoxSizer
// =============================================================================

// StaticBoxSizer is a sizer with a static box (frame) around it.
type StaticBoxSizer struct {
	BaseSizer
}

// NewStaticBoxSizer creates a new static box sizer.
//
// Parameters:
//   - orient: Orientation (HORIZONTAL or VERTICAL)
//   - parent: Parent window
//   - label: Label for the static box
//
// Example:
//
//	sizer := wx.NewStaticBoxSizer(wx.VERTICAL, panel, "Options")
//	sizer.Add(checkbox1, 0, wx.ALL, 5)
//	sizer.Add(checkbox2, 0, wx.ALL, 5)
//	panel.SetSizer(sizer)
func NewStaticBoxSizer(orient int, parent Window, label string) *StaticBoxSizer {
	wxLabel := NewWxString(label)
	defer wxLabel.Free()

	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	// Create the static box first, then wrap it in a sizer
	box := C.wxStaticBox_Create(parentPtr, C.int(ID_ANY), wxLabel.Ptr(),
		C.int(-1), C.int(-1), C.int(-1), C.int(-1), C.int(0))
	ptr := C.wxStaticBoxSizer_Create(box, C.int(orient))
	return &StaticBoxSizer{BaseSizer{ptr: ptr}}
}

// GetStaticBox returns the underlying static box window.
func (s *StaticBoxSizer) GetStaticBox() unsafe.Pointer {
	return C.wxStaticBoxSizer_GetStaticBox(s.ptr)
}
