// Package wx menu implementation
//
// This file implements Menu and MenuBar wrappers for building application menus.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "kwx_def.h"
#include "kwx_types.h"

// Menu functions
extern void* wxMenu_Create(void* title, long style);
extern void wxMenu_Append(void* self, int id, void* text, void* help, int isCheckable);
extern void wxMenu_AppendSeparator(void* self);

// MenuBar functions
extern void* wxMenuBar_Create(int style);
extern int wxMenuBar_Append(void* self, void* menu, void* title);
*/
import "C"
import "unsafe"

// Menu represents a wxMenu — a list of menu items shown in a menu bar or as a popup.
type Menu struct {
	ptr unsafe.Pointer
}

// NewMenu creates a new menu with the given title.
//
// Example:
//
//	fileMenu := wx.NewMenu("&File")
//	fileMenu.Append(wx.ID_EXIT, "E&xit", "Exit application")
func NewMenu(title string) *Menu {
	wxTitle := NewWxString(title)
	defer wxTitle.Free()

	ptr := C.wxMenu_Create(wxTitle.Ptr(), 0)
	return &Menu{ptr: ptr}
}

// Ptr returns the underlying wxMenu pointer.
func (m *Menu) Ptr() unsafe.Pointer {
	return m.ptr
}

// Append adds a normal (non-checkable) menu item.
//
// Parameters:
//   - id: Menu item ID (e.g., wx.ID_EXIT, wx.ID_NEW)
//   - text: Item label with optional accelerator (e.g., "&Open\tCtrl+O")
//   - help: Status bar help string (shown when item is highlighted)
func (m *Menu) Append(id int, text, help string) {
	wxText := NewWxString(text)
	defer wxText.Free()

	wxHelp := NewWxString(help)
	defer wxHelp.Free()

	C.wxMenu_Append(m.ptr, C.int(id), wxText.Ptr(), wxHelp.Ptr(), 0)
}

// AppendSeparator adds a separator line to the menu.
func (m *Menu) AppendSeparator() {
	C.wxMenu_AppendSeparator(m.ptr)
}

// MenuBar represents a wxMenuBar — the horizontal bar of menus at the top of a frame.
type MenuBar struct {
	ptr unsafe.Pointer
}

// NewMenuBar creates a new empty menu bar.
//
// Example:
//
//	menuBar := wx.NewMenuBar()
//	menuBar.Append(fileMenu, "&File")
//	menuBar.Append(helpMenu, "&Help")
//	frame.SetMenuBar(menuBar)
func NewMenuBar() *MenuBar {
	ptr := C.wxMenuBar_Create(0)
	return &MenuBar{ptr: ptr}
}

// Ptr returns the underlying wxMenuBar pointer.
func (mb *MenuBar) Ptr() unsafe.Pointer {
	return mb.ptr
}

// Append adds a menu to the menu bar with the given title.
//
// Parameters:
//   - menu: The menu to add
//   - title: The title shown in the menu bar (use & for mnemonic)
func (mb *MenuBar) Append(menu *Menu, title string) {
	wxTitle := NewWxString(title)
	defer wxTitle.Free()

	C.wxMenuBar_Append(mb.ptr, menu.Ptr(), wxTitle.Ptr())
}
