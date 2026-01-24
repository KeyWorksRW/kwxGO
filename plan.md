# kwxGO Implementation Plan

## Project Goal

Create a Go binding to the wxWidgets library via wxFFI, enabling Go developers to build native cross-platform GUI applications.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Go Application                          │
├─────────────────────────────────────────────────────────────────┤
│                    kwxGO Package (wx/)                          │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────┐ │
│  │ High-level   │ │ Constants    │ │ Event System             │ │
│  │ OO Wrappers  │ │ (wx.ID_OK)   │ │ (closures, callbacks)    │ │
│  └──────────────┘ └──────────────┘ └──────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┤
│  │ Low-level FFI Bindings (cgo)                                 │
│  │ - C function declarations                                    │
│  │ - Type conversions (string ↔ wxString*)                      │
│  │ - Memory management helpers                                  │
│  └──────────────────────────────────────────────────────────────┤
├─────────────────────────────────────────────────────────────────┤
│                    wxFFI (C DLL/SO)                             │
│  - wxffi.dll (Windows) / libwxffi.so (Linux/macOS)              │
├─────────────────────────────────────────────────────────────────┤
│                    wxWidgets 3.3+                               │
└─────────────────────────────────────────────────────────────────┘
```

---

## Go Toolchain Requirements

### Required Tools

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.21+ | Language toolchain with cgo support |
| **cgo** | (bundled) | C interop - links to wxFFI |
| **GCC/MinGW** | 13+ | C compiler for cgo (Windows: MinGW-w64) |
| **CMake** | 3.20+ | Build wxFFI library |
| **Ninja** | 1.11+ | Build system for wxFFI |

### Windows-Specific Setup

```powershell
# Option 1: MSYS2 (recommended)
pacman -S mingw-w64-ucrt-x86_64-gcc

# Option 2: Pre-built MinGW-w64
# Download from https://github.com/niXman/mingw-builds-binaries/releases

# Verify cgo works
go env CGO_ENABLED  # Should be "1"
```

### Environment Variables

```powershell
# Windows
$env:CGO_ENABLED = "1"
$env:CC = "gcc"  # or path to MinGW gcc

# For linking wxFFI
$env:CGO_LDFLAGS = "-L/path/to/wxffi/build -lwxffi"
$env:CGO_CFLAGS = "-I/path/to/wxffi/include"
```

### Linux/macOS Setup

```bash
# Linux (Ubuntu/Debian)
sudo apt install gcc build-essential

# macOS
xcode-select --install

# Verify
go env CGO_ENABLED  # Should be "1"
```

---

## Implementation Phases

### Phase 1: Project Foundation (Complexity: Simple)

**Objective:** Set up Go module structure and verify cgo compilation.

**Directory Structure:**
```
kwxGO-dev/
├── go.mod                 # Module definition
├── wx/                    # Main package
│   ├── wx.go             # Package doc, initialization
│   ├── ffi.go            # Low-level cgo bindings
│   ├── types.go          # Go type definitions
│   ├── string.go         # wxString conversion helpers
│   └── constants.go      # wxWidgets constants
├── examples/
│   └── minimal/
│       └── main.go       # Minimal app example
├── build/                 # Build artifacts
├── external/
│   └── wxFFI/            # Submodule (existing)
└── CMakeLists.txt        # Build wxFFI + examples
```

**Files to Create:**
1. `go.mod` - Module initialization
2. `wx/wx.go` - Package initialization and cgo linkage
3. `wx/ffi.go` - Raw cgo function declarations

**Verification:**
- `go build ./wx/...` compiles without errors
- cgo can find wxFFI headers

---

### Phase 2: Core Type System (Complexity: Moderate)

**Objective:** Implement type conversions between Go and wxFFI.

**Key Conversions:**

| Go Type | wxFFI Type | Conversion |
|---------|------------|------------|
| `string` | `wxString*` | `wxString_CreateUTF8()` / `wxString_GetString()` |
| `int` | `int` | Direct |
| `bool` | `int` (TBool) | 0/1 conversion |
| `unsafe.Pointer` | `void*` | Direct cast |
| `Point{X, Y int}` | `int x, int y` | Struct decomposition |
| `Size{W, H int}` | `int w, int h` | Struct decomposition |
| `Rect{X, Y, W, H}` | 4 ints | Struct decomposition |

**String Handling (Critical):**

```go
// wx/string.go

// WxString wraps a wxString* with automatic cleanup
type WxString struct {
    ptr unsafe.Pointer
}

