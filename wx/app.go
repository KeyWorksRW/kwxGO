// Package wx application lifecycle functions
//
// This file provides high-level Go wrappers for kwxApp functions beyond the
// basic Initialize/MainLoop/Shutdown lifecycle (which are in wx.go).

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include -I${SRCDIR}/../build/_deps/kwxffi-src/examples/CApp

#include <stdlib.h>
#include "kwxApp.h"
*/
import "C"
import (
	"unsafe"
)

// Application Properties

// GetAppName returns the application name.
// The application name is used in dialogs, config files, etc.
func GetAppName() string {
	cname := C.kwxApp_GetAppName()
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// SetAppName sets the application name.
// The application name is used in dialogs, config files, etc.
func SetAppName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.kwxApp_SetAppName(cname)
}

// GetVendorName returns the vendor/organization name.
func GetVendorName() string {
	cname := C.kwxApp_GetVendorName()
	if cname == nil {
		return ""
	}
	return C.GoString(cname)
}

// SetVendorName sets the vendor/organization name.
func SetVendorName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.kwxApp_SetVendorName(cname)
}

// GetTopWindow returns the top-level window pointer.
// Returns nil if no top window is set.
func GetTopWindow() unsafe.Pointer {
	return C.kwxApp_GetTopWindow()
}

// SetExitOnFrameDelete controls whether the application exits when the
// last top-level frame is deleted. Default is true.
func SetExitOnFrameDelete(flag bool) {
	cflag := C.int(0)
	if flag {
		cflag = 1
	}
	C.kwxApp_SetExitOnFrameDelete(cflag)
}

// GetExitOnFrameDelete returns true if the application will exit when
// the last top-level frame is deleted.
func GetExitOnFrameDelete() bool {
	return C.kwxApp_GetExitOnFrameDelete() != 0
}

// System Information

// GetDisplaySize returns the display size in pixels.
func GetDisplaySize() (width, height int) {
	var w, h C.int
	C.kwxApp_GetDisplaySize(&w, &h)
	return int(w), int(h)
}

// GetMousePosition returns the current mouse position in screen coordinates.
func GetMousePosition() (x, y int) {
	var cx, cy C.int
	C.kwxApp_GetMousePosition(&cx, &cy)
	return int(cx), int(cy)
}

// GetOsVersion returns the OS version information.
// Returns (osType, major, minor) where osType is the wxWidgets OS constant.
func GetOsVersion() (osType, major, minor int) {
	var maj, min C.int
	ostype := C.kwxApp_GetOsVersion(&maj, &min)
	return int(ostype), int(maj), int(min)
}

// GetOsDescription returns a human-readable OS description string.
func GetOsDescription() string {
	cstr := C.kwxApp_GetOsDescription()
	if cstr == nil {
		return ""
	}
	defer C.kwxApp_FreeString(cstr)
	return C.GoString(cstr)
}

// GetUserId returns the user's login ID.
func GetUserId() string {
	cstr := C.kwxApp_GetUserId()
	if cstr == nil {
		return ""
	}
	defer C.kwxApp_FreeString(cstr)
	return C.GoString(cstr)
}

// GetUserName returns the user's full name.
func GetUserName() string {
	cstr := C.kwxApp_GetUserName()
	if cstr == nil {
		return ""
	}
	defer C.kwxApp_FreeString(cstr)
	return C.GoString(cstr)
}

// Event Processing

// Pending returns true if there are pending events in the event queue.
func Pending() bool {
	return C.kwxApp_Pending() != 0
}

// Dispatch dispatches the next pending event.
// Call this in a custom event loop if not using MainLoop().
func Dispatch() {
	C.kwxApp_Dispatch()
}

// Yield yields control to process pending events and then returns.
// Use SafeYield when calling from an event handler to avoid recursion.
// Returns true if events were processed.
func Yield() bool {
	return C.kwxApp_Yield() != 0
}

// SafeYield yields control to process pending events, with protection
// against re-entrant event processing. Pass the window that triggered
// the yield (or nil for the top window).
// Returns true if events were processed.
func SafeYield(window unsafe.Pointer) bool {
	return C.kwxApp_SafeYield(window) != 0
}

// Utilities

// Bell rings the system bell (beep sound).
func Bell() {
	C.kwxApp_Bell()
}

// Sleep pauses execution for the specified number of seconds.
// This does NOT process events - use Yield() if you need event processing.
func Sleep(seconds int) {
	C.kwxApp_Sleep(C.int(seconds))
}

// MilliSleep pauses execution for the specified number of milliseconds.
// This does NOT process events - use Yield() if you need event processing.
func MilliSleep(milliseconds int) {
	C.kwxApp_MilliSleep(C.int(milliseconds))
}

// EnableTooltips enables or disables tooltips globally.
func EnableTooltips(enable bool) {
	cflag := C.int(0)
	if enable {
		cflag = 1
	}
	C.kwxApp_EnableTooltips(cflag)
}

// SetTooltipDelay sets the delay in milliseconds before tooltips are shown.
func SetTooltipDelay(milliseconds int) {
	C.kwxApp_SetTooltipDelay(C.int(milliseconds))
}

// FindWindowById finds a window by its ID.
// If parent is nil, searches globally. Returns nil if not found.
func FindWindowById(id int, parent unsafe.Pointer) unsafe.Pointer {
	return C.kwxApp_FindWindowById(C.int(id), parent)
}

// FindWindowByLabel finds a window by its label/title.
// If parent is nil, searches globally. Returns nil if not found.
func FindWindowByLabel(label string, parent unsafe.Pointer) unsafe.Pointer {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	return C.kwxApp_FindWindowByLabel(clabel, parent)
}

// FindWindowByName finds a window by its internal name.
// If parent is nil, searches globally. Returns nil if not found.
func FindWindowByName(name string, parent unsafe.Pointer) unsafe.Pointer {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.kwxApp_FindWindowByName(cname, parent)
}
