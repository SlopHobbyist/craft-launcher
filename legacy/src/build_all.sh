#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Error: Go compiler not found!"
    echo ""
    echo "Please install Go from: https://go.dev/dl/"
    echo ""
    echo "After installation, run this script again."
    exit 1
fi

echo "Go compiler found: $(go version)"
echo ""

# Get the directory of the script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DIST_DIR="$SCRIPT_DIR/../dist"

# Create dist directory if it doesn't exist
mkdir -p "$DIST_DIR"

echo "Building game launcher for all platforms..."
echo "Output directory: $DIST_DIR"
echo ""

# macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
cd "$SCRIPT_DIR"
GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/launcher-macos-arm64" main.go launch_unix.go
echo "✓ launcher-macos-arm64"

# macOS x86-64 (Intel)
echo "Building for macOS x86-64..."
GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/launcher-macos-amd64" main.go launch_unix.go
echo "✓ launcher-macos-amd64"

# Windows x86-64
echo "Building for Windows x86-64..."
GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/launcher-windows-amd64.exe" main.go launch_windows.go
echo "✓ launcher-windows-amd64.exe"

# Windows x86 32-bit
echo "Building for Windows x86 32-bit..."
GOOS=windows GOARCH=386 go build -o "$DIST_DIR/launcher-windows-386.exe" main.go launch_windows.go
echo "✓ launcher-windows-386.exe"

# Windows ARM64
echo "Building for Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -o "$DIST_DIR/launcher-windows-arm64.exe" main.go launch_windows.go
echo "✓ launcher-windows-arm64.exe"

# Windows ARM32
echo "Building for Windows ARM32..."
GOOS=windows GOARCH=arm go build -o "$DIST_DIR/launcher-windows-arm.exe" main.go launch_windows.go
echo "✓ launcher-windows-arm.exe"

# Linux x86-64
echo "Building for Linux x86-64..."
GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/launcher-linux-amd64" main.go launch_unix.go
echo "✓ launcher-linux-amd64"

# Linux x86 32-bit
echo "Building for Linux x86 32-bit..."
GOOS=linux GOARCH=386 go build -o "$DIST_DIR/launcher-linux-386" main.go launch_unix.go
echo "✓ launcher-linux-386"

# Linux ARM64
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o "$DIST_DIR/launcher-linux-arm64" main.go launch_unix.go
echo "✓ launcher-linux-arm64"

# Linux ARM32
echo "Building for Linux ARM32..."
GOOS=linux GOARCH=arm go build -o "$DIST_DIR/launcher-linux-arm" main.go launch_unix.go
echo "✓ launcher-linux-arm"

echo ""
echo "Creating universal launchers..."

# Create Windows universal launcher (batch script)
cat > "$DIST_DIR/launcher-windows.bat" << 'EOF'
@echo off
REM Universal Windows launcher - auto-detects architecture

REM Get the directory where this batch file is located
set "SCRIPT_DIR=%~dp0"

REM Detect processor architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "LAUNCHER=%SCRIPT_DIR%launcher-windows-amd64.exe"
    goto :launch
)
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set "LAUNCHER=%SCRIPT_DIR%launcher-windows-arm64.exe"
    goto :launch
)
if "%PROCESSOR_ARCHITECTURE%"=="ARM" (
    set "LAUNCHER=%SCRIPT_DIR%launcher-windows-arm.exe"
    goto :launch
)
if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    REM Check if running 32-bit on 64-bit Windows
    if defined PROCESSOR_ARCHITEW6432 (
        set "LAUNCHER=%SCRIPT_DIR%launcher-windows-amd64.exe"
        goto :launch
    ) else (
        set "LAUNCHER=%SCRIPT_DIR%launcher-windows-386.exe"
        goto :launch
    )
)

REM Unknown architecture
echo Error: Unknown processor architecture: %PROCESSOR_ARCHITECTURE%
echo Supported architectures: AMD64 (x86-64), x86 (32-bit), ARM64, ARM
pause
exit /b 1

:launch
REM Check if launcher exists
if not exist "%LAUNCHER%" (
    echo Error: Launcher not found: %LAUNCHER%
    echo Please ensure all launcher files are in the same folder
    pause
    exit /b 1
)

