#!/bin/bash

# Exit on error
set -e

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

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
echo "==========================================="
echo "Build Complete!"
echo "macOS:   $BUILD_DIR/$APP_NAME.app"
echo "Windows: $BUILD_DIR/$APP_NAME.exe"
echo "==========================================="
