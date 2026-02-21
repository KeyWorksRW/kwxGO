// Package wx event type wrappers.
//
// This file provides typed wrappers for specific wxWidgets event types,
// giving access to event-specific data beyond what the base Event provides.
//
// Use the As*Event conversion functions to cast a base Event to a specific type:
//
//	button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
//	    cmd := wx.AsCommandEvent(evt)
//	    fmt.Printf("Selection: %d\n", cmd.GetSelection())
//	})
//
// Available event wrappers:
//   - CommandEvent — button clicks, menu selections, checkbox toggles, etc.
//   - CloseEvent — window close requests (can be vetoed)
//   - KeyEvent — keyboard events (key down, key up, char)
//   - MouseEvent — mouse events (click, move, enter, leave)
//   - SizeEvent — window resize events

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include "kwx_def.h"
#include "kwx_types.h"

// CommandEvent functions
extern int wxCommandEvent_GetSelection(void* self);
extern void* wxCommandEvent_GetString(void* self);
extern long wxCommandEvent_GetInt(void* self);
extern long wxCommandEvent_GetExtraLong(void* self);
extern int wxCommandEvent_IsChecked(void* self);
extern int wxCommandEvent_IsSelection(void* self);

// CloseEvent functions
extern int wxCloseEvent_CanVeto(void* self);
extern int wxCloseEvent_GetVeto(void* self);
extern void wxCloseEvent_Veto(void* self, int veto);
extern int wxCloseEvent_GetLoggingOff(void* self);
extern void wxCloseEvent_SetCanVeto(void* self, int canVeto);
extern void wxCloseEvent_SetLoggingOff(void* self, int logOff);

// KeyEvent functions
extern int wxKeyEvent_GetKeyCode(void* self);
extern int wxKeyEvent_GetX(void* self);
extern int wxKeyEvent_GetY(void* self);
extern int wxKeyEvent_AltDown(void* self);
extern int wxKeyEvent_ControlDown(void* self);
extern int wxKeyEvent_ShiftDown(void* self);
extern int wxKeyEvent_MetaDown(void* self);
extern int wxKeyEvent_GetModifiers(void* self);
extern int wxKeyEvent_HasModifiers(void* self);

// MouseEvent functions
extern int wxMouseEvent_GetX(void* self);
extern int wxMouseEvent_GetY(void* self);
extern int wxMouseEvent_AltDown(void* self);
extern int wxMouseEvent_ControlDown(void* self);
extern int wxMouseEvent_ShiftDown(void* self);
extern int wxMouseEvent_MetaDown(void* self);
extern int wxMouseEvent_LeftDown(void* self);
extern int wxMouseEvent_LeftUp(void* self);
extern int wxMouseEvent_LeftDClick(void* self);
extern int wxMouseEvent_LeftIsDown(void* self);
extern int wxMouseEvent_MiddleDown(void* self);
extern int wxMouseEvent_MiddleUp(void* self);
extern int wxMouseEvent_MiddleIsDown(void* self);
extern int wxMouseEvent_RightDown(void* self);
extern int wxMouseEvent_RightUp(void* self);
extern int wxMouseEvent_Dragging(void* self);
extern int wxMouseEvent_Moving(void* self);
extern int wxMouseEvent_Entering(void* self);
extern int wxMouseEvent_Leaving(void* self);
extern int wxMouseEvent_Button(void* self, int but);
extern int wxMouseEvent_ButtonDown(void* self, int but);
extern int wxMouseEvent_ButtonUp(void* self, int but);
extern int wxMouseEvent_ButtonDClick(void* self, int but);
extern int wxMouseEvent_ButtonIsDown(void* self, int but);
extern int wxMouseEvent_IsButton(void* self);

// SizeEvent functions
extern void* wxSizeEvent_GetSize(void* self);
extern int wxSize_GetWidth(void* self);
extern int wxSize_GetHeight(void* self);
extern void wxSize_Delete(void* self);
*/
import "C"

// =============================================================================
// CommandEvent — button clicks, menu selections, text changes, etc.
// =============================================================================

// CommandEvent wraps wxCommandEvent for command events.
// Command events propagate upward through the window hierarchy.
type CommandEvent struct {
	Event
}

