// Package wx memory management patterns
//
// This file documents memory management approaches for kwxGO bindings.

package wx

// Memory Management Overview
//
// kwxGO uses a mix of automatic and manual memory management:
//
// # 1. wxWidgets Objects (Windows, Controls, etc.)
//
// Most wxWidgets objects are automatically managed by wxWidgets' parent-child
// ownership model. When you create a window with a parent, wxWidgets takes
// ownership and will delete it when the parent is destroyed.
//
//	frame := wx.NewFrame(nil, wx.ID_ANY, "Title", ...)  // Top-level, user destroys
//	button := wx.NewButton(frame, wx.ID_OK, "OK", ...)  // Owned by frame, auto-deleted
//
// You only need to explicitly destroy top-level windows (frames, dialogs) when
// closing them programmatically. Child windows are cleaned up automatically.
//
// # 2. wxString Objects - MANUAL MANAGEMENT REQUIRED
//
// wxString objects created on the Go side MUST be explicitly freed using Free().
// This is because:
//   - They are allocated by kwxFFI in C++ heap memory
//   - Go's garbage collector cannot track them
//   - cgo does not allow safe finalizers for C++ objects
//
// Always use defer to ensure cleanup:
//
//	wxStr := wx.NewWxString("Hello")
//	defer wxStr.Free()  // REQUIRED
//	// use wxStr.Ptr() in FFI calls
//
// Helper patterns for convenience:
//
//	// Pattern 1: One-use strings (automatic cleanup)
//	withWxString("Hello", func(ptr unsafe.Pointer) {
//	    C.wxSomeFunction(ptr)
//	})
//
//	// Pattern 2: Extract and free in one call
//	wxStrPtr := C.wxTextCtrl_GetValue(ctrl.Ptr())
//	goStr := wx.WxStringToGoAndFree(wxStrPtr)  // Extracts and frees
//
// # 3. Event Closures and Callbacks
//
// Event handler closures are managed by the event system and cleaned up when
// the window is destroyed. You don't need to manually free them.
//
// # 4. Finalizers - NOT USED
//
// Go's runtime.SetFinalizer is NOT used for wxString or other C++ objects because:
//   - Finalizers run in a separate goroutine at unpredictable times
//   - wxWidgets requires all operations on the main thread (runtime.LockOSThread)
//   - cgo calls from finalizers can cause deadlocks
//   - Finalizers are not guaranteed to run before program exit
//
// For these reasons, kwxGO uses explicit Free() calls (RAII-like via defer).
//
// # 5. Thread Safety
//
// ALL wxWidgets operations must happen on the main OS thread. The wx.Run()
// function calls runtime.LockOSThread() to ensure this. Do not create
// goroutines that directly call wxWidgets functions - use wx.CallAfter()
// to marshal calls to the main thread (Phase 9 in implementation plan).
