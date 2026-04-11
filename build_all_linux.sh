#!/bin/bash

# Exit on error
set -e

# Ensure we can find 'wails3' if it's in the standard Go bin location
export PATH=$PATH:$HOME/go/bin

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

LDFLAGS=""

echo "==========================================="
echo "Building $APP_NAME for all platforms (Linux)"
echo "==========================================="
echo ""

# Process icons
if [ -f "icons/source/launcher-icon.png" ]; then
    echo "==========================================="
    echo "Processing icons..."
    echo "==========================================="
    node icons/process-icons.js || echo "⚠ Icon processing skipped - continuing build..."
    echo ""
fi


# macOS ARM64 (M1/M2/etc)
echo "==========================================="
echo "Building for macOS ARM64 (Apple Silicon)"
echo "==========================================="
GOOS=darwin GOARCH=arm64 wails3 build
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-arm64.app"
fi
echo "✓ macOS ARM64 build complete"
echo ""

# macOS x86-64 (Intel)
echo "==========================================="
echo "Building for macOS x86-64 (Intel)"
echo "==========================================="
GOOS=darwin GOARCH=amd64 wails3 build
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-amd64.app"
fi
echo "✓ macOS x86-64 build complete"
echo ""

# Windows x86-64
echo "==========================================="
echo "Building for Windows x86-64"
echo "==========================================="
GOOS=windows GOARCH=amd64 wails3 build -ldflags "$LDFLAGS" -o craft-launcher-windows-amd64.exe
echo "✓ Windows x86-64 build complete"
echo ""

# Windows x86 (32-bit)
echo "==========================================="
echo "Building for Windows x86 (32-bit)"
echo "==========================================="
GOOS=windows GOARCH=386 wails3 build -ldflags "$LDFLAGS" -o craft-launcher-windows-386.exe
echo "✓ Windows x86 32-bit build complete"
echo ""

# Windows ARM
echo "==========================================="
echo "Building for Windows ARM"
echo "==========================================="
GOOS=windows GOARCH=arm64 wails3 build -ldflags "$LDFLAGS" -o craft-launcher-windows-arm64.exe
echo "✓ Windows ARM build complete"
echo ""

# Linux x86-64
echo "==========================================="
echo "Building for Linux x86-64"
echo "==========================================="
GOOS=linux GOARCH=amd64 wails3 build -ldflags "$LDFLAGS" -o craft-launcher-linux-amd64
echo "✓ Linux x86-64 build complete"
echo ""

# Linux ARM
echo "==========================================="
echo "Building for Linux ARM"
echo "==========================================="
GOOS=linux GOARCH=arm64 wails3 build -ldflags "$LDFLAGS" -o craft-launcher-linux-arm64
echo "✓ Linux ARM build complete"
echo ""


# Copy Linux install script
if [ -f "install_linux.sh" ]; then
    cp install_linux.sh "$BUILD_DIR/"
    echo "✓ install_linux.sh copied to build directory"
fi

# Copy launcher icon for Linux
if [ -f "icons/source/launcher-icon.png" ]; then
    cp icons/source/launcher-icon.png "$BUILD_DIR/"
    echo "✓ launcher-icon.png copied to build directory"
fi

echo "==========================================="
echo "All builds complete!"
echo "==========================================="
echo "Build outputs are in: $BUILD_DIR/"
echo ""
echo "Files created:"
ls -1 "$BUILD_DIR" | grep -E '\.(app|exe)$|^craft-launcher-linux' | sed 's/^/  • /'
echo ""
echo "==========================================="
