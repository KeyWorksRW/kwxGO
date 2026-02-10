package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// String functions
extern void* wxString_CreateUTF8(char* buffer);
extern void wxString_Delete(void* s);
extern int wxString_GetString(void* s, char* buffer);
extern size_t wxString_Length(void* s);
*/
import "C"
import "unsafe"

// WxString wraps a wxString pointer for use with wxFFI functions.
// Call Free() when done to release the underlying memory.
type WxString struct {
	ptr unsafe.Pointer
}

// NewWxString creates a new wxString from a Go string.
// The caller must call Free() when done with the string.
//
// Example:
//
//	wxStr := wx.NewWxString("Hello")
//	defer wxStr.Free()
//	// use wxStr.Ptr() in FFI calls
func NewWxString(s string) *WxString {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	ptr := C.wxString_CreateUTF8(cstr)
	return &WxString{ptr: ptr}
}

// Free releases the wxString memory.
// It's safe to call Free() multiple times.
func (w *WxString) Free() {
	if w.ptr != nil {
		C.wxString_Delete(w.ptr)
		w.ptr = nil
	}
}

// Ptr returns the underlying wxString pointer for use in FFI calls.
func (w *WxString) Ptr() unsafe.Pointer {
	return w.ptr
}

// String returns the Go string representation.
// This does NOT free the wxString - call Free() separately.
func (w *WxString) String() string {
	if w.ptr == nil {
		return ""
	}
	return wxStringToGo(w.ptr)
}

// wxStringToGo converts a wxString pointer to a Go string.
// This is a helper for internal use - it does NOT free the wxString.
func wxStringToGo(ptr unsafe.Pointer) string {
	if ptr == nil {
		return ""
	}
	// First call to get length
	length := int(C.wxString_Length(ptr))
	if length == 0 {
		return ""
	}

	// Allocate buffer and get string content
	buf := make([]byte, length)
	C.wxString_GetString(ptr, (*C.char)(unsafe.Pointer(&buf[0])))
	return string(buf)
}

// WxStringToGoAndFree extracts the Go string and frees the wxString.
// Use this for wxFFI functions that return new wxString pointers.
//
// Example:
//
//	wxStrPtr := C.wxTextCtrl_GetValue(ctrl.Ptr())
//	goStr := wx.WxStringToGoAndFree(wxStrPtr)
func WxStringToGoAndFree(ptr unsafe.Pointer) string {
	if ptr == nil {
		return ""
	}
	s := wxStringToGo(ptr)
	C.wxString_Delete(ptr)
	return s
}

// withWxString is a helper that creates a wxString, calls the function with it,
// and automatically frees the wxString. Use for single-use string parameters.
func withWxString(s string, fn func(unsafe.Pointer)) {
	wxStr := NewWxString(s)
	defer wxStr.Free()
	fn(wxStr.Ptr())
}
