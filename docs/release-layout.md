# Release bundle layout

This document describes the layout of a released, self-contained kwxGO bundle and
how a consumer links against it **without** setting any `CGO_CFLAGS` /
`CGO_LDFLAGS` environment variables by hand.

In-tree development builds are unaffected: they get their flags from the CMake
`go_program()` helper via `CGO_CFLAGS` / `CGO_LDFLAGS`. The directives that make a
bundle self-contained live in a build-tagged file, `wx/cgo_release_windows_amd64.go`,
and are active only when `-tags kwxrelease` is passed.

## Bundle tree

```
kwxgo-bundle/
├── wx/                      Go sources of the `wx` package (from the kwxGO repo)
├── kwxFFI/
│   └── include/             the 44 kwxFFI C headers (kwx_classes.h, kwx_events.h,
│                            kwx_keys.h, kwx_constants.h, kwx_def.h, ...)
└── lib/
    └── Release/
        └── libkwxall.a      ONE merged archive: libkwxGO.a + libkwxFFI.a +
                             every wxWidgets static archive
```

`lib/Release/` deliberately keeps parity with the in-repo CMake output directory,
so a bundle can be produced by copying `wx/` and `kwxFFI/include/` next to the
existing `lib/Release/` build output.

The two directories above are resolved relative to the *package* directory by the
`#cgo` directives:

```go
#cgo CFLAGS:  -I${SRCDIR}/../kwxFFI/include
#cgo LDFLAGS: -L${SRCDIR}/../lib/Release
#cgo LDFLAGS: -Wl,--start-group -lkwxall -Wl,--end-group
```

`${SRCDIR}` is the directory of `wx/`, so the expected bundle root is the parent
of `wx/`. Nothing else in the tree is loadable at link time.

`libkwxall.a` is produced by the CMake target `kwxall`, which merges
`libkwxGO.a` (the `wx.Run()` lifecycle and event trampoline), `libkwxFFI.a` (the C
FFI entry points) and all wxWidgets static archives into a single archive.

## Consumer build

```sh
go build -tags kwxrelease ./your/app
```

Nothing else is required. Without `-tags kwxrelease` the release file is not
compiled at all (`IgnoredGoFiles`), and the package behaves exactly like the
normal in-tree build.

If you would rather not rely on the tag, the equivalent manual form is:

```sh
export CGO_CFLAGS="-I /path/to/bundle/kwxFFI/include"
export CGO_LDFLAGS="-L/path/to/bundle/lib/Release -lkwxall -lgdiplus -lmsimg32 ..."
go build ./your/app
```

The whole point of the bundle file is to make this unnecessary.

## No wxWidgets headers or sources required

Consumers do **not** need the wxWidgets headers, sources or build tree. The kwxFFI
headers are C-safe wrappers: `kwx_classes.h` includes only kwxFFI's own headers,
never `wx/*.h`. Everything C++ lives inside the prebuilt archives.

## ABI and toolchain locking

A bundle is locked to one platform/compiler combination. Building it with a
different toolchain produces link errors or, worse, silent ABI mismatches. Build
one bundle per row, and never mix archives across rows:

| Bundle id               | OS      | Arch  | Compiler / runtime        |
| ----------------------- | ------- | ----- | ------------------------- |
| `windows-amd64-clang-mingw` | Windows | amd64 | clang (MinGW target), libc++ |
| `windows-amd64-gcc-mingw`   | Windows | amd64 | GCC (MinGW-w64), libstdc++   |
| `linux-amd64-gcc`           | Linux   | amd64 | GCC                       |
| `darwin-arm64-clang`        | macOS   | arm64 | Apple clang               |

For the Windows bundles, CGO must use the **MinGW** linker; MSVC is not supported
(Go's external linking cannot consume MSVC `.lib` archives).

Additionally, all three components — the `wx` Go sources, `kwxFFI/include`, and
`libkwxall.a` — must come from the **same commit**. The kwxFFI ABI is not stable,
so a header set from one commit linked against an archive from another may fail to
link or crash at runtime.

## Windows: the Common Controls v6 manifest is still your job

The bundle cannot embed a resource manifest on your behalf, because the manifest
belongs to the final executable, not to a library.

Without it, `wxApp::Initialize()` shows a "no Common Controls v6" warning dialog
and all controls render in the Windows-2000 look. To fix it, ship a resource
script and let `windres` compile it:

```
yourapp.rc  --windres-->  yourapp_windows.syso  --go build-->  embedded in yourapp.exe
```

```sh
windres yourapp.rc -O coff -o yourapp_windows.syso
```

The `.rc` file needs at minimum:

```
#include <windows.h>
1 RT_MANIFEST "yourapp.manifest"
```

`go build` automatically links any `*_windows.syso` file found in the `main`
package directory, so no extra build step or flag is needed. See
`examples/minimal/minimal.rc` for a complete example (manifest plus version info).

## Archive size

`libkwxall.a` is roughly **78 MB**, which is normal for a static wxWidgets build.
Archives of this kind compress roughly **2-2.5x**, so a distributed bundle is
typically 30-40 MB before your own code is added.

## Why one archive is not inefficient

Merging everything into a single archive does not bloat the output executable.

`ar` merges *archive members*, not code: `libkwxall.a` contains exactly the same
object files that the individual archives contained, just indexed together. At
link time `ld` still performs member-level selection — it pulls in only the
members that resolve currently-undefined symbols, and iterates until a fixed
point. Therefore the final `.exe` size is identical to linking the separate
archives; at most a few extra member scan iterations are avoided.

This is also why no `--start-group` rescanning of the wx archive set is needed:
with one archive the inter-member dependency graph is already available to the
linker in a single index. The `-Wl,--start-group` / `--end-group` pair is kept in
the directives anyway because it costs nothing and guards against ordering
surprises.

## Summary

| Concern                    | Handled by                        |
| -------------------------- | --------------------------------- |
| Header path                | `#cgo CFLAGS` in the release file |
| Library path               | `#cgo LDFLAGS -L`                 |
| wx + FFI + binding symbols | `libkwxall.a`                     |
| Windows system libraries   | `#cgo LDFLAGS -l...`              |
| Runtime / C++ runtime      | `-static -lc++ -lunwind`          |
| Common Controls v6 manifest| **consumer** (`windres` + `.syso`)|