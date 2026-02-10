# kwxGO Implementation Plan
Memory saved: kwxgo-plan-v2
## Project Goal

Create a Go binding to the wxWidgets library via kwxFFI, enabling Go developers to build native cross-platform GUI applications.

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
│  │ - C function declarations via kwxFFI headers                 │
│  │ - kwxApp C interface for app lifecycle                       │
│  │ - Type conversions (string ↔ wxString*)                      │
│  │ - Memory management helpers                                  │
│  └──────────────────────────────────────────────────────────────┤
├─────────────────────────────────────────────────────────────────┤
│              kwxFFI (static lib, fetched via CMake)             │
│  - kwxFFI.lib (Windows) / libkwxFFI.a (Linux/macOS)            │
│  - kwxApp.cpp compiled into kwxGO static lib                   │
├─────────────────────────────────────────────────────────────────┤
│                    wxWidgets 3.3.1+                             │
│  - Fetched automatically by kwxFFI's CMakeLists.txt             │
└─────────────────────────────────────────────────────────────────┘
```

### Key Architecture Decisions

1. **FetchContent, not submodules**: kwxFFI and wxWidgets are fetched at CMake configure time via FetchContent (same pattern as kwxFortran). No git submodules.

2. **kwxApp for app lifecycle**: The `kwxApp` pure C interface (from `kwxFFI/examples/CApp/`) handles wxApp initialization. This replaces the old `ELJApp`/`SampleApp` approach. Key functions: `kwxApp_Initialize()`, `kwxApp_MainLoop()`, `kwxApp_ExitMainLoop()`, `kwxApp_Shutdown()`.

3. **Static library default**: kwxFFI builds as a static library by default. The Go build links against `kwxFFI.lib` and `kwxApp.o` plus wxWidgets static libs.

4. **cgo links to build output**: cgo CFLAGS/LDFLAGS point to the CMake build output directory (`build/_deps/kwxffi-src/include` for headers, `build/` for libraries).

---

## Go Toolchain Requirements

### Required Tools

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.21+ | Language toolchain with cgo support |
| **cgo** | (bundled) | C interop - links to kwxFFI |
| **MSVC** | 2022+ | C/C++ compiler (Windows, for building kwxFFI/wxWidgets) |
| **GCC/MinGW** | 13+ | C compiler for cgo on Windows (MinGW-w64) |
| **CMake** | 3.30+ | Build kwxFFI + wxWidgets (CMP0168 for FetchContent perf) |
| **Ninja** | 1.11+ | Build system for kwxFFI |

### Windows Build Notes

**Important**: kwxFFI and wxWidgets are built with MSVC (via CMake/Ninja). The Go side uses cgo with GCC/MinGW to link against the resulting libraries. This means:
- CMake build: MSVC compiler
- cgo compilation: MinGW GCC
- Linking: MinGW links against MSVC-produced `.lib` files

This is a known working pattern — MSVC static libs are compatible with MinGW linking for C-exported symbols.

### Windows-Specific Setup

```powershell
# Option 1: MSYS2 (recommended for MinGW)
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

# cgo flags are set in wx.go via #cgo directives
# No manual env vars needed if build/ directory exists
```

### Linux/macOS Setup

```bash
# Linux (Ubuntu/Debian)
sudo apt install gcc build-essential cmake ninja-build

# macOS
xcode-select --install
brew install cmake ninja