// AsCommandEvent converts a base Event to a CommandEvent.
// Use this in handlers for command events (EVT_BUTTON, EVT_MENU,
// EVT_TEXT, EVT_CHECKBOX, etc.).
//
// Example:
//
//	button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
//	    cmd := wx.AsCommandEvent(evt)
//	    fmt.Println("Selection:", cmd.GetSelection())
//	})
func AsCommandEvent(e *Event) *CommandEvent {
	return &CommandEvent{Event: *e}
}

// GetSelection returns the index of the selected item.
// For list boxes, choice controls, and similar multi-item controls.
func (e *CommandEvent) GetSelection() int {
	return int(C.wxCommandEvent_GetSelection(e.ptr))
}

// GetString returns the string associated with the command event.
// Returns the label of the selected item for choice/list controls,
// or the value of the text control for text events.
func (e *CommandEvent) GetString() string {
	wxStr := C.wxCommandEvent_GetString(e.ptr)
	if wxStr == nil {
		return ""
	}
	return WxStringToGoAndFree(wxStr)
}

// GetInt returns the integer value associated with the command event.
// The meaning depends on the event type (e.g., checkbox state).
func (e *CommandEvent) GetInt() int {
	return int(C.wxCommandEvent_GetInt(e.ptr))
}

// GetExtraLong returns the extra long value.
// Used by scroll events and similar for additional data.
func (e *CommandEvent) GetExtraLong() int {
	return int(C.wxCommandEvent_GetExtraLong(e.ptr))
}

// IsChecked returns true if the item is checked.
// Use for checkbox click events and check menu item events.
func (e *CommandEvent) IsChecked() bool {
	return C.wxCommandEvent_IsChecked(e.ptr) != 0
}

// IsSelection returns true if this is a selection event (not deselection).
// For list box and similar controls that support deselection.
func (e *CommandEvent) IsSelection() bool {
	return C.wxCommandEvent_IsSelection(e.ptr) != 0
}

// =============================================================================
// CloseEvent — window close requests
// =============================================================================

// CloseEvent wraps wxCloseEvent for window close events.
// Close events can optionally be vetoed to prevent the window from closing.
type CloseEvent struct {
	Event
}

// AsCloseEvent converts a base Event to a CloseEvent.
// Use this in handlers for EVT_CLOSE_WINDOW events.
//
// Example:
//
//	frame.Bind(wx.EVT_CLOSE_WINDOW, func(evt *wx.Event) {
//	    closeEvt := wx.AsCloseEvent(evt)
//	    if closeEvt.CanVeto() {
//	        // Ask user for confirmation, then:
//	        closeEvt.Veto(true)  // prevent close
//	        // or evt.Skip()     // allow close
//	    } else {
//	        evt.Skip() // forced close, must allow
//	    }
//	})
func AsCloseEvent(e *Event) *CloseEvent {
	return &CloseEvent{Event: *e}
}

// CanVeto returns true if the close event can be vetoed.
// If false, the window MUST be closed (e.g., system shutdown).
func (e *CloseEvent) CanVeto() bool {
	return C.wxCloseEvent_CanVeto(e.ptr) != 0
}

// GetVeto returns true if the close has been vetoed.
func (e *CloseEvent) GetVeto() bool {
	return C.wxCloseEvent_GetVeto(e.ptr) != 0
}

// Veto vetoes or allows the close event.
// Pass true to prevent closing, false to allow it.
// Only effective if CanVeto() returns true.
func (e *CloseEvent) Veto(veto bool) {
	flag := C.int(0)
	if veto {
		flag = 1
	}
	C.wxCloseEvent_Veto(e.ptr, flag)
}

// GetLoggingOff returns true if the application is logging off.
// This is only relevant for EVT_END_SESSION events.
func (e *CloseEvent) GetLoggingOff() bool {
	return C.wxCloseEvent_GetLoggingOff(e.ptr) != 0
}

// =============================================================================
// KeyEvent — keyboard events
// =============================================================================

// KeyEvent wraps wxKeyEvent for keyboard events.
// Use with EVT_KEY_DOWN, EVT_KEY_UP, and EVT_CHAR events.
type KeyEvent struct {
	Event
}

// AsKeyEvent converts a base Event to a KeyEvent.
// Use this in handlers for keyboard events.
func AsKeyEvent(e *Event) *KeyEvent {
	return &KeyEvent{Event: *e}
}

