package wx

import (
	"runtime"
	"unsafe"
)

// NewIconFromXPM creates a wxIcon from XPM data, given as the lines of an .xpm
// file (no terminating NULL entry needed). The lines are copied into a single
// NUL-terminated blob plus a NULL-terminated pointer table that is valid for the
// duration of the native call; wxWidgets copies the pixels.
func NewIconFromXPM(lines []string) *Icon {
	if len(lines) == 0 {
		return nil
	}

	total := 0
	for _, line := range lines {
		total += len(line) + 1
	}

	blob := make([]byte, 0, total)
	offsets := make([]int, len(lines))
	for i, line := range lines {
		offsets[i] = len(blob)
		blob = append(blob, line...)
		blob = append(blob, 0)
	}

	// NULL-terminated table of C string pointers, holding no Go pointers
	// (a []uintptr is safe to hand to cgo, unlike a []*byte).
	table := make([]uintptr, len(lines)+1)
	for i, offset := range offsets {
		table[i] = uintptr(unsafe.Pointer(&blob[offset]))
	}

	icon := NewIconFromXPMData(unsafe.Pointer(&table[0]))

	runtime.KeepAlive(blob)
	runtime.KeepAlive(table)

	return icon
}