# Verify
go env CGO_ENABLED  # Should be "1"
```

---

## Current State

### Files Already Created
- `go.mod` — Module definition
- `wx/wx.go` — Package doc, cgo linkage, `Run()`, `ExitMainLoop()`, `Window` interface
- `wx/ffi.go` — Raw cgo declarations for string, closure, window, frame, panel, button, statictext, textctrl, sizer, messagebox functions
- `wx/types.go` — `Point`, `Size`, `Rect`, `DefaultPosition`, `DefaultSize`, `BaseWindow`
- `wx/string.go` — `WxString`, `NewWxString()`, `Free()`, `wxStringToGo()`, `WxStringToGoAndFree()`, `withWxString()`
- `examples/minimal/main.go` — Placeholder example (not functional yet)

### Issues in Existing Code

1. **Uses `ELJApp_*` functions**: `wx.go` references `ELJApp_MainLoop`, `ELJApp_ExitMainLoop`, `ELJApp_InitAllImageHandlers`, `ELJApp_Initialized`, `ELJApp_SetTopWindow`. These should be `kwxApp_*` functions from the CApp interface.

2. **Wrong include paths**: cgo directives point to `external/kwxFFI/include` (removed submodule). Must point to FetchContent paths: `build/_deps/kwxffi-src/include` for kwxFFI headers, and `build/_deps/kwxffi-src/examples/CApp` for `kwxApp.h`.

3. **No CMakeLists.txt**: Project has no CMakeLists.txt. Need one modeled on kwxFortran to fetch kwxFFI+wxWidgets and build a static Go helper library containing kwxApp.

4. **Missing `kwxApp_Initialize()` call**: The `Run()` function doesn't call `kwxApp_Initialize()` before entering the main loop. This is required.

5. **kwxApp.cpp not compiled**: The kwxApp source from `examples/CApp/` needs to be compiled and linked. The CMakeLists.txt must handle this (same pattern as kwxFortran's `kwxApp.cpp` inclusion).

---

## Implementation Phases

### Phase 1: CMakeLists.txt and Build System (Complexity: Moderate) ✅

**Objective:** Create CMakeLists.txt using FetchContent to download kwxFFI and wxWidgets, compile kwxApp.cpp, and produce libraries Go can link to.

**Template:** `kwxRust-dev/CMakeLists.txt` (attached in context)

**CMakeLists.txt must:**
1. Require CMake 3.30 (CMP0168 for FetchContent performance)
2. Use FetchContent to fetch kwxFFI from `https://github.com/KeyWorksRW/kwxFFI.git`
3. kwxFFI's own CMakeLists.txt fetches wxWidgets as a subdependency
4. Create a `kwxGO` static library target containing:
   - `examples/CApp/kwxApp.cpp` from the fetched kwxFFI source
5. Link kwxGO against kwxFFI, wxbase, wxcore
6. Set output directories so cgo can find libraries in `bin/<Config>/`
7. Include wxWidgets headers (platform-specific setup.h location)

**Build verification:**
```powershell
cmake -B build -G Ninja
ninja -C build -f build-Debug.ninja
```

**Key paths after build:**
- kwxFFI headers: `build/_deps/kwxffi-src/include/`
- kwxApp header: `build/_deps/kwxffi-src/examples/CApp/kwxApp.h`
- kwxFFI lib: `build/` (kwxFFI.lib or libkwxFFI.a)
- kwxGO lib: `lib/<Config>/` (kwxGO.lib containing kwxApp)
- wxWidgets libs: `build/_deps/wxwidgets-build/lib/`

**Files to create:**
1. `CMakeLists.txt` — Main build file

---

### Phase 2: Fix cgo Linkage and App Lifecycle (Complexity: Moderate) ✅

**Objective:** Update all Go source files to use correct include paths and kwxApp functions.

**Changes to `wx/wx.go`:**
- Replace `#cgo CFLAGS: -I${SRCDIR}/../external/kwxFFI/include` with paths to FetchContent output
- Replace all `ELJApp_*` references with `kwxApp_*` equivalents:
  - `ELJApp_MainLoop` → `kwxApp_MainLoop`
  - `ELJApp_ExitMainLoop` → `kwxApp_ExitMainLoop`
  - `ELJApp_InitAllImageHandlers` → `kwxApp_InitAllImageHandlers`
  - `ELJApp_Initialized` → `kwxApp_IsTerminating` (inverted logic) or remove
  - `ELJApp_SetTopWindow` → `kwxApp_SetTopWindow`