// NewWxString creates a wxString from Go string
func NewWxString(s string) *WxString {
    cstr := C.CString(s)
    defer C.free(unsafe.Pointer(cstr))
    return &WxString{ptr: C.wxString_CreateUTF8(cstr)}
}

// Free releases the wxString memory
func (w *WxString) Free() {
    if w.ptr != nil {
        C.wxString_Delete((*C.wxString)(w.ptr))
        w.ptr = nil
    }
}

// GoString extracts the string and frees the wxString
func WxStringToGo(ptr unsafe.Pointer) string {
    if ptr == nil {
        return ""
    }
    wxStr := (*C.wxString)(ptr)
    length := C.wxString_Length(wxStr)
    buf := make([]byte, length)
    C.wxString_GetString(wxStr, (*C.char)(unsafe.Pointer(&buf[0])))
    C.wxString_Delete(wxStr)
    return string(buf)
}
```

**Files to Create:**
1. `wx/types.go` - Go struct definitions (Point, Size, Rect, etc.)
2. `wx/string.go` - wxString conversion helpers
3. `wx/memory.go` - Memory management utilities

**Verification:**
- Unit tests for string round-trip conversion
- No memory leaks (use Go's race detector)

---

### Phase 3: Constants and Enums (Complexity: Simple)

**Objective:** Export wxWidgets constants to Go.

**Approach:**
- Constants in wxFFI are exported as functions (e.g., `expwxDEFAULT_FRAME_STYLE()`)
- Call these once at init and cache the values
- Use Go constants where possible for compile-time checks

**Example:**
```go
// wx/constants.go

var (
    ID_ANY           = int(C.expwxID_ANY())
    ID_OK            = int(C.expwxID_OK())
    ID_CANCEL        = int(C.expwxID_CANCEL())
    DEFAULT_FRAME_STYLE = int(C.expwxDEFAULT_FRAME_STYLE())
    // ... hundreds more
)
```

**Files to Create:**
1. `wx/constants.go` - All wxWidgets constants

**Consideration:** May need code generation for 2000+ constants.

---

### Phase 4: Application Lifecycle (Complexity: Moderate)

**Objective:** Implement wxApp integration for Go.

**Challenge:** wxWidgets requires running on the main thread, and Go's runtime manages threads differently.

**Solution Pattern:**
```go
// wx/app.go

// App represents the wxWidgets application
type App struct {
    ptr unsafe.Pointer
}

// Run starts the application with the given init function
// This MUST be called from main()
func Run(init func() bool) int {
    // Lock to main thread (required for GUI)
    runtime.LockOSThread()

    // Create closure for init callback
    initClosure := createClosure(init)

    // Call wxEntry or equivalent
    return int(C.ELJApp_MainLoop())
}
```

**Key Functions to Wrap:**
- `ELJApp` initialization
- Main event loop
- Application termination

**Files to Create:**
1. `wx/app.go` - Application lifecycle management

---

### Phase 5: Window Hierarchy (Complexity: Moderate)

**Objective:** Implement base window types.

**Class Hierarchy:**
```
Object (base)
└── EvtHandler
    └── Window
        ├── Panel
        ├── TopLevelWindow
        │   ├── Frame
        │   └── Dialog
        └── Control (base for all controls)
```

**Go Interface Pattern:**
```go
// wx/window.go

// Window interface represents any wxWindow
type Window interface {
    Ptr() unsafe.Pointer
    Show(show bool)
    Hide()
    Destroy()
    SetSize(x, y, w, h int)
    // ... common methods
}

// BaseWindow provides common implementation
type BaseWindow struct {
    ptr unsafe.Pointer
}

func (w *BaseWindow) Ptr() unsafe.Pointer { return w.ptr }

func (w *BaseWindow) Show(show bool) {
    flag := 0
    if show { flag = 1 }
    C.wxWindow_Show((*C.wxWindow)(w.ptr), C.int(flag))
}
```

**Files to Create:**
1. `wx/object.go` - Base object wrapper
2. `wx/evthandler.go` - Event handler base
3. `wx/window.go` - Window interface and base
4. `wx/frame.go` - Frame implementation
5. `wx/panel.go` - Panel implementation
6. `wx/dialog.go` - Dialog implementation

---

### Phase 6: Event System (Complexity: Complex)

**Objective:** Implement event handling with Go callbacks.

**Challenge:** wxFFI uses C closures that call back into foreign code. Go's cgo has restrictions on callbacks.

**Solution:** Use `//export` for callback functions and a registry pattern.

