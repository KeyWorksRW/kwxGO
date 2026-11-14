# Windows Manifest Files for kwxGO

This directory contains Windows-specific resource files that embed an application manifest with Common Controls v6 support.

## Files

- **minimal.manifest**: XML manifest declaring dependencies and settings
  - Common Controls v6 for modern themed controls
  - DPI awareness (PerMonitorV2)
  - OS compatibility (Windows 7-11)
  - Execution level (asInvoker)

- **minimal.rc**: Resource script referencing the manifest
  - Embeds manifest as RT_MANIFEST resource
  - Includes version information
  - Can be extended with icons

- **minimal_windows.syso**: Compiled resource object (generated, not checked in)
  - Created by `windres -i minimal.rc -o minimal_windows.syso -O coff`
  - The `_windows` suffix is a GOOS build constraint: go build links it only on
    Windows and ignores it everywhere else
  - Automatically linked by Go during build
  - Ignored by .gitignore

## Build Process

`cmake/functions.cmake` (`go_program` / `go_example`) compiles the .rc file to `minimal_windows.syso` when configuring a Windows build, before `go build` runs. The .syso file is then automatically linked by the Go toolchain. On macOS and Linux that step does not run at all, and a leftover .syso is skipped because of its `_windows` suffix.

For manual builds:
```powershell
cd examples/minimal
windres -i minimal.rc -o minimal_windows.syso -O coff
go build -o ../../bin/Release/minimal.exe .
```

## Verification

To verify the manifest is embedded:
```powershell
mt.exe -inputresource:bin/Release/minimal.exe -out:extracted.manifest
```

## Linux/macOS

Manifest files are Windows-specific and are ignored on other platforms. The .syso file is never compiled there, and even if one is left over from a Windows build it is not linked, because its name ends in `_windows.syso`.
