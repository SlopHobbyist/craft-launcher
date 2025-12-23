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
wails build -platform darwin/arm64 -ldflags "$LDFLAGS"

echo ""
echo "==========================================="
echo "Building $APP_NAME for Windows (x64)"
echo "==========================================="
# Note: This requires a C cross-compiler (usually mingw-w64) if specific CGO features are used,
# but Wails often creates a 'portable' internal build or requires the user to have xcode-select/brew tools.
# If this fails, install mingw-w64: brew install mingw-w64
wails build -platform windows/amd64 -ldflags "$LDFLAGS"

echo ""

echo "==========================================="
echo "Build Complete!"
echo "macOS:   $BUILD_DIR/$APP_NAME.app"
echo "Windows: $BUILD_DIR/$APP_NAME.exe"
echo "==========================================="
