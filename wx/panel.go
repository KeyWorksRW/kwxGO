// Package wx panel implementation
//
// This file implements the Panel type, representing a window on which
// controls are placed. Panels are typically children of frames or dialogs.

package wx

/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include

#include <stdlib.h>
#include "wxffi_def.h"
#include "wxffi_types.h"

// Panel functions
extern void* wxPanel_Create(void* parent, int id,
                            int x, int y, int w, int h, int style);
*/
import "C"
import (
	"unsafe"
)

// Panel represents a window on which controls are placed.
// Panels are typically used as containers for controls within frames or dialogs.
// They provide automatic TAB traversal between controls.
type Panel struct {
	BaseWindow
}

// NewPanel creates a new panel.
//
// Parameters:
//   - parent: Parent window (usually a frame or dialog)
//   - id: Window identifier (use ID_ANY for automatic)
//   - pos: Panel position (use DefaultPosition for default)
//   - size: Panel size (use DefaultSize for default)
//   - style: Window style flags (use 0 for default)
//
// Example:
//
//	panel := wx.NewPanel(frame, wx.ID_ANY,
//	    wx.DefaultPosition, wx.DefaultSize, 0)
func NewPanel(parent Window, id int, pos Point, size Size, style int) *Panel {
	var parentPtr unsafe.Pointer
	if parent != nil {
		parentPtr = parent.Ptr()
	}

	ptr := C.wxPanel_Create(
		parentPtr,
		C.int(id),
		C.int(pos.X),
		C.int(pos.Y),
		C.int(size.W),
		C.int(size.H),
		C.int(style),
	)

	panel := &Panel{}
	panel.SetPtr(ptr)
	return panel
}
