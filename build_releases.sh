#!/bin/bash

# Exit on error
set -e

# Ensure we can find 'wails' if it's in the standard Go bin location
export PATH=$PATH:$HOME/go/bin

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

# Process icons if source files exist
if [ -f "icons/source/launcher-icon.png" ]; then
    echo "==========================================="
    echo "Processing icons..."
    echo "==========================================="
    node icons/process-icons.js || echo "⚠ Icon processing skipped - continuing build..."
    echo ""
fi

echo "==========================================="
echo "Building $APP_NAME for macOS (ARM64)"
echo "==========================================="
wails build -platform darwin/arm64

echo ""
echo "==========================================="
echo "Building $APP_NAME for Windows (x64)"
echo "==========================================="
# Note: This requires a C cross-compiler (usually mingw-w64) if specific CGO features are used,
# but Wails often creates a 'portable' internal build or requires the user to have xcode-select/brew tools.
# If this fails, install mingw-w64: brew install mingw-w64
wails build -platform windows/amd64

echo ""
# Handle Bundled Data
if [ -d "bundled" ]; then
    echo "==========================================="
    echo "Bundling pre-configured data..."
    echo "==========================================="
    
    # For macOS: Copy to .app/Contents/MacOS/data
    MAC_APP="$BUILD_DIR/$APP_NAME.app"
    if [ -d "$MAC_APP" ]; then
        echo "Adding to macOS app..."
        mkdir -p "$MAC_APP/Contents/MacOS/data"
        cp -r bundled/* "$MAC_APP/Contents/MacOS/data/"
    fi
    
    # For Windows: Copy to data folder next to exe
    echo "Adding to Windows build..."
    mkdir -p "$BUILD_DIR/data"
    cp -r bundled/* "$BUILD_DIR/data/"
    
    echo "✓ Bundled data included"
    echo ""
fi

echo "==========================================="
echo "Build Complete!"
echo "macOS:   $BUILD_DIR/$APP_NAME.app"
echo "Windows: $BUILD_DIR/$APP_NAME.exe"
echo "==========================================="