- Add `kwxApp_Initialize()` call in `Run()` before anything else
- Include `kwxApp.h` in the cgo preamble

**Changes to `wx/ffi.go`:**
- Update `#cgo CFLAGS` and `#cgo LDFLAGS` paths
- Replace `ELJApp_*` declarations with `kwxApp_*` equivalents
- Add remaining kwxApp function declarations

**Changes to `wx/string.go`:**
- Update cgo include paths (currently uses relative include which may still work)

**kwxApp C API summary (from kwxApp.h):**
```c
// Lifecycle
int kwxApp_Initialize(int argc, char** argv);
int kwxApp_MainLoop(void);
void kwxApp_ExitMainLoop(void);
void kwxApp_Shutdown(void);
int kwxApp_IsTerminating(void);

// Properties
const char* kwxApp_GetAppName(void);
void kwxApp_SetAppName(const char* name);
const char* kwxApp_GetVendorName(void);
void kwxApp_SetVendorName(const char* name);
void* kwxApp_GetTopWindow(void);
void kwxApp_SetTopWindow(void* window);
void kwxApp_SetExitOnFrameDelete(int flag);

// Utilities
void kwxApp_InitAllImageHandlers(void);
void kwxApp_Bell(void);
void kwxApp_Sleep(int seconds);
void kwxApp_MilliSleep(int milliseconds);
int kwxApp_Yield(void);
int kwxApp_SafeYield(void* window);

// Idle timer
void kwxApp_SetIdleCallback(int interval_ms, kwxApp_IdleCallback callback, void* data);
int kwxApp_GetIdleInterval(void);
```

**Updated `Run()` pattern:**
```go
func Run(init InitFunc) int {
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Initialize wxWidgets via kwxApp
    if C.kwxApp_Initialize(0, nil) == 0 {
        return -1
    }

    // Initialize all image handlers
    C.kwxApp_InitAllImageHandlers()

    // Call user's init function
    if init != nil && !init() {
        C.kwxApp_Shutdown()
        return 0
    }

    // Run the main event loop
    result := int(C.kwxApp_MainLoop())

    // Cleanup
    C.kwxApp_Shutdown()
    return result
}
```

**Verification:**
- `go build ./wx/...` compiles without errors
- cgo can find kwxFFI headers and kwxApp.h

---

### Phase 3: Core Type System Refinement (Complexity: Simple) ✅

**Objective:** Verify and polish type conversions between Go and kwxFFI.

**Already done:** `wx/types.go` and `wx/string.go` exist with core types.

**Remaining:**
- Verify string round-trip works after Phase 2 linkage is fixed
- Add `memory.go` if finalizer-based cleanup is needed
- Unit tests for string conversion (deferred until linkage works)

**Key Conversions (unchanged):**

| Go Type | kwxFFI Type | Conversion |
|---------|------------|------------|
| `string` | `wxString*` | `wxString_CreateUTF8()` / `wxString_GetString()` |
| `int` | `int` | Direct |
| `bool` | `int` (TBool) | 0/1 conversion |
| `unsafe.Pointer` | `void*` | Direct cast |
| `Point{X, Y int}` | `int x, int y` | Struct decomposition |
| `Size{W, H int}` | `int w, int h` | Struct decomposition |
| `Rect{X, Y, W, H}` | 4 ints | Struct decomposition |

---

### Phase 4: Constants and Enums (Complexity: Simple) ✅

**Objective:** Export wxWidgets constants to Go.

**Approach:**
- Constants in kwxFFI are exported as functions (e.g., `expwxDEFAULT_FRAME_STYLE()`)
- Call these once at init and cache the values
- Use Go `var` declarations (not `const` — values come from C at runtime)

