//go:build windows && amd64 && kwxrelease

package wx

// Release-mode cgo linking against a self-contained kwxGO bundle.
//
// Enabled only with -tags kwxrelease so that ordinary development builds,
// which receive their flags from the CMake go_program() helper via
// CGO_CFLAGS/CGO_LDFLAGS, are unaffected.
//
// Expects this layout relative to the bundle root (the parent of wx/):
//   kwxFFI/include/   44 kwxFFI C headers
//   lib/Release/      libkwxall.a
//
// libkwxall.a is the merged archive produced by the CMake target `kwxall`:
// libkwxGO.a + libkwxFFI.a + all wxWidgets static archives, so the linker
// needs no --start-group rescanning of the wx set. Member-level selection is
// preserved, so executable size is unchanged.
//
// NOTE: the blank line below is REQUIRED. cgo feeds every comment that is part
// of the comment group immediately preceding import "C" to the C compiler, so
// without it the prose above would be parsed as C and fail to build.

/*
#cgo CFLAGS:  -I${SRCDIR}/../kwxFFI/include
#cgo LDFLAGS: -L${SRCDIR}/../lib/Release
#cgo LDFLAGS: -Wl,--start-group -lkwxall -Wl,--end-group
#cgo LDFLAGS: -lgdiplus -lmsimg32 -lcomctl32 -luuid -loleacc -lgdi32 -luser32
#cgo LDFLAGS: -lole32 -loleaut32 -luxtheme -lwinspool -lwinmm -lshell32
#cgo LDFLAGS: -lcomdlg32 -ladvapi32 -lwsock32 -lws2_32 -lshlwapi -lversion -lopengl32
#cgo LDFLAGS: -static -lc++ -lunwind
*/
import "C"
