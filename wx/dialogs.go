// Package wx dialogs implementation
//
// This file implements common dialog functions.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// Dialog functions - use base class functions via inheritance
extern int wxDialog_ShowModal(void* self);
extern int wxWindow_Destroy(void* self);

extern int wxcMessageBox(void* message, void* caption, int style, void* parent, int x, int y);
extern void* wxFileDialog_Create(void* parent, void* message, void* defaultDir,
                                 void* defaultFile, void* wildcard, int x, int y, int style);
extern void* wxFileDialog_GetPath(void* self);
extern void* wxDirDialog_Create(void* parent, void* message, void* defaultPath, int x, int y, int style);
extern void* wxDirDialog_GetPath(void* self);
extern void* wxTextEntryDialog_Create(void* parent, void* message, void* caption,
                                      void* defaultValue, long style, int x, int y);
extern void* wxTextEntryDialog_GetValue(void* self);
extern void* wxColourDialog_Create(void* parent, void* data);
extern void* wxColourDialog_GetColourData(void* self);
extern void* wxFontDialog_Create(void* parent, void* data);
extern void* wxFontDialog_GetFontData(void* self);
*/
import "C"
import (
	"unsafe"
)

// =============================================================================
// Message Box
// =============================================================================

// MessageBox displays a message box dialog.
//
// Parameters:
//   - message: The message to display
//   - caption: The dialog title
//   - style: Dialog style flags (OK, YES_NO, ICON_INFORMATION, etc.)
//   - parent: Parent window (can be nil)
//
// Returns the button pressed (ID_OK, ID_CANCEL, ID_YES, ID_NO).
//
// Example:
//
//	result := wx.MessageBox("Save changes?", "Confirm",
//	    wx.YES_NO|wx.ICON_QUESTION, frame)
//	if result == wx.ID_YES {
//	    // Save changes
//	}
func MessageBox(message, caption string, style int, parent Window) int {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxMessage := NewWxString(message)
	defer wxMessage.Free()

	wxCaption := NewWxString(caption)
	defer wxCaption.Free()

	return int(C.wxcMessageBox(wxMessage.Ptr(), wxCaption.Ptr(),
		C.int(style), parentPtr, -1, -1))
}

// =============================================================================
// File Dialog
// =============================================================================

// FileDialog represents a file selection dialog.
type FileDialog struct {
	ptr unsafe.Pointer
}

// NewFileDialog creates a new file dialog.
//
// Parameters:
//   - parent: Parent window
//   - message: Dialog title
//   - defaultDir: Initial directory
//   - defaultFile: Initial filename
//   - wildcard: File type filter (e.g., "Text files (*.txt)|*.txt|All files (*.*)|*.*")
//   - style: Dialog style flags (FD_OPEN, FD_SAVE, FD_MULTIPLE, etc.)
//
// Example:
//
//	dlg := wx.NewFileDialog(frame, "Choose a file", "", "",
//	    "Text files (*.txt)|*.txt|All files (*.*)|*.*",
//	    wx.FD_OPEN|wx.FD_FILE_MUST_EXIST)
//	if dlg.ShowModal() == wx.ID_OK {
//	    path := dlg.GetPath()
//	    fmt.Println("Selected:", path)
//	}
//	dlg.Destroy()
func NewFileDialog(parent Window, message, defaultDir, defaultFile, wildcard string, style int) *FileDialog {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxMessage := NewWxString(message)
	defer wxMessage.Free()

	wxDefaultDir := NewWxString(defaultDir)
	defer wxDefaultDir.Free()

	wxDefaultFile := NewWxString(defaultFile)
	defer wxDefaultFile.Free()

	wxWildcard := NewWxString(wildcard)
	defer wxWildcard.Free()

	ptr := C.wxFileDialog_Create(parentPtr, wxMessage.Ptr(), wxDefaultDir.Ptr(),
		wxDefaultFile.Ptr(), wxWildcard.Ptr(), -1, -1, C.int(style))

	return &FileDialog{ptr: ptr}
}

// ShowModal displays the dialog and waits for user response.
// Returns ID_OK if accepted, ID_CANCEL if cancelled.
func (d *FileDialog) ShowModal() int {
	return int(C.wxDialog_ShowModal(d.ptr))
}

// GetPath returns the selected file path.
func (d *FileDialog) GetPath() string {
	wxStr := C.wxFileDialog_GetPath(d.ptr)
	return WxStringToGoAndFree(wxStr)
}

// Destroy destroys the dialog and frees resources.
func (d *FileDialog) Destroy() {
	if d.ptr != nil {
		C.wxWindow_Destroy(d.ptr)
		d.ptr = nil
	}
}

// =============================================================================
// Directory Dialog
// =============================================================================

// DirDialog represents a directory selection dialog.
type DirDialog struct {
	ptr unsafe.Pointer
}

// NewDirDialog creates a new directory dialog.
//
// Parameters:
//   - parent: Parent window
//   - message: Dialog title
//   - defaultPath: Initial directory path
//   - style: Dialog style flags
//
// Example:
//
//	dlg := wx.NewDirDialog(frame, "Choose a directory", "", 0)
//	if dlg.ShowModal() == wx.ID_OK {
//	    path := dlg.GetPath()
//	    fmt.Println("Selected:", path)
//	}
//	dlg.Destroy()
func NewDirDialog(parent Window, message, defaultPath string, style int) *DirDialog {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxMessage := NewWxString(message)
	defer wxMessage.Free()

	wxDefaultPath := NewWxString(defaultPath)
	defer wxDefaultPath.Free()

	ptr := C.wxDirDialog_Create(parentPtr, wxMessage.Ptr(), wxDefaultPath.Ptr(),
		-1, -1, C.int(style))

	return &DirDialog{ptr: ptr}
}

