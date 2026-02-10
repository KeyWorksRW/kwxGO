#!/usr/bin/env bash
# Build script for kwxGO minimal example
# This script builds the C++ libraries (kwxFFI, kwxGO) and then the Go minimal example

set -e

echo -e "\033[36mBuilding kwxGO Minimal Example\033[0m"
echo -e "\033[36m==============================\033[0m"
echo ""

# Get script directory and workspace root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to workspace root
cd "$WORKSPACE_ROOT"

# Step 1: Configure CMake if build directory doesn't exist
if [ ! -d "build" ]; then
    echo -e "\033[33m[1/3] Configuring CMake...\033[0m"
    cmake -B build -G Ninja
    echo -e "\033[32m      Configuration complete\033[0m"
else
    echo -e "\033[32m[1/3] CMake already configured\033[0m"
fi

# Step 2: Build C++ libraries (Release to match kwxFFI)
echo -e "\033[33m[2/3] Building C++ libraries (kwxFFI, kwxGO)...\033[0m"
ninja -C build -f build-Release.ninja
echo -e "\033[32m      C++ libraries built successfully\033[0m"
echo -e "\033[90m      Output: bin/Release/ (libkwxFFI.so, libkwxGO.so)\033[0m"

# Step 3: Build Go minimal example to bin/Release (where libraries are)
echo -e "\033[33m[3/3] Building Go minimal example...\033[0m"

go build -o bin/Release/minimal ./examples/minimal/
echo -e "\033[32m      Go application built successfully\033[0m"

echo ""
echo -e "\033[36mBuild complete! Executable: bin/Release/minimal\033[0m"