```go
// wx/events.go

// Event types (loaded from wxFFI)
var (
    EVT_BUTTON     = int(C.expEVT_COMMAND_BUTTON_CLICKED())
    EVT_CLOSE      = int(C.expEVT_CLOSE_WINDOW())
    EVT_MENU       = int(C.expEVT_MENU())
    // ...
)

// EventHandler is a Go function that handles events
type EventHandler func(event *Event)

// Global registry for callbacks (cgo limitation)
var callbackRegistry = make(map[uintptr]EventHandler)
var callbackMutex sync.RWMutex

// Bind connects an event to a handler
func (w *BaseWindow) Bind(eventType int, handler EventHandler) {
    // Create closure wrapper
    closure := C.wxClosure_Create(
        C.ClosureFun(C.goEventCallback),
        unsafe.Pointer(uintptr(registerCallback(handler))),
    )
    C.wxEvtHandler_Connect(
        (*C.wxEvtHandler)(w.ptr),
        C.int(-1),  // ID
        C.int(-1),  // lastID
        C.int(eventType),
        closure,
    )
}

//export goEventCallback
func goEventCallback(fun, data, evt unsafe.Pointer) {
    id := uintptr(data)
    callbackMutex.RLock()
    handler, ok := callbackRegistry[id]
    callbackMutex.RUnlock()
    if ok && evt != nil {
        handler(&Event{ptr: evt})
    }
}
```

**Files to Create:**
1. `wx/events.go` - Event types and binding
2. `wx/callback.go` - Callback registry and closure helpers
3. `wx/event_types.go` - Event wrapper structs

---

### Phase 7: Common Controls (Complexity: Moderate)

**Objective:** Wrap frequently-used controls.

**Priority Controls:**
1. Button
2. StaticText (label)
3. TextCtrl (text input)
4. CheckBox
5. RadioButton
6. ComboBox / Choice
7. ListBox
8. Panel
9. Sizer (layout)

**Example - Button:**
```go
// wx/button.go

type Button struct {
    BaseWindow
}

func NewButton(parent Window, id int, label string, pos Point, size Size, style int) *Button {
    wxLabel := NewWxString(label)
    defer wxLabel.Free()

    ptr := C.wxButton_Create(
        (*C.wxWindow)(parent.Ptr()),
        C.int(id),
        (*C.wxString)(wxLabel.ptr),
        C.int(pos.X), C.int(pos.Y),
        C.int(size.W), C.int(size.H),
        C.int(style),
    )
    return &Button{BaseWindow{ptr: unsafe.Pointer(ptr)}}
}

func (b *Button) SetLabel(label string) {
    wxLabel := NewWxString(label)
    defer wxLabel.Free()
    C.wxButton_SetLabel((*C.wxButton)(b.ptr), (*C.wxString)(wxLabel.ptr))
}

func (b *Button) SetDefault() {
    C.wxButton_SetDefault((*C.wxButton)(b.ptr))
}
```

**Files to Create:**
1. `wx/button.go`
2. `wx/statictext.go`
3. `wx/textctrl.go`
4. `wx/checkbox.go`
5. `wx/sizer.go` (BoxSizer, FlexGridSizer)

---

### Phase 8: Layout System (Complexity: Moderate)

**Objective:** Implement sizers for automatic layout.

**Key Sizers:**
- BoxSizer (horizontal/vertical)
- FlexGridSizer
- GridBagSizer
- StaticBoxSizer

```go
// wx/sizer.go

type Sizer interface {
    Ptr() unsafe.Pointer
    Add(item interface{}, proportion int, flag int, border int)
    AddSpacer(size int)
    Fit(window Window)
}

type BoxSizer struct {
    ptr unsafe.Pointer
}

func NewBoxSizer(orient int) *BoxSizer {
    return &BoxSizer{ptr: unsafe.Pointer(C.wxBoxSizer_Create(C.int(orient)))}
}

func (s *BoxSizer) Add(window Window, proportion, flag, border int) {
    C.wxSizer_Add(
        (*C.wxSizer)(s.ptr),
        (*C.wxWindow)(window.Ptr()),
        C.int(proportion),
        C.int(flag),
        C.int(border),
    )
}
```

---

### Phase 9: Dialogs (Complexity: Simple)

**Objective:** Wrap common dialogs.