// GetKeyCode returns the key code for this event.
// Returns a wxKeyCode value (e.g., WXK_RETURN, WXK_ESCAPE, or ASCII value).
func (e *KeyEvent) GetKeyCode() int {
	return int(C.wxKeyEvent_GetKeyCode(e.ptr))
}

// GetX returns the X position of the key event in client coordinates.
func (e *KeyEvent) GetX() int {
	return int(C.wxKeyEvent_GetX(e.ptr))
}

// GetY returns the Y position of the key event in client coordinates.
func (e *KeyEvent) GetY() int {
	return int(C.wxKeyEvent_GetY(e.ptr))
}

// GetPosition returns the position of the key event as a Point.
func (e *KeyEvent) GetPosition() Point {
	return Point{X: e.GetX(), Y: e.GetY()}
}

// AltDown returns true if the Alt key was held during the event.
func (e *KeyEvent) AltDown() bool {
	return C.wxKeyEvent_AltDown(e.ptr) != 0
}

// ControlDown returns true if the Control key was held during the event.
func (e *KeyEvent) ControlDown() bool {
	return C.wxKeyEvent_ControlDown(e.ptr) != 0
}

// ShiftDown returns true if the Shift key was held during the event.
func (e *KeyEvent) ShiftDown() bool {
	return C.wxKeyEvent_ShiftDown(e.ptr) != 0
}

// MetaDown returns true if the Meta/Windows key was held during the event.
func (e *KeyEvent) MetaDown() bool {
	return C.wxKeyEvent_MetaDown(e.ptr) != 0
}

// GetModifiers returns the modifier keys as a bitmask.
func (e *KeyEvent) GetModifiers() int {
	return int(C.wxKeyEvent_GetModifiers(e.ptr))
}

// HasModifiers returns true if any modifier key (Ctrl, Alt) was held.
// Note: Shift alone does not count as a modifier for this function.
func (e *KeyEvent) HasModifiers() bool {
	return C.wxKeyEvent_HasModifiers(e.ptr) != 0
}

// =============================================================================
// MouseEvent — mouse events
// =============================================================================

// MouseEvent wraps wxMouseEvent for mouse events.
// Use with EVT_LEFT_DOWN, EVT_LEFT_UP, EVT_MOTION, etc.
type MouseEvent struct {
	Event
}

// AsMouseEvent converts a base Event to a MouseEvent.
// Use this in handlers for mouse events.
func AsMouseEvent(e *Event) *MouseEvent {
	return &MouseEvent{Event: *e}
}

// GetX returns the X coordinate of the mouse event in client coordinates.
func (e *MouseEvent) GetX() int {
	return int(C.wxMouseEvent_GetX(e.ptr))
}

// GetY returns the Y coordinate of the mouse event in client coordinates.
func (e *MouseEvent) GetY() int {
	return int(C.wxMouseEvent_GetY(e.ptr))
}

// GetPosition returns the mouse position as a Point.
func (e *MouseEvent) GetPosition() Point {
	return Point{X: e.GetX(), Y: e.GetY()}
}

// AltDown returns true if the Alt key was held during the mouse event.
func (e *MouseEvent) AltDown() bool {
	return C.wxMouseEvent_AltDown(e.ptr) != 0
}

// ControlDown returns true if the Control key was held during the mouse event.
func (e *MouseEvent) ControlDown() bool {
	return C.wxMouseEvent_ControlDown(e.ptr) != 0
}

// ShiftDown returns true if the Shift key was held during the mouse event.
func (e *MouseEvent) ShiftDown() bool {
	return C.wxMouseEvent_ShiftDown(e.ptr) != 0
}

// MetaDown returns true if the Meta/Windows key was held during the mouse event.
func (e *MouseEvent) MetaDown() bool {
	return C.wxMouseEvent_MetaDown(e.ptr) != 0
}

// LeftDown returns true if the left mouse button was just pressed.
func (e *MouseEvent) LeftDown() bool {
	return C.wxMouseEvent_LeftDown(e.ptr) != 0
}

// LeftUp returns true if the left mouse button was just released.
func (e *MouseEvent) LeftUp() bool {
	return C.wxMouseEvent_LeftUp(e.ptr) != 0
}

// LeftDClick returns true if the left mouse button was double-clicked.
func (e *MouseEvent) LeftDClick() bool {
	return C.wxMouseEvent_LeftDClick(e.ptr) != 0
}