// ShowModal displays the dialog and waits for user response.
// Returns ID_OK if accepted, ID_CANCEL if cancelled.
func (d *DirDialog) ShowModal() int {
	return int(C.wxDialog_ShowModal(d.ptr))
}

// GetPath returns the selected directory path.
func (d *DirDialog) GetPath() string {
	wxStr := C.wxDirDialog_GetPath(d.ptr)
	return WxStringToGoAndFree(wxStr)
}

// Destroy destroys the dialog and frees resources.
func (d *DirDialog) Destroy() {
	if d.ptr != nil {
		C.wxWindow_Destroy(d.ptr)
		d.ptr = nil
	}
}

// =============================================================================
// Text Entry Dialog
// =============================================================================

// TextEntryDialog represents a simple text input dialog.
type TextEntryDialog struct {
	ptr unsafe.Pointer
}

// NewTextEntryDialog creates a new text entry dialog.
//
// Parameters:
//   - parent: Parent window
//   - message: Prompt message
//   - caption: Dialog title
//   - defaultValue: Initial text value
//   - style: Dialog style flags
//
// Example:
//
//	dlg := wx.NewTextEntryDialog(frame, "Enter your name:", "Name", "", wx.OK|wx.CANCEL)
//	if dlg.ShowModal() == wx.ID_OK {
//	    name := dlg.GetValue()
//	    fmt.Println("Name:", name)
//	}
//	dlg.Destroy()
func NewTextEntryDialog(parent Window, message, caption, defaultValue string, style int) *TextEntryDialog {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	wxMessage := NewWxString(message)
	defer wxMessage.Free()

	wxCaption := NewWxString(caption)
	defer wxCaption.Free()

	wxDefaultValue := NewWxString(defaultValue)
	defer wxDefaultValue.Free()

	ptr := C.wxTextEntryDialog_Create(parentPtr, wxMessage.Ptr(), wxCaption.Ptr(),
		wxDefaultValue.Ptr(), C.long(style), -1, -1)

	return &TextEntryDialog{ptr: ptr}
}

// ShowModal displays the dialog and waits for user response.
// Returns ID_OK if accepted, ID_CANCEL if cancelled.
func (d *TextEntryDialog) ShowModal() int {
	return int(C.wxDialog_ShowModal(d.ptr))
}

// GetValue returns the entered text.
func (d *TextEntryDialog) GetValue() string {
	wxStr := C.wxTextEntryDialog_GetValue(d.ptr)
	return WxStringToGoAndFree(wxStr)
}

// Destroy destroys the dialog and frees resources.
func (d *TextEntryDialog) Destroy() {
	if d.ptr != nil {
		C.wxWindow_Destroy(d.ptr)
		d.ptr = nil
	}
}

// =============================================================================
// Color Dialog
// =============================================================================

// ColourDialog represents a color selection dialog.
type ColourDialog struct {
	ptr unsafe.Pointer
}

// NewColourDialog creates a new color dialog.
//
// Parameters:
//   - parent: Parent window
//   - data: Optional color data (can be nil)
//
// Example:
//
//	dlg := wx.NewColourDialog(frame, nil)
//	if dlg.ShowModal() == wx.ID_OK {
//	    // Get selected color from dlg.GetColourData()
//	}
//	dlg.Destroy()
func NewColourDialog(parent Window, data unsafe.Pointer) *ColourDialog {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	ptr := C.wxColourDialog_Create(parentPtr, data)
	return &ColourDialog{ptr: ptr}
}

// ShowModal displays the dialog and waits for user response.
// Returns ID_OK if accepted, ID_CANCEL if cancelled.
func (d *ColourDialog) ShowModal() int {
	return int(C.wxDialog_ShowModal(d.ptr))
}

// GetColourData returns the color data containing the selected color.
func (d *ColourDialog) GetColourData() unsafe.Pointer {
	return C.wxColourDialog_GetColourData(d.ptr)
}

// Destroy destroys the dialog and frees resources.
func (d *ColourDialog) Destroy() {
	if d.ptr != nil {
		C.wxWindow_Destroy(d.ptr)
		d.ptr = nil
	}
}

// ColorDialog is an alias for ColourDialog (American spelling).
type ColorDialog = ColourDialog

// NewColorDialog is an alias for NewColourDialog (American spelling).
var NewColorDialog = NewColourDialog

// =============================================================================
// Font Dialog
// =============================================================================

// FontDialog represents a font selection dialog.
type FontDialog struct {
	ptr unsafe.Pointer
}

// NewFontDialog creates a new font dialog.
//
// Parameters:
//   - parent: Parent window
//   - data: Optional font data (can be nil)
//
// Example:
//
//	dlg := wx.NewFontDialog(frame, nil)
//	if dlg.ShowModal() == wx.ID_OK {
//	    // Get selected font from dlg.GetFontData()
//	}
//	dlg.Destroy()
func NewFontDialog(parent Window, data unsafe.Pointer) *FontDialog {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	ptr := C.wxFontDialog_Create(parentPtr, data)
	return &FontDialog{ptr: ptr}
}

// ShowModal displays the dialog and waits for user response.
// Returns ID_OK if accepted, ID_CANCEL if cancelled.
func (d *FontDialog) ShowModal() int {
	return int(C.wxDialog_ShowModal(d.ptr))
}

// GetFontData returns the font data containing the selected font.
func (d *FontDialog) GetFontData() unsafe.Pointer {
	return C.wxFontDialog_GetFontData(d.ptr)
}

// Destroy destroys the dialog and frees resources.
func (d *FontDialog) Destroy() {
	if d.ptr != nil {
		C.wxWindow_Destroy(d.ptr)
		d.ptr = nil
	}
}
