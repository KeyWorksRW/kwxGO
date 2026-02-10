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

- **minimal.syso**: Compiled resource object (generated, not checked in)
  - Created by `windres -i minimal.rc -o minimal.syso -O coff`
  - Automatically linked by Go during build
  - Ignored by .gitignore

## Build Process

The `scripts/build_minimal.ps1` script automatically compiles the .rc file to .syso before building the Go application. The .syso file is then automatically linked by the Go toolchain.

For manual builds:
```powershell
cd examples/minimal
windres -i minimal.rc -o minimal.syso -O coff
go build -o ../../bin/Release/minimal.exe .
```

## Verification

To verify the manifest is embedded:
```powershell
mt.exe -inputresource:bin/Release/minimal.exe -out:extracted.manifest
```

## Linux/macOS

Manifest files are Windows-specific and are ignored on other platforms. The .syso file will not be linked on non-Windows systems.