REM Launch the appropriate executable
"%LAUNCHER%"
EOF
echo "✓ launcher-windows.bat (universal)"

# Create macOS universal launcher (shell script)
cat > "$DIST_DIR/launcher-macos.command" << 'EOF'
#!/bin/bash
# Universal macOS launcher - auto-detects architecture

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Detect processor architecture
ARCH=$(uname -m)

if [ "$ARCH" = "arm64" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-macos-arm64"
elif [ "$ARCH" = "x86_64" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-macos-amd64"
else
    echo "Error: Unknown processor architecture: $ARCH"
    echo "Supported architectures: arm64 (Apple Silicon), x86_64 (Intel)"
    read -p "Press Enter to exit..."
    exit 1
fi

# Check if launcher exists
if [ ! -f "$LAUNCHER" ]; then
    echo "Error: Launcher not found: $LAUNCHER"
    echo "Please ensure all launcher files are in the same folder"
    read -p "Press Enter to exit..."
    exit 1
fi

# Make sure launcher is executable
chmod +x "$LAUNCHER"

# Launch the appropriate executable
"$LAUNCHER"
EOF
chmod +x "$DIST_DIR/launcher-macos.command"
echo "✓ launcher-macos.command (universal)"

# Create Linux universal launcher (shell script)
cat > "$DIST_DIR/launcher-linux.sh" << 'EOF'
#!/bin/bash
# Universal Linux launcher - auto-detects architecture

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Detect processor architecture
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-linux-amd64"
elif [ "$ARCH" = "i686" ] || [ "$ARCH" = "i386" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-linux-386"
elif [ "$ARCH" = "aarch64" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-linux-arm64"
elif [ "$ARCH" = "armv7l" ] || [ "$ARCH" = "armv6l" ]; then
    LAUNCHER="$SCRIPT_DIR/launcher-linux-arm"
else
    echo "Error: Unknown processor architecture: $ARCH"
    echo "Supported architectures: x86_64, i686/i386 (32-bit), aarch64 (ARM64), armv7l/armv6l (ARM32)"
    read -p "Press Enter to exit..."
    exit 1
fi

# Check if launcher exists
if [ ! -f "$LAUNCHER" ]; then
    echo "Error: Launcher not found: $LAUNCHER"
    echo "Please ensure all launcher files are in the same folder"
    read -p "Press Enter to exit..."
    exit 1
fi

# Make sure launcher is executable
chmod +x "$LAUNCHER"

# Launch the appropriate executable
"$LAUNCHER"
EOF
chmod +x "$DIST_DIR/launcher-linux.sh"
echo "✓ launcher-linux.sh (universal)"

echo ""
echo "All builds complete!"
echo ""
echo "Distribution files in dist/:"
echo ""
echo "RECOMMENDED FOR USERS (auto-detect architecture):"
echo "  - launcher-windows.bat         Windows (all versions)"
echo "  - launcher-macos.command       macOS (all versions)"
echo "  - launcher-linux.sh            Linux (all versions)"
echo ""
echo "Architecture-specific (advanced users):"
echo "  - launcher-macos-arm64         macOS Apple Silicon"
echo "  - launcher-macos-amd64         macOS Intel"
echo "  - launcher-windows-amd64.exe   Windows 7/8/8.1/10/11 x86-64"
echo "  - launcher-windows-386.exe     Windows 7/8/8.1/10 x86 32-bit"
echo "  - launcher-windows-arm64.exe   Windows 10/11 ARM64"
echo "  - launcher-windows-arm.exe     Windows 8/8.1/10 ARM 32-bit"
echo "  - launcher-linux-amd64         Linux x86-64"
echo "  - launcher-linux-386           Linux x86 32-bit"
echo "  - launcher-linux-arm64         Linux ARM64"
echo "  - launcher-linux-arm           Linux ARM 32-bit"
echo ""
echo "To distribute: Copy ALL launcher files + test.jar to users"
echo "Users choose: launcher-windows.bat OR launcher-macos.command OR launcher-linux.sh"
