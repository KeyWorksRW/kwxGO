// Package wx window implementation
//
// This file implements the Window interface and common window operations.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// Window functions
extern int wxWindow_Show(void* self);
extern int wxWindow_Hide(void* self);
extern int wxWindow_Close(void* self, int force);
extern int wxWindow_Destroy(void* self);
extern void wxWindow_SetSize(void* self, int x, int y, int width, int height, int sizeFlags);
extern void* wxWindow_GetSize(void* self);
extern void wxWindow_Move(void* self, int x, int y);
extern void* wxWindow_GetPosition(void* self);
extern void wxWindow_SetLabel(void* self, void* label);
extern int wxWindow_Enable(void* self);
extern int wxWindow_Disable(void* self);
extern int wxWindow_IsEnabled(void* self);
extern int wxWindow_IsShown(void* self);
extern void wxWindow_SetFocus(void* self);
extern void wxWindow_Refresh(void* self, int eraseBackground);
extern void wxWindow_Update(void* self);
extern int wxWindow_GetId(void* self);
extern void wxWindow_SetId(void* self, int id);
extern void* wxWindow_GetParent(void* self);
extern void wxWindow_SetSizer(void* self, void* sizer);
extern void wxWindow_SetSizerAndFit(void* self, void* sizer, int deleteOld);

// wxSize/wxPoint accessor functions
extern int wxSize_GetWidth(void* s);
extern int wxSize_GetHeight(void* s);
extern void wxSize_Delete(void* s);
extern int wxPoint_GetX(void* p);
extern int wxPoint_GetY(void* p);
extern void wxPoint_Delete(void* p);
*/
import "C"
import (
	"unsafe"
)

// Window interface represents any wxWindow.
// All window types (Frame, Panel, Dialog, Controls) implement this interface.
type Window interface {
	// Ptr returns the underlying wxWindow pointer
	Ptr() unsafe.Pointer

	// Show shows or hides the window
	Show(show bool) bool

	// Hide hides the window (equivalent to Show(false))
	Hide() bool

	// Close attempts to close the window
	Close(force bool) bool

	// Destroy destroys the window immediately
	Destroy() bool

	// IsShown returns true if the window is shown
	IsShown() bool

	// Enable enables or disables the window
	Enable(enable bool) bool

	// Disable disables the window (equivalent to Enable(false))
	Disable() bool

	// IsEnabled returns true if the window is enabled
	IsEnabled() bool

	// SetFocus sets keyboard focus to this window
	SetFocus()

	// Refresh refreshes the window
	Refresh(eraseBackground bool)

	// Update causes immediate repainting
	Update()

	// GetId returns the window identifier
	GetId() int

	// SetId sets the window identifier
	SetId(id int)

	// GetParent returns the parent window pointer
	GetParent() unsafe.Pointer

	// SetSize sets window position and size
	SetSize(x, y, width, height int)

	// GetSize returns window size
	GetSize() (width, height int)

	// Move moves the window
	Move(x, y int)

	// GetPosition returns window position
	GetPosition() (x, y int)

	// SetLabel sets the window label/title
	SetLabel(label string)

	// SetSizer sets the window's sizer
	SetSizer(sizer unsafe.Pointer)

	// SetSizerAndFit sets the sizer and resizes window to fit
	SetSizerAndFit(sizer unsafe.Pointer, deleteOld bool)
}

// Show shows the window.
// Returns true if the window's state was changed.
func (w *BaseWindow) Show(show bool) bool {
	if show {
		return C.wxWindow_Show(w.ptr) != 0
	}
	return C.wxWindow_Hide(w.ptr) != 0
}

// Hide hides the window.
func (w *BaseWindow) Hide() bool {
	return C.wxWindow_Hide(w.ptr) != 0
}

// Close attempts to close the window.
// If force is true, the window is closed unconditionally.
// If force is false, the window can veto the close.
func (w *BaseWindow) Close(force bool) bool {
	flag := C.int(0)
	if force {
		flag = 1
	}
	return C.wxWindow_Close(w.ptr, flag) != 0
}

// Destroy destroys the window and its children immediately.
// Returns true if the window was destroyed.
func (w *BaseWindow) Destroy() bool {
	return C.wxWindow_Destroy(w.ptr) != 0
}

// IsShown returns true if the window is currently shown.
func (w *BaseWindow) IsShown() bool {
	return C.wxWindow_IsShown(w.ptr) != 0
}

// Enable enables or disables the window.
// A disabled window cannot receive user input.
func (w *BaseWindow) Enable(enable bool) bool {
	if enable {
		return C.wxWindow_Enable(w.ptr) != 0
	}
	return C.wxWindow_Disable(w.ptr) != 0
}

