#!/usr/bin/env pwsh
# Build script for kwxGO minimal example
# This script builds the C++ libraries (kwxFFI, kwxGO) and then the Go minimal example

$ErrorActionPreference = "Stop"

Write-Host "Building kwxGO Minimal Example" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Get script directory and workspace root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$WorkspaceRoot = Split-Path -Parent $ScriptDir

# Change to workspace root
Push-Location $WorkspaceRoot

try {
    # Step 1: Configure CMake if build directory doesn't exist
    if (!(Test-Path "build")) {
        Write-Host "[1/3] Configuring CMake..." -ForegroundColor Yellow
        cmake -B build -G Ninja
        if ($LASTEXITCODE -ne 0) {
            throw "CMake configuration failed"
        }
        Write-Host "      Configuration complete" -ForegroundColor Green
    } else {
        Write-Host "[1/3] CMake already configured" -ForegroundColor Green
    }

    # Step 2: Build C++ libraries (Release to match kwxFFI)
    Write-Host "[2/3] Building C++ libraries (kwxFFI, kwxGO)..." -ForegroundColor Yellow
    ninja -C build -f build-Release.ninja
    if ($LASTEXITCODE -ne 0) {
        throw "C++ build failed"
    }
    Write-Host "      C++ libraries built successfully" -ForegroundColor Green
    Write-Host "      Output: bin/Release/kwxFFI.dll, kwxGO.dll" -ForegroundColor Gray

    # Step 3: Compile Windows resources (manifest for Common Controls v6)
    Write-Host "[3/4] Compiling Windows resources..." -ForegroundColor Yellow
    Push-Location examples/minimal
    windres -i minimal.rc -o minimal.syso -O coff
    if ($LASTEXITCODE -ne 0) {
        Pop-Location
        throw "Resource compilation failed (windres not found or error)"
    }
    Pop-Location
    Write-Host "      Resources compiled (Common Controls v6 manifest embedded)" -ForegroundColor Green

    # Step 4: Build Go minimal example to bin/Release (where DLLs are)
    # -gcflags="all=-N -l" disables optimizations/inlining for Delve debugging
    Write-Host "[4/4] Building Go minimal example..." -ForegroundColor Yellow

    go build -gcflags="all=-N -l" -o bin/Release/minimal.exe ./examples/minimal/
    if ($LASTEXITCODE -ne 0) {
        throw "Go build failed"
    }
    Write-Host "      Go application built successfully" -ForegroundColor Green

    Write-Host ""
    Write-Host "Build complete! Executable: bin/Release/minimal.exe" -ForegroundColor Cyan
    exit 0

} catch {
    Write-Host ""
    Write-Host "Build failed: $_" -ForegroundColor Red
    exit 1
} finally {
    Pop-Location
}
