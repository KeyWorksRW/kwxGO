// Package wx provides Go bindings to wxWidgets via wxFFI.
//
// wxWidgets is a cross-platform GUI toolkit that gives applications a native
// look and feel on Windows, macOS, and Linux. This package wraps the wxFFI
// C interface to enable Go applications to use wxWidgets.
//
// # Quick Start
//
// A minimal wxGO application:
//
//	package main
//
//	import "github.com/KeyWorksRW/kwxGO/wx"
//
//	func main() {
//	    wx.Run(func() bool {
//	        frame := wx.NewFrame(nil, wx.ID_ANY, "Hello",
//	            wx.DefaultPosition, wx.Size{W: 400, H: 300}, wx.DEFAULT_FRAME_STYLE)
//	        frame.Show(true)
//	        return true
//	    })
//	}
//
// # Thread Safety
//
// wxWidgets requires all GUI operations to run on the main thread. The wx.Run
// function handles this automatically by calling runtime.LockOSThread().
//
// # Memory Management
//
// Most wxWidgets objects (windows, controls) are owned by their parent and
// deleted automatically. Some objects like wxString must be explicitly freed.
// This package provides helper functions to manage wxString conversion.
package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include -I${SRCDIR}/../build/_deps/kwxffi-src/examples/CApp

#include "kwx_def.h"
#include "kwx_types.h"
#include "kwxApp.h"

// String functions
extern void* wxString_CreateUTF8(char* buffer);
extern void wxString_Delete(void* s);
extern int wxString_GetString(void* s, char* buffer);
extern size_t wxString_Length(void* s);
*/
import "C"
import (
	"runtime"
)

// init ensures the package is properly initialized
func init() {
	// Package initialization happens here
	// Actual wxWidgets init happens in Run()
}

// InitFunc is the application initialization function type.
// Return true to continue running the application, false to exit.
type InitFunc func() bool

// Run starts the wxWidgets application with the given initialization function.
// This MUST be called from main() and will not return until the application exits.
//
// The init function is called after wxWidgets is initialized but before the
// main event loop starts. Create your main window and other UI elements here.
// Return true to enter the main event loop, or false to exit immediately.
//
// Example:
//
//	func main() {
//	    wx.Run(func() bool {
//	        frame := wx.NewFrame(nil, wx.ID_ANY, "My App", ...)
//	        frame.Show(true)
//	        return true
//	    })
//	}
func Run(init InitFunc) int {
	// Lock to OS thread - required for GUI operations
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Initialize wxWidgets via kwxApp
	if C.kwxApp_Initialize(0, nil) == 0 {
		return -1
	}
	defer C.kwxApp_Shutdown()

	// Initialize all image handlers (PNG, JPEG, etc.)
	C.kwxApp_InitAllImageHandlers()

	// Call user's init function
	if init != nil {
		if !init() {
			return 0
		}
	}

	// Run the main event loop
	return int(C.kwxApp_MainLoop())
}

// ExitMainLoop signals the main event loop to exit.
// The application will exit after processing pending events.
func ExitMainLoop() {
	C.kwxApp_ExitMainLoop()
}

// IsTerminating returns true if the application is terminating.
func IsTerminating() bool {
	return C.kwxApp_IsTerminating() != 0
}

// SetTopWindow sets the application's top-level window.
// This is typically called automatically when the first frame is shown.
func SetTopWindow(window Window) {
	if window != nil {
		C.kwxApp_SetTopWindow(window.Ptr())
	}
}
