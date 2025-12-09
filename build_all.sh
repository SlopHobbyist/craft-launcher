#!/bin/bash

# Exit on error
set -e

# Ensure we can find 'wails' if it's in the standard Go bin location
export PATH=$PATH:$HOME/go/bin

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

echo "==========================================="
echo "Building $APP_NAME for all platforms"
echo "==========================================="
echo ""

# macOS ARM64 (M1/M2/etc)
echo "==========================================="
echo "Building for macOS ARM64 (Apple Silicon)"
echo "==========================================="
wails build -platform darwin/arm64
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-arm64.app"
fi
echo "✓ macOS ARM64 build complete"
echo ""

# macOS x86-64 (Intel)
echo "==========================================="
echo "Building for macOS x86-64 (Intel)"
echo "==========================================="
wails build -platform darwin/amd64
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-amd64.app"
fi
echo "✓ macOS x86-64 build complete"
echo ""

# Windows x86-64
echo "==========================================="
echo "Building for Windows x86-64"
echo "==========================================="
wails build -platform windows/amd64 -o craft-launcher-windows-amd64.exe
echo "✓ Windows x86-64 build complete"
echo ""

# Windows x86 (32-bit)
echo "==========================================="
echo "Building for Windows x86 (32-bit)"
echo "==========================================="
wails build -platform windows/386 -o craft-launcher-windows-386.exe
echo "✓ Windows x86 32-bit build complete"
echo ""

# Windows ARM
echo "==========================================="
echo "Building for Windows ARM"
echo "==========================================="
wails build -platform windows/arm64 -o craft-launcher-windows-arm64.exe
echo "✓ Windows ARM build complete"
echo ""

# Linux x86-64
echo "==========================================="
echo "Building for Linux x86-64"
echo "==========================================="
if wails build -platform linux/amd64 -o craft-launcher-linux-amd64 2>&1; then
    echo "✓ Linux x86-64 build complete"
else
    echo "⚠ Linux x86-64 build skipped (cross-compilation not supported on macOS)"
fi
echo ""

# Linux ARM
echo "==========================================="
echo "Building for Linux ARM"
echo "==========================================="
if wails build -platform linux/arm64 -o craft-launcher-linux-arm64 2>&1; then
    echo "✓ Linux ARM build complete"
else
    echo "⚠ Linux ARM build skipped (cross-compilation not supported on macOS)"
fi
echo ""

echo "==========================================="
echo "All builds complete!"
echo "==========================================="
echo "Build outputs are in: $BUILD_DIR/"
echo ""
echo "Files created:"
ls -1 "$BUILD_DIR" | grep -E '\.(app|exe)$|^craft-launcher-linux' | sed 's/^/  • /'
echo ""
echo "Note: Linux builds require building on a Linux machine"
echo "      or using Docker with a Linux environment."
echo "==========================================="