**Example:**
```go
// wx/constants.go

var (
    ID_ANY              = int(C.expwxID_ANY())
    ID_OK               = int(C.expwxID_OK())
    ID_CANCEL           = int(C.expwxID_CANCEL())
    DEFAULT_FRAME_STYLE = int(C.expwxDEFAULT_FRAME_STYLE())
    // ... hundreds more
)
```

**Source:** Read `kwxFFI/src/defs.cpp` (5000+ lines of `EXPORT int expXXX()` functions)

**Consideration:** Code generation recommended for 2000+ constants. Can be delegated to a subagent.

**Files to Create:**
1. `wx/constants.go` — All wxWidgets constants

---

### Phase 5: Application Lifecycle (Complexity: Moderate) ✅

**Objective:** Implement wxApp integration for Go using kwxApp.

Already have `Run()` in `wx.go` but it uses wrong API. Phase 2 fixes the API calls.

**Additional app functions to wrap:**
```go
// wx/app.go

func SetAppName(name string) {
    cname := C.CString(name)
    defer C.free(unsafe.Pointer(cname))
    C.kwxApp_SetAppName(cname)
}

func SetVendorName(name string) { ... }
func GetTopWindow() unsafe.Pointer { return C.kwxApp_GetTopWindow() }
func SetExitOnFrameDelete(flag bool) { ... }
func Bell() { C.kwxApp_Bell() }
func Yield() bool { return C.kwxApp_Yield() != 0 }
```

**Thread safety note:** `runtime.LockOSThread()` in `Run()` is critical. All wxWidgets calls must happen on the main OS thread. Go's goroutine scheduler can move goroutines between threads, so the lock prevents that.

---

### Phase 6: Window Hierarchy (Complexity: Moderate) ✅

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
}

// BaseWindow provides common implementation
type BaseWindow struct {
    ptr unsafe.Pointer
}

func (w *BaseWindow) Ptr() unsafe.Pointer { return w.ptr }

func (w *BaseWindow) Show(show bool) {
    flag := C.int(0)
    if show { flag = 1 }
    C.wxWindow_Show(w.ptr, flag)
}
```

**Files to Create:**
1. `wx/object.go` — Base object wrapper
2. `wx/evthandler.go` — Event handler base
3. `wx/window.go` — Window interface and base
4. `wx/frame.go` — Frame implementation
5. `wx/panel.go` — Panel implementation
6. `wx/dialog.go` — Dialog implementation

---

### Phase 7: Event System (Complexity: Complex) ✅

**Objective:** Implement event handling with Go callbacks.

**Challenge:** kwxFFI uses C closures (`wxClosure`) that call back into foreign code. Go's cgo has restrictions on callbacks — Go functions passed to C must use `//export` and be in a separate file.

**Solution:** Use `//export` for a single callback dispatcher and a registry pattern:

```go
// wx/events.go

type EventHandler func(event *Event)

// Global registry for callbacks (cgo limitation)
var callbackRegistry = make(map[uintptr]EventHandler)
var callbackMutex sync.RWMutex
var nextCallbackID uintptr

// Bind connects an event to a handler
func (w *BaseWindow) Bind(eventType int, handler EventHandler) {
    id := registerCallback(handler)
    closure := C.wxClosure_Create(
        C.ClosureFun(C.goEventCallback),
        unsafe.Pointer(id),
    )
    C.wxEvtHandler_Connect(
        w.ptr,
        C.int(-1), C.int(-1),
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
1. `wx/events.go` — Event types and binding
2. `wx/callback.go` — Callback registry and closure helpers (with `//export`)
3. `wx/event_types.go` — Event wrapper structs

---

### Phase 8: Common Controls (Complexity: Moderate) ✅

**Objective:** Wrap frequently-used controls.

**Priority Controls:**
1. Button
2. StaticText (label)
3. TextCtrl (text input)
4. CheckBox
5. RadioButton
6. ComboBox / Choice
7. ListBox