**Dialogs:**
- MessageBox
- FileDialog
- DirDialog
- ColorDialog
- FontDialog
- TextEntryDialog

```go
// wx/dialogs.go

func MessageBox(message, caption string, style int, parent Window) int {
    wxMsg := NewWxString(message)
    defer wxMsg.Free()
    wxCap := NewWxString(caption)
    defer wxCap.Free()

    var parentPtr *C.wxWindow
    if parent != nil {
        parentPtr = (*C.wxWindow)(parent.Ptr())
    }

    return int(C.wxMessageBox(
        (*C.wxString)(wxMsg.ptr),
        (*C.wxString)(wxCap.ptr),
        C.int(style),
        parentPtr,
    ))
}
```

---

### Phase 10: Minimal Example (Complexity: Simple)

**Objective:** Create a working "Hello World" application.

```go
// examples/minimal/main.go

package main

import (
    "github.com/KeyWorksRW/kwxGO/wx"
)

type MyFrame struct {
    *wx.Frame
}

func NewMyFrame() *MyFrame {
    frame := wx.NewFrame(nil, wx.ID_ANY, "Hello wxGO",
        wx.DefaultPosition, wx.Size{W: 400, H: 300}, wx.DEFAULT_FRAME_STYLE)

    panel := wx.NewPanel(frame, wx.ID_ANY)
    button := wx.NewButton(panel, wx.ID_ANY, "Click Me!",
        wx.Point{X: 50, Y: 50}, wx.DefaultSize, 0)

    button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
        wx.MessageBox("Hello from Go!", "wxGO", wx.OK|wx.ICON_INFORMATION, frame)
    })

    frame.Bind(wx.EVT_CLOSE, func(evt *wx.Event) {
        frame.Destroy()
    })

    return &MyFrame{Frame: frame}
}

func main() {
    wx.Run(func() bool {
        frame := NewMyFrame()
        frame.Show(true)
        return true
    })
}
```

---

## Build System

### CMakeLists.txt Updates

```cmake
# Add to existing CMakeLists.txt

# Build wxFFI as shared library for Go
add_subdirectory(external/wxFFI)

# Set output directory for DLL/SO
set_target_properties(kwxFFI PROPERTIES
    LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib
)
```

### Go Build Script (build.ps1)

```powershell
# build.ps1 - Build script for Windows

param(
    [switch]$Release,
    [switch]$Clean
)

$ErrorActionPreference = "Stop"
$ProjectRoot = $PSScriptRoot

# Build wxFFI first
Write-Host "Building wxFFI..." -ForegroundColor Cyan
Push-Location "$ProjectRoot/build"
if ($Release) {
    ninja -f build-Release.ninja
} else {
    ninja -f build-Debug.ninja
}
Pop-Location

# Set cgo environment
$env:CGO_ENABLED = "1"
$env:CGO_CFLAGS = "-I$ProjectRoot/external/wxFFI/include"
$env:CGO_LDFLAGS = "-L$ProjectRoot/build/lib -lwxffi"

# Build Go package
Write-Host "Building Go package..." -ForegroundColor Cyan
go build ./wx/...

# Build examples
Write-Host "Building examples..." -ForegroundColor Cyan
go build -o build/minimal.exe ./examples/minimal/
```

---

## Risks and Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| cgo callback limitations | High | Use callback registry pattern; test thoroughly |
| Thread safety (Go scheduler vs main thread) | High | `runtime.LockOSThread()` in main; document requirements |
| Memory leaks (wxString, objects) | Medium | Implement finalizers; provide explicit Free() methods |
| Build complexity (cgo + cmake) | Medium | Provide build scripts; document setup clearly |
| Platform differences (Windows/Linux/macOS) | Medium | CI testing on all platforms; conditional compilation |
| wxFFI API changes | Low | Pin to specific wxFFI version; test against updates |

---

## Verification Strategy

### Per-Phase Verification
1. **Compile test:** `go build ./wx/...`
2. **Unit tests:** `go test ./wx/...`
3. **Race detection:** `go test -race ./wx/...`
4. **Example runs:** Manual testing of example apps

### Integration Testing
- Build and run minimal example on Windows
- Verify event handling works correctly
- Check for memory leaks with profiling

### CI Pipeline (Future)
- GitHub Actions workflow for Windows/Linux/macOS
- Automated build of wxFFI + Go bindings
- Run tests on each platform

---

## Estimated Timeline