// Disable disables the window.
func (w *BaseWindow) Disable() bool {
	return C.wxWindow_Disable(w.ptr) != 0
}

// IsEnabled returns true if the window is currently enabled.
func (w *BaseWindow) IsEnabled() bool {
	return C.wxWindow_IsEnabled(w.ptr) != 0
}

// SetFocus sets the keyboard focus to this window.
func (w *BaseWindow) SetFocus() {
	C.wxWindow_SetFocus(w.ptr)
}

// Refresh refreshes the window by causing a repaint.
// If eraseBackground is true, the background is erased first.
func (w *BaseWindow) Refresh(eraseBackground bool) {
	flag := C.int(0)
	if eraseBackground {
		flag = 1
	}
	C.wxWindow_Refresh(w.ptr, flag)
}

// Update causes an immediate repaint of the window.
// Unlike Refresh, this processes paint events immediately.
func (w *BaseWindow) Update() {
	C.wxWindow_Update(w.ptr)
}

// GetId returns the window identifier.
func (w *BaseWindow) GetId() int {
	return int(C.wxWindow_GetId(w.ptr))
}

// SetId sets the window identifier.
func (w *BaseWindow) SetId(id int) {
	C.wxWindow_SetId(w.ptr, C.int(id))
}

// GetParent returns the parent window pointer, or nil if this is a top-level window.
func (w *BaseWindow) GetParent() unsafe.Pointer {
	return C.wxWindow_GetParent(w.ptr)
}

// SetSize sets the window position and size.
// Pass -1 for parameters you don't want to change.
func (w *BaseWindow) SetSize(x, y, width, height int) {
	C.wxWindow_SetSize(w.ptr, C.int(x), C.int(y), C.int(width), C.int(height), 0)
}

// GetSize returns the window size.
func (w *BaseWindow) GetSize() (width, height int) {
	szPtr := C.wxWindow_GetSize(w.ptr)
	width = int(C.wxSize_GetWidth(szPtr))
	height = int(C.wxSize_GetHeight(szPtr))
	C.wxSize_Delete(szPtr)
	return
}

// Move moves the window to the specified position.
func (w *BaseWindow) Move(x, y int) {
	C.wxWindow_Move(w.ptr, C.int(x), C.int(y))
}

// GetPosition returns the window position.
func (w *BaseWindow) GetPosition() (x, y int) {
	ptPtr := C.wxWindow_GetPosition(w.ptr)
	x = int(C.wxPoint_GetX(ptPtr))
	y = int(C.wxPoint_GetY(ptPtr))
	C.wxPoint_Delete(ptPtr)
	return
}

// SetLabel sets the window label/title.
func (w *BaseWindow) SetLabel(label string) {
	wxStr := NewWxString(label)
	defer wxStr.Free()
	C.wxWindow_SetLabel(w.ptr, wxStr.Ptr())
}

// SetSizer sets the window's sizer.
func (w *BaseWindow) SetSizer(sizer unsafe.Pointer) {
	C.wxWindow_SetSizer(w.ptr, sizer)
}

// SetSizerAndFit sets the window's sizer and resizes the window to fit the sizer.
// If deleteOld is true, the old sizer (if any) is deleted.
func (w *BaseWindow) SetSizerAndFit(sizer unsafe.Pointer, deleteOld bool) {
	flag := C.int(0)
	if deleteOld {
		flag = 1
	}
	C.wxWindow_SetSizerAndFit(w.ptr, sizer, flag)
}

// SetSizerInterface sets the window's sizer using a Sizer interface.
// This is a convenience method that accepts a Sizer instead of unsafe.Pointer.
//
// Example:
//
//	sizer := wx.NewBoxSizer(wx.VERTICAL)
//	sizer.Add(button, 0, wx.ALL, 5)
//	panel.SetSizerInterface(sizer)
func (w *BaseWindow) SetSizerInterface(sizer Sizer) {
	w.SetSizer(sizer.Ptr())
}

// SetSizerAndFitInterface sets the window's sizer and resizes using a Sizer interface.
// This is a convenience method that accepts a Sizer instead of unsafe.Pointer.
//
// Example:
//
//	sizer := wx.NewBoxSizer(wx.VERTICAL)
//	sizer.Add(button, 0, wx.ALL, 5)
//	panel.SetSizerAndFitInterface(sizer, true)
func (w *BaseWindow) SetSizerAndFitInterface(sizer Sizer, deleteOld bool) {
	w.SetSizerAndFit(sizer.Ptr(), deleteOld)
}