**Example — Button:**
```go
// wx/button.go

type Button struct {
    BaseWindow
}

func NewButton(parent Window, id int, label string, pos Point, size Size, style int) *Button {
    wxLabel := NewWxString(label)
    defer wxLabel.Free()
    ptr := C.wxButton_Create(
        parent.Ptr(),
        C.int(id),
        wxLabel.Ptr(),
        C.int(pos.X), C.int(pos.Y),
        C.int(size.W), C.int(size.H),
        C.int(style),
    )
    return &Button{BaseWindow{ptr: unsafe.Pointer(ptr)}}
}
```

**Files to Create:**
1. `wx/button.go`
2. `wx/statictext.go`
3. `wx/textctrl.go`
4. `wx/checkbox.go`

---

### Phase 9: Layout System (Complexity: Moderate) ✅

**Objective:** Implement sizers for automatic layout.

**Key Sizers:**
- BoxSizer (horizontal/vertical)
- FlexGridSizer
- GridBagSizer
- StaticBoxSizer

**Files to Create:**
1. `wx/sizer.go` (BoxSizer, FlexGridSizer, StaticBoxSizer)

---

### Phase 10: Dialogs (Complexity: Simple) — NOT STARTED

**Objective:** Wrap common dialogs.

**Dialogs:**
- MessageBox
- FileDialog
- DirDialog
- ColorDialog
- FontDialog
- TextEntryDialog

**Files to Create:**
1. `wx/dialogs.go`

---

### Phase 11: Working Minimal Example (Complexity: Simple) — NOT STARTED

**Objective:** Create a working "Hello World" application.

```go
// examples/minimal/main.go
package main

import "github.com/KeyWorksRW/kwxGO/wx"

func main() {
    wx.Run(func() bool {
        frame := wx.NewFrame(nil, wx.ID_ANY, "Hello wxGO",
            wx.DefaultPosition, wx.Size{W: 400, H: 300}, wx.DEFAULT_FRAME_STYLE)

        panel := wx.NewPanel(frame, wx.ID_ANY)
        button := wx.NewButton(panel, wx.ID_ANY, "Click Me!",
            wx.Point{X: 50, Y: 50}, wx.DefaultSize, 0)

        button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
            wx.MessageBox("Hello from Go!", "wxGO", wx.OK|wx.ICON_INFORMATION, frame)
        })

        frame.Show(true)
        wx.SetTopWindow(frame)
        return true
    })
}
```

---

## Build System

### CMakeLists.txt (Phase 1)

Model on kwxFortran's CMakeLists.txt. Key differences:
- Language is CXX only (no Fortran)
- No Go-specific targets — CMake just builds C/C++ libraries
- Go build is separate (`go build`)

### Go Build Process

```powershell
# Step 1: Build C/C++ libraries (one-time, or after kwxFFI updates)
cmake -B build -G Ninja
ninja -C build -f build-Debug.ninja

# Step 2: Build Go package
go build ./wx/...

# Step 3: Build examples
go build -o bin/minimal.exe ./examples/minimal/
```

### cgo Directives (in wx.go)

```go
/*
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/include
#cgo CFLAGS: -I${SRCDIR}/../build/_deps/kwxffi-src/examples/CApp
#cgo windows LDFLAGS: -L${SRCDIR}/../build -lkwxFFI -lkwxGO
#cgo windows LDFLAGS: -L${SRCDIR}/../build/_deps/wxwidgets-build/lib/vc_x64_dll
#cgo windows LDFLAGS: -lwxbase33ud -lwxmsw33ud_core
#cgo linux LDFLAGS: -L${SRCDIR}/../build -lkwxFFI -lkwxGO -lwxbase -lwxcore
#cgo darwin LDFLAGS: -L${SRCDIR}/../build -lkwxFFI -lkwxGO -lwxbase -lwxcore

#include "kwxApp.h"
#include "wxffi_def.h"
#include "wxffi_types.h"
*/
import "C"
```