| Phase | Description | Estimated Effort |
|-------|-------------|------------------|
| 1 | Project Foundation | 2 hours |
| 2 | Core Type System | 4 hours |
| 3 | Constants | 2 hours (may need codegen) |
| 4 | App Lifecycle | 4 hours |
| 5 | Window Hierarchy | 6 hours |
| 6 | Event System | 8 hours (complex) |
| 7 | Common Controls | 6 hours |
| 8 | Layout System | 4 hours |
| 9 | Dialogs | 2 hours |
| 10 | Minimal Example | 2 hours |

**Total: ~40 hours** (for initial working binding)

---

## References

- [wxFFI Architecture](.github/wxffi-architecture.md)
- [Go cgo documentation](https://pkg.go.dev/cmd/cgo)
- [wxWidgets documentation](https://docs.wxwidgets.org/)
- [Similar project: andlabs/ui](https://github.com/andlabs/ui) (Go GUI, different approach)

---

## GitHub.com Coding Agent Suitability

The following evaluates which phases could be assigned to the GitHub.com coding agent (Opus 4.5) via an issue. The agent lacks Serena symbolic tools and cannot build/test locally, so phases requiring iterative debugging or runtime verification are poor candidates.

### ✅ Suitable for GitHub.com Agent

| Phase | Rationale |
|-------|-----------|
| **Phase 3: Constants** | Mechanical extraction of constants from wxFFI. Clear pattern: wrap `expwxXXX()` calls. No runtime testing needed—just generate the Go file. Agent can read `external/wxFFI/src/defs.cpp` and produce `wx/constants.go`. |
| **Phase 9: Dialogs** | Simple wrappers following established patterns from Phase 7. Once Button/TextCtrl patterns exist, dialogs are straightforward. Agent can generate `wx/dialogs.go` with MessageBox, FileDialog, etc. |

### ⚠️ Possibly Suitable (with constraints)

| Phase | Rationale |
|-------|-----------|
| **Phase 5: Window Hierarchy** | Repetitive wrapper code following a template. However, depends on Phases 2/4 being complete and tested. Could work if given explicit templates and the completed `wx/types.go`. |
| **Phase 7: Common Controls** | Similar pattern repetition (Button, TextCtrl, CheckBox, etc.). Depends on event system being functional. Could be split: agent generates initial code, human verifies against runtime. |
| **Phase 8: Layout System** | Sizer wrappers are mostly mechanical. Depends on Phase 5 window types. Lower risk than controls since sizers have simpler interactions. |

### ❌ Not Suitable for GitHub.com Agent

| Phase | Rationale |
|-------|-----------|
| **Phase 1: Project Foundation** | Requires iterative cgo debugging, environment setup verification, and confirming build actually works. Agent cannot run `go build` or verify cgo linkage. |
| **Phase 2: Core Type System** | String conversion is critical and subtle. Memory management bugs only manifest at runtime. Needs interactive testing with actual wxFFI library. |
| **Phase 4: App Lifecycle** | Thread safety and `runtime.LockOSThread()` behavior requires runtime verification. Main loop integration is highly platform-specific. |
| **Phase 6: Event System** | Most complex phase. cgo callback restrictions, closure registry, and event dispatch all require iterative debugging. Memory safety issues only visible at runtime. |
| **Phase 10: Minimal Example** | Integration test—requires all prior phases working together. Must be run and visually verified. |

### Recommended Approach

1. **Complete locally:** Phases 1, 2, 4, 6, 10 (foundation + complex runtime-dependent work)
2. **Assign to agent:** Phase 3 (constants) can be done early in parallel
3. **Assign to agent after templates exist:** Phases 5, 7, 8, 9 once patterns are established

### Example Issue for Phase 3

```markdown
## Generate wx/constants.go from wxFFI

**Context:** This project binds wxWidgets to Go via wxFFI. Constants in wxFFI
are exported as C functions (e.g., `expwxID_OK()` returns the value of `wxID_OK`).

**Task:** Generate `wx/constants.go` that:
1. Declares cgo imports for constant functions from wxFFI
2. Exports Go variables for each constant (e.g., `var ID_OK = int(C.expwxID_OK())`)

**Source:** Read `external/wxFFI/src/defs.cpp` for the list of `EXPORT int expwxXXX()` functions.

**Pattern:**
```go
// #include "wxffi_def.h"
// int expwxID_OK();
import "C"

var ID_OK = int(C.expwxID_OK())
```

**Scope:** Generate constants for IDs, styles, and flags. Skip string constants (wxART_*) for now.
```