// LeftIsDown returns true if the left mouse button is currently pressed.
func (e *MouseEvent) LeftIsDown() bool {
	return C.wxMouseEvent_LeftIsDown(e.ptr) != 0
}

// MiddleDown returns true if the middle mouse button was just pressed.
func (e *MouseEvent) MiddleDown() bool {
	return C.wxMouseEvent_MiddleDown(e.ptr) != 0
}

// MiddleUp returns true if the middle mouse button was just released.
func (e *MouseEvent) MiddleUp() bool {
	return C.wxMouseEvent_MiddleUp(e.ptr) != 0
}

// MiddleIsDown returns true if the middle mouse button is currently pressed.
func (e *MouseEvent) MiddleIsDown() bool {
	return C.wxMouseEvent_MiddleIsDown(e.ptr) != 0
}

// RightDown returns true if the right mouse button was just pressed.
func (e *MouseEvent) RightDown() bool {
	return C.wxMouseEvent_RightDown(e.ptr) != 0
}

// RightUp returns true if the right mouse button was just released.
func (e *MouseEvent) RightUp() bool {
	return C.wxMouseEvent_RightUp(e.ptr) != 0
}

// Dragging returns true if the mouse is being dragged (moved with a button down).
func (e *MouseEvent) Dragging() bool {
	return C.wxMouseEvent_Dragging(e.ptr) != 0
}

// Moving returns true if the mouse is moving (without any buttons pressed).
func (e *MouseEvent) Moving() bool {
	return C.wxMouseEvent_Moving(e.ptr) != 0
}

// Entering returns true if the mouse is entering the window.
func (e *MouseEvent) Entering() bool {
	return C.wxMouseEvent_Entering(e.ptr) != 0
}

// Leaving returns true if the mouse is leaving the window.
func (e *MouseEvent) Leaving() bool {
	return C.wxMouseEvent_Leaving(e.ptr) != 0
}

// IsButton returns true if this is any type of button event.
func (e *MouseEvent) IsButton() bool {
	return C.wxMouseEvent_IsButton(e.ptr) != 0
}

// Button returns true if the specified button (1=left, 2=right, 3=middle)
// triggered this event.
func (e *MouseEvent) Button(but int) bool {
	return C.wxMouseEvent_Button(e.ptr, C.int(but)) != 0
}

// ButtonDown returns true if the specified button was just pressed.
func (e *MouseEvent) ButtonDown(but int) bool {
	return C.wxMouseEvent_ButtonDown(e.ptr, C.int(but)) != 0
}

// ButtonUp returns true if the specified button was just released.
func (e *MouseEvent) ButtonUp(but int) bool {
	return C.wxMouseEvent_ButtonUp(e.ptr, C.int(but)) != 0
}

// ButtonDClick returns true if the specified button was double-clicked.
func (e *MouseEvent) ButtonDClick(but int) bool {
	return C.wxMouseEvent_ButtonDClick(e.ptr, C.int(but)) != 0
}

// ButtonIsDown returns true if the specified button is currently pressed.
func (e *MouseEvent) ButtonIsDown(but int) bool {
	return C.wxMouseEvent_ButtonIsDown(e.ptr, C.int(but)) != 0
}

// =============================================================================
// SizeEvent — window resize events
// =============================================================================

// SizeEvent wraps wxSizeEvent for window resize events.
// Use with EVT_SIZE events.
type SizeEvent struct {
	Event
}

// AsSizeEvent converts a base Event to a SizeEvent.
// Use this in handlers for EVT_SIZE events.
//
// Example:
//
//	frame.Bind(wx.EVT_SIZE, func(evt *wx.Event) {
//	    sizeEvt := wx.AsSizeEvent(evt)
//	    size := sizeEvt.GetSize()
//	    fmt.Printf("New size: %dx%d\n", size.W, size.H)
//	    evt.Skip() // Allow default size handling
//	})
func AsSizeEvent(e *Event) *SizeEvent {
	return &SizeEvent{Event: *e}
}

// GetSize returns the new window size after the resize.
func (e *SizeEvent) GetSize() Size {
	wxSize := C.wxSizeEvent_GetSize(e.ptr)
	if wxSize == nil {
		return Size{}
	}
	w := int(C.wxSize_GetWidth(wxSize))
	h := int(C.wxSize_GetHeight(wxSize))
	C.wxSize_Delete(wxSize)
	return Size{W: w, H: h}
}
