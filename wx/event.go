// Package wx — event binding for kwxGO.
//
// Bridges Go closures to the wxWidgets event system via the C callback
// mechanism in src/go_events.cpp.
//
// Flow:
//  1. Run() calls initEventCallback() which registers goEventDispatch
//     as the global C event callback via kwxGO_SetEventCallback.
//  2. Bind/BindId stores the Go handler in handlerRegistry keyed by a
//     runtime/cgo.Handle, then calls kwxGO_Connect with the handle value.
//  3. When an event fires, C calls goEventDispatch(fun, data, evt).
//     We recover the handle from data and call the stored Go closure.
package wx

/*
#include <stdint.h>

// Callback type and registration (src/go_events.cpp)
typedef void (*GoEventCallbackFunc)(void*, void*, void*);
extern void kwxGO_SetEventCallback(GoEventCallbackFunc callback);
extern int  kwxGO_Connect(void *obj, int first, int last, int type, void *data);
extern int  kwxGO_Disconnect(void *obj, int first, int last, int type);

// Forward-declare the Go-exported dispatcher with matching export attributes
// so CGo's generated redeclaration in _cgo_export.c is identical and harmless.
#if defined(_WIN32) || defined(__CYGWIN__)
extern __declspec(dllexport) void goEventDispatch(void *fun, void *data, void *evt);
#else
extern void __attribute__((visibility("default"))) goEventDispatch(void *fun, void *data, void *evt);
#endif

// C-level trampoline so Go code never needs to form C.goEventDispatch.
static void initGoEventCallback(void) {
    kwxGO_SetEventCallback(goEventDispatch);
}

// Clear the callback so no events dispatch into Go during shutdown.
static void clearGoEventCallback(void) {
    kwxGO_SetEventCallback(NULL);
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// EventHandler is the Go callback signature for wx events.
// The *Event is only valid for the duration of the handler call.
type EventHandler func(evt *Event)

// initEventCallback registers goEventDispatch with the C layer.
// Called once by Run() before the event loop starts.
func initEventCallback() {
	C.initGoEventCallback()
}

// clearEventCallback removes the Go dispatch function so no events
// are routed into Go during wxWidgets shutdown / cleanup.
func clearEventCallback() {
	C.clearGoEventCallback()
}

//export goEventDispatch
func goEventDispatch(_ unsafe.Pointer, data unsafe.Pointer, evt unsafe.Pointer) {
	h := cgo.Handle(uintptr(data))
	if fn, ok := h.Value().(EventHandler); ok {
		fn(&Event{BaseObject: BaseObject{ptr: evt}})
	}
}

// Bind connects a Go handler to an event type on this window/control.
// The handler fires whenever an event of the given type originates from
// this object (any ID).
func (o *EvtHandler) Bind(evtType int, fn EventHandler) {
	h := cgo.NewHandle(fn)
	C.kwxGO_Connect(o.Ptr(), -1, -1, C.int(evtType), unsafe.Pointer(uintptr(h)))
}

// BindId connects a Go handler to an event type for a specific command ID.
// Used for menu items, buttons with known IDs, etc.
func (o *EvtHandler) BindId(evtType int, id int, fn EventHandler) {
	h := cgo.NewHandle(fn)
	C.kwxGO_Connect(o.Ptr(), C.int(id), C.int(id), C.int(evtType), unsafe.Pointer(uintptr(h)))
}

// Unbind disconnects all handlers for the given event type on this object.
func (o *EvtHandler) Unbind(evtType int) {
	C.kwxGO_Disconnect(o.Ptr(), -1, -1, C.int(evtType))
}

// UnbindId disconnects all handlers for the given event type and ID.
func (o *EvtHandler) UnbindId(evtType int, id int) {
	C.kwxGO_Disconnect(o.Ptr(), C.int(id), C.int(id), C.int(evtType))
}
