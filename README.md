# kwxGO

Go language bindings for [wxWidgets](https://www.wxwidgets.org/) via [kwxFFI](https://github.com/KeyWorksRW/kwxFFI), enabling cross-platform GUI applications that render with native controls on Windows, macOS, and Linux.

## ⚠️ Warning: Not Production Ready

**Do not use this Go/wxWidgets interface in production.**

- The kwxFFI ABI interface is **not stable** and can change without warning.
- Very little testing has been done, and won't be until late Q2 of 2026.
- API surface may change as idioms are refined.

## Overview

kwxGO provides a two-layer binding architecture:

1. **FFI Layer** (`wx/cgo_glue_gen.go`) — Raw CGo declarations with `// #include` directives mapping directly to the kwxFFI C function exports. All pointers are `unsafe.Pointer`.
2. **Wrapper Layer** (`wx/*_gen.go`) — Idiomatic Go wrappers with proper types, automatic string conversion, and sensible defaults.

Users should work with the wrapper layer. The FFI layer is available for direct access when needed.

```
Go code  →  CGo (unsafe.Pointer)  →  kwxFFI (C functions)  →  wxWidgets (C++)
```

## Quick Example

```go
// examples/minimal/main.go - Minimal wxWidgets application in Go
package main

import (
	"github.com/KeyWorksRW/kwxGO/wx"
)

func main() {
	wx.Run(func() bool {
		// Create main frame
		frame := wx.NewFrame(nil, wx.ID_ANY, "Hello wxGO",
			-1, -1, 400, 300, wx.DEFAULT_FRAME_STYLE)

		// Create panel to hold controls
		panel := wx.NewPanel(frame.Ptr(), wx.ID_ANY,
			-1, -1, -1, -1, 0)

		// Create a button
		button := wx.NewButton(panel.Ptr(), wx.ID_ANY, "Click Me!",
			50, 50, -1, -1, 0)

		// Bind button click event
		button.Bind(wx.EVT_BUTTON, func(evt *wx.Event) {
			wx.MessageBox("Hello from Go!", "wxGO",
				wx.OK|wx.ICON_INFORMATION, frame.Ptr(), -1, -1)
		})

		// Show the frame
		frame.Show()
		wx.SetTopWindow(frame.Ptr())

		return true
	})
}
```

More examples are in the [`examples/`](examples/) directory.

## Type Hierarchy

Types mirror the wxWidgets C++ class hierarchy using Go struct embedding:

```
BaseObject
└── EvtHandler
    └── Window
        ├── TopLevelWindow
        │   ├── Frame
        │   └── Dialog
        ├── Control
        │   ├── Button
        │   ├── StaticText
        │   ├── TextCtrl
        │   ├── CheckBox
        │   ├── RadioButton
        │   ├── Choice
        │   ├── ListBox
        │   ├── ComboBox
        │   └── BitmapButton
        └── Panel

BaseObject
└── Sizer
    ├── BoxSizer
    ├── FlexGridSizer
    └── GridSizer
```

All types wrap an opaque `unsafe.Pointer` via `BaseObject` with a `Ptr()` accessor.

## Naming Conventions

| Entity | Convention | Example |
| --- | --- | --- |
| Package | `wx` | `import "github.com/KeyWorksRW/kwxGO/wx"` |
| Types | `PascalCase` | `Frame`, `Button`, `Window` |
| Constructors | `New<Type>` | `NewFrame()`, `NewButton()` |
| Methods | `PascalCase` | `frame.SetTitle()`, `btn.SetLabel()` |
| Constants | `UPPER_SNAKE` | `DEFAULT_FRAME_STYLE`, `ID_ANY` |
| Event types | `EVT_*` | `EVT_BUTTON`, `EVT_MENU` |

## Why Go + wxWidgets?

Go's native compilation, goroutine model, and growing use in systems and tooling make it attractive for desktop GUI applications. kwxGO provides native-look cross-platform GUIs with static linking — one executable, no runtime dependencies.

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

## Building

Prerequisites:
- Go 1.22+
- CMake 3.30+
- A supported C/C++ compiler (see above)
- Internet access for FetchContent (kwxFFI and wxWidgets are fetched automatically)

### Consuming a prebuilt release bundle

A released, self-contained bundle needs no `CGO_CFLAGS` / `CGO_LDFLAGS` setup —
just build with the `kwxrelease` tag:

```sh
go build -tags kwxrelease ./your/app
```

See [docs/release-layout.md](docs/release-layout.md) for the bundle tree, ABI/toolchain
locking rules, and the Windows Common Controls v6 manifest caveat.

## License

Apache License 2.0 — see [LICENSE](LICENSE) for details.
