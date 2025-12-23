#!/bin/bash

# Exit on error
set -e

# Ensure we can find 'wails' if it's in the standard Go bin location
export PATH=$PATH:$HOME/go/bin

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

# Read Server URL
SERVER_URL=""
if [ -f ".server_url" ]; then
    SERVER_URL=$(cat .server_url | tr -d '\n\r')
    echo "Using Server URL: $SERVER_URL"
else
    echo "Warning: .server_url not found. Using default."
fi

LDFLAGS="-X 'craft-launcher/launcher/integrity.ServerURL=$SERVER_URL'"

echo "==========================================="
echo "Building $APP_NAME for all platforms"
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
wails build -platform darwin/arm64 -ldflags "$LDFLAGS"
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-arm64.app"
fi
echo "✓ macOS ARM64 build complete"
echo ""

# macOS x86-64 (Intel)
echo "==========================================="
echo "Building for macOS x86-64 (Intel)"
echo "==========================================="
wails build -platform darwin/amd64 -ldflags "$LDFLAGS"
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-amd64.app"
fi
echo "✓ macOS x86-64 build complete"
echo ""

# Windows x86-64
echo "==========================================="
echo "Building for Windows x86-64"
echo "==========================================="
wails build -platform windows/amd64 -ldflags "$LDFLAGS" -o craft-launcher-windows-amd64.exe
echo "✓ Windows x86-64 build complete"
echo ""

# Windows x86 (32-bit)
echo "==========================================="
echo "Building for Windows x86 (32-bit)"
echo "==========================================="
wails build -platform windows/386 -ldflags "$LDFLAGS" -o craft-launcher-windows-386.exe
echo "✓ Windows x86 32-bit build complete"
echo ""

# Windows ARM
echo "==========================================="
echo "Building for Windows ARM"
echo "==========================================="
wails build -platform windows/arm64 -ldflags "$LDFLAGS" -o craft-launcher-windows-arm64.exe
echo "✓ Windows ARM build complete"
echo ""

# Linux x86-64
echo "==========================================="
echo "Building for Linux x86-64"
echo "==========================================="
if wails build -platform linux/amd64 -ldflags "$LDFLAGS" -o craft-launcher-linux-amd64 2>&1; then
    echo "✓ Linux x86-64 build complete"
else
    echo "⚠ Linux x86-64 build skipped (cross-compilation not supported on macOS)"
fi
echo ""

# Linux ARM
echo "==========================================="
echo "Building for Linux ARM"
echo "==========================================="
if wails build -platform linux/arm64 -ldflags "$LDFLAGS" -o craft-launcher-linux-arm64 2>&1; then
    echo "✓ Linux ARM build complete"
else
    echo "⚠ Linux ARM build skipped (cross-compilation not supported on macOS)"
fi
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
echo "Note: Linux builds require building on a Linux machine"
echo "      or using Docker with a Linux environment."
echo "==========================================="
