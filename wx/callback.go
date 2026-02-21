// Package wx callback registry for event handling.
//
// This file implements the Go-side callback registry that maps uintptr IDs
// to Go event handler functions. The //export goEventCallback function is
// the single entry point from C into Go for all event callbacks.
//
// IMPORTANT: This file contains //export directives. Due to cgo rules,
// the C preamble must only contain includes and declarations (no function
// definitions).
//
// Event callback flow:
//
//  1. Go: Bind() registers handler, gets uintptr ID
//  2. Go: calls wxClosure_Create(goEventCallback, id) then wxEvtHandler_Connect
//  3. C++: kwxEventHandler callback fires → goEventCallback(fun, id, evt)
//  4. Go: goEventCallback looks up handler by ID, calls it
//  5. C++: ~wxClosure destructor → goEventCallback(fun, id, nil) → cleanup

package wx

/*
#include <stddef.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

// EventHandler is the callback function type for event handlers.
// The Event parameter provides access to event details. Use the As*Event
// conversion functions to access type-specific data.
//
// Example:
//
//	button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
//	    cmd := wx.AsCommandEvent(evt)
//	    fmt.Println("Selection:", cmd.GetSelection())
//	})
type EventHandler func(event *Event)

// Callback registry — maps uintptr IDs to Go event handler functions.
// This is necessary because cgo cannot pass Go function pointers to C.
// The uintptr ID is used as the void* data parameter in the C closure,
// avoiding cgo pointer rules violations (it's not a real Go pointer).
var (
	callbackRegistry = make(map[uintptr]EventHandler)
	callbackMutex    sync.RWMutex
	nextCallbackID   uintptr
)

// registerCallback stores a Go event handler and returns a unique ID.
// The ID is used as the data pointer in the C closure.
func registerCallback(handler EventHandler) uintptr {
	callbackMutex.Lock()
	defer callbackMutex.Unlock()
	nextCallbackID++
	id := nextCallbackID
	callbackRegistry[id] = handler
	return id
}

// unregisterCallback removes a handler from the registry.
func unregisterCallback(id uintptr) {
	callbackMutex.Lock()
	defer callbackMutex.Unlock()
	delete(callbackRegistry, id)
}

// goEventCallback is the single C-callable entry point for all event callbacks.
// It is called by wxClosure::Invoke when an event fires, and by ~wxClosure
// when the closure is destroyed (with evt == nil as cleanup signal).
//
// Parameters (from C):
//   - fun:  the closure function pointer (unused in Go dispatch)
//   - data: the callback registry ID (uintptr cast to void*)
//   - evt:  the wxEvent pointer, or nil for cleanup
//
//export goEventCallback
func goEventCallback(fun unsafe.Pointer, data unsafe.Pointer, evt unsafe.Pointer) {
	id := uintptr(data)

	if evt == nil {
		// Cleanup call from ~wxClosure destructor.
		// Remove the handler from the registry to prevent memory leaks.
		unregisterCallback(id)
		return
	}

	callbackMutex.RLock()
	handler, ok := callbackRegistry[id]
	callbackMutex.RUnlock()

	if ok {
		handler(&Event{ptr: evt})
	}
}
