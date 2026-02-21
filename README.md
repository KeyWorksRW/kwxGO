# kwxGO

Go language bindings for [wxWidgets](https://www.wxwidgets.org/) via [kwxFFI](https://github.com/KeyWorksRW/kwxFFI).

kwxFFI exposes a pure C API for wxWidgets. kwxGO wraps that API with idiomatic Go types and provides CGo interop so Go programs can create native GUI applications on Windows, macOS, and Linux without any runtime framework dependency.

---

## Requirements

### All Platforms

| Tool | Minimum Version |
|------|----------------|
| Go   | 1.22           |
| CMake | 3.30          |
| Ninja | any recent    |

### Compiler — **MSVC is not supported**

CGo on Windows uses the **MinGW-w64** linker, which is incompatible with the MSVC object format and ABI. You must use one of:

| Toolchain | Preset | Notes |
|-----------|--------|-------|
| GCC (MinGW-w64) | `ninja-gcc` | Install [MSYS2](https://www.msys2.org/) → `pacman -S mingw-w64-ucrt-x86_64-gcc` |
| Clang (MinGW target) | `ninja-clang` | Install [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw) or MSYS2 → `pacman -S mingw-w64-ucrt-x86_64-clang` |

> **Why no MSVC?** Go's CGo toolchain links using a MinGW-hosted GCC. It cannot consume `.lib` files produced by `cl.exe`. A static library built with GCC or MinGW-targeted Clang produces a `.a` archive that CGo can link directly, with no DLL dependency at runtime.

---

## Building

```powershell
# 1. Configure (choose ninja-gcc or ninja-clang)
cmake --fresh --preset ninja-gcc

# 2. Build Debug
cmake --build build --config Debug

# 3. Build Release
cmake --build build --config Release
```

Output is placed in `bin/Debug/` and `bin/Release/`.

---

## Architecture

```
kwxGO.a          ← static library (kwxApp.cpp)
  └── depends on kwxFFI.a   (pure-C wxWidgets wrapper, fetched via CMake FetchContent)
        └── depends on wxWidgets .a libs (also fetched automatically)
```

CGo in your Go package links all of the above at `go build` time via `cgo LDFLAGS`. There is no shared library required at runtime — the wxWidgets code is statically linked into your final executable.

---

## Project Layout

```
wx/                  # Go package — idiomatic Go API
  ffi.go             # CGo declarations (raw C bridge)
  types.go           # Go structs wrapping C void* pointers
  string.go          # Go ↔ wxString* conversion
  constants.go       # wxWidgets constants
  app.go             # Application lifecycle
  frame.go           # wxFrame wrapper
  ...
examples/
  minimal/           # Minimal window example
src/                 # C++ shim sources (kwxApp.cpp wrapper)
build/               # CMake binary dir (generated, not committed)
bin/                 # Library output (Debug/, Release/)
```

---

## Related Projects

- [kwxFFI](https://github.com/KeyWorksRW/kwxFFI) — C API layer for wxWidgets (upstream)
- [wxWidgets](https://github.com/wxWidgets/wxWidgets) — fetched automatically as a subdependency of kwxFFI
