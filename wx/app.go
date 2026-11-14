// Package wx — app lifecycle wrappers for kwxGO.
//
// Calls kwxApp_* functions compiled into libkwxGO.a (src/kwxApp.cpp).
package wx

/*
#include <stdlib.h>

extern int  kwxApp_Initialize(int argc, char **argv);
extern int  kwxApp_MainLoop(void);
extern void kwxApp_ExitMainLoop(void);
extern void kwxApp_Shutdown(void);
extern void kwxApp_SetAppName(const char *name);
extern const char *kwxApp_GetAppName(void);
extern void kwxApp_SetVendorName(const char *name);
extern void kwxApp_SetTopWindow(void *window);
extern void *kwxApp_GetTopWindow(void);
extern void kwxApp_SetExitOnFrameDelete(int flag);
extern int  kwxApp_Yield(void);
extern int  kwxApp_SafeYield(void *window);
extern void kwxApp_InitAllImageHandlers(void);
extern void kwxApp_Bell(void);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Run initialises wxWidgets, calls onInit (which should create the top-level
// frame and return true), then blocks in the wx main event loop.
// It is the idiomatic entry point for a kwxGO application.
func Run(onInit func() bool) {
	// wxWidgets requires every GUI call to happen on the OS thread that
	// initialised the toolkit (wxThread::IsMain()). Between cgo calls the Go
	// scheduler is free to move the main goroutine to another thread, which
	// makes wx assert ("only main thread may call
	// wxMutexGuiLeaveOrEnter()") or hang before the first window appears.
	// Pin the goroutine for the whole lifetime of the event loop.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Bootstrap wx with no real argc/argv — kwxApp_Initialize accepts 0/nil.
	if C.kwxApp_Initialize(0, nil) == 0 {
		panic("kwxGO: kwxApp_Initialize failed")
	}
	// Register the global Go event callback so Bind/BindId work.
	initEventCallback()

	if onInit != nil && !onInit() {
		C.kwxApp_Shutdown()
		return
	}
	C.kwxApp_MainLoop()
	clearEventCallback() // Stop dispatching events before cleanup
	C.kwxApp_Shutdown()
}

// ExitMainLoop exits the wx event loop. Call from an event handler to quit.
func ExitMainLoop() { C.kwxApp_ExitMainLoop() }

// SetAppName sets the application name (used in dialogs, config files, etc.).
func SetAppName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.kwxApp_SetAppName(cs)
}

// GetAppName returns the current application name.
func GetAppName() string { return C.GoString(C.kwxApp_GetAppName()) }

// SetVendorName sets the vendor / organisation name.
func SetVendorName(name string) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	C.kwxApp_SetVendorName(cs)
}

// SetTopWindow tells wx which window is the "main" window.
// Pass frame.Ptr() as the argument.
func SetTopWindow(win unsafe.Pointer) { C.kwxApp_SetTopWindow(win) }

// GetTopWindow returns the top-level window pointer set via SetTopWindow.
func GetTopWindow() unsafe.Pointer { return C.kwxApp_GetTopWindow() }

// SetExitOnFrameDelete controls whether the app exits when the last frame
// is deleted (default: true).
func SetExitOnFrameDelete(exit bool) {
	flag := 0
	if exit {
		flag = 1
	}
	C.kwxApp_SetExitOnFrameDelete(C.int(flag))
}

// Yield processes pending events. Safe to call from non-event context.
func Yield() bool { return C.kwxApp_Yield() != 0 }

// SafeYield processes pending events safely from within an event handler.
// Pass the window that should remain interactive, or nil for any window.
func SafeYield(win unsafe.Pointer) bool { return C.kwxApp_SafeYield(win) != 0 }

// InitAllImageHandlers initialises PNG, JPEG, GIF, and other image handlers.
func InitAllImageHandlers() { C.kwxApp_InitAllImageHandlers() }

// Bell sounds the system bell.
func Bell() { C.kwxApp_Bell() }