> **Note:** The exact LDFLAGS for wxWidgets will need tweaking depending on whether wxWidgets is built as shared or static libraries, and the exact library naming. This will be determined during Phase 1 when the CMake build is working.

---

## Risks and Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| MSVC/MinGW ABI compatibility | High | kwxFFI uses `extern "C"` — C ABI is compatible across compilers. Test early in Phase 1. |
| cgo callback limitations | High | Use callback registry pattern with `//export`; test thoroughly in Phase 7 |
| Thread safety (Go scheduler vs main thread) | High | `runtime.LockOSThread()` in `Run()`; document requirements clearly |
| Memory leaks (wxString, objects) | Medium | `WxStringToGoAndFree()` pattern; explicit `Free()` methods; consider finalizers |
| Build complexity (cgo + CMake) | Medium | Provide build scripts; document setup clearly |
| wxWidgets library naming varies by config | Medium | Use CMake generator expressions; detect paths at build time |
| kwxApp.h/cpp not part of kwxFFI lib | Medium | Must compile kwxApp.cpp ourselves (like kwxFortran does) |

---

## Verification Strategy

### Per-Phase Verification
1. **Phase 1**: `ninja -C build -f build-Debug.ninja` succeeds; libraries exist
2. **Phase 2**: `go build ./wx/...` compiles without errors
3. **Phase 3**: Unit tests for string round-trip (when linkage works)
4. **Phase 4**: Compile check — all constants resolve
5. **Phase 5-10**: `go build ./wx/...` compiles; no cgo errors
6. **Phase 11**: Example builds and launches a window

### Integration Testing
- Build and run minimal example on Windows
- Verify event handling works correctly
- Check for memory leaks with Go profiling tools

---

## Subagent Delegation Suitability

### ✅ Suitable for Subagents

| Phase | Agent | Rationale |
|-------|-------|-----------|
| **Phase 4: Constants** | raptor/gpt | Mechanical extraction from defs.cpp. Clear pattern. |
| **Phase 10: Dialogs** | raptor/gpt | Simple wrappers following established patterns. |
| **Phase 8: Controls** | raptor/gpt | Repetitive wrapper code once patterns established. |

### ❌ Requires Local Execution (architect/code agent)

| Phase | Rationale |
|-------|-----------|
| **Phase 1: CMakeLists.txt** | Iterative CMake debugging, build verification |
| **Phase 2: Fix Linkage** | cgo path debugging, compile verification |
| **Phase 5: App Lifecycle** | Runtime verification of thread safety, kwxApp integration |
| **Phase 7: Event System** | Complex cgo callback debugging, memory safety |
| **Phase 11: Minimal Example** | Integration test — requires running GUI |

---

## Estimated Timeline (Agent-Hours)

| Phase | Description | Estimate |
|-------|-------------|----------|
| 1 | CMakeLists.txt + Build System | 30 min |
| 2 | Fix cgo Linkage + kwxApp | 30 min |
| 3 | Core Type System Refinement | 15 min |
| 4 | Constants (delegatable) | 20 min |
| 5 | Application Lifecycle | 20 min |
| 6 | Window Hierarchy | 30 min |
| 7 | Event System | 45 min |
| 8 | Common Controls | 30 min |
| 9 | Layout System | 20 min |
| 10 | Dialogs (delegatable) | 15 min |
| 11 | Minimal Example | 15 min |

**Total: ~4.5 agent-hours** for initial working binding

---

## References

- kwxFFI Architecture: `../.github/kwxffi-architecture.md`
- kwxFFI CApp Interface: `kwxFFI/examples/CApp/kwxApp.h` (fetched at build time)
- kwxFortran CMake template: `../kwxFortran-dev/CMakeLists.txt`
- [Go cgo documentation](https://pkg.go.dev/cmd/cgo)
- [wxWidgets documentation](https://docs.wxwidgets.org/)
