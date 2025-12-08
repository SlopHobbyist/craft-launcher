@echo off
REM Build script for Windows

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go compiler not found!
    echo.
    echo Please install Go from: https://go.dev/dl/
    echo.
    echo After installation, run this script again.
    pause
    exit /b 1
)

REM Display Go version
echo Go compiler found:
go version
echo.

REM Get the directory of the script
set "SCRIPT_DIR=%~dp0"
set "DIST_DIR=%SCRIPT_DIR%..\dist"

REM Create dist directory if it doesn't exist
if not exist "%DIST_DIR%" mkdir "%DIST_DIR%"

echo Building game launcher for all platforms...
echo Output directory: %DIST_DIR%
echo.

REM macOS ARM64 (Apple Silicon)
echo Building for macOS ARM64...
cd /d "%SCRIPT_DIR%"
set GOOS=darwin
set GOARCH=arm64
go build -o "%DIST_DIR%\launcher-macos-arm64" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build macOS ARM64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-macos-arm64

REM macOS x86-64 (Intel)
echo Building for macOS x86-64...
set GOOS=darwin
set GOARCH=amd64
go build -o "%DIST_DIR%\launcher-macos-amd64" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build macOS x86-64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-macos-amd64

REM Windows x86-64
echo Building for Windows x86-64...
set GOOS=windows
set GOARCH=amd64
go build -o "%DIST_DIR%\launcher-windows-amd64.exe" main.go launch_windows.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Windows x86-64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-windows-amd64.exe

REM Windows x86 32-bit
echo Building for Windows x86 32-bit...
set GOOS=windows
set GOARCH=386
go build -o "%DIST_DIR%\launcher-windows-386.exe" main.go launch_windows.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Windows x86 32-bit
    pause
    exit /b 1
)
echo [32m✓[0m launcher-windows-386.exe

REM Windows ARM64
echo Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -o "%DIST_DIR%\launcher-windows-arm64.exe" main.go launch_windows.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Windows ARM64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-windows-arm64.exe

REM Windows ARM32
echo Building for Windows ARM32...
set GOOS=windows
set GOARCH=arm
go build -o "%DIST_DIR%\launcher-windows-arm.exe" main.go launch_windows.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Windows ARM32
    pause
    exit /b 1
)
echo [32m✓[0m launcher-windows-arm.exe

REM Linux x86-64
echo Building for Linux x86-64...
set GOOS=linux
set GOARCH=amd64
go build -o "%DIST_DIR%\launcher-linux-amd64" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Linux x86-64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-linux-amd64

REM Linux x86 32-bit
echo Building for Linux x86 32-bit...
set GOOS=linux
set GOARCH=386
go build -o "%DIST_DIR%\launcher-linux-386" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Linux x86 32-bit
    pause
    exit /b 1
)
echo [32m✓[0m launcher-linux-386

REM Linux ARM64
echo Building for Linux ARM64...
set GOOS=linux
set GOARCH=arm64
go build -o "%DIST_DIR%\launcher-linux-arm64" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Linux ARM64
    pause
    exit /b 1
)
echo [32m✓[0m launcher-linux-arm64

REM Linux ARM32
echo Building for Linux ARM32...
set GOOS=linux
set GOARCH=arm
go build -o "%DIST_DIR%\launcher-linux-arm" main.go launch_unix.go
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Linux ARM32
    pause
    exit /b 1
)
echo [32m✓[0m launcher-linux-arm

echo.
echo Creating universal launchers...

REM Create Windows universal launcher (batch script)
(
echo @echo off
echo REM Universal Windows launcher - auto-detects architecture
echo.
echo REM Get the directory where this batch file is located
echo set "SCRIPT_DIR=%%~dp0"
echo.
echo REM Detect processor architecture
echo if "%%PROCESSOR_ARCHITECTURE%%"=="AMD64" ^(
echo     set "LAUNCHER=%%SCRIPT_DIR%%launcher-windows-amd64.exe"
echo     goto :launch
echo ^)
echo if "%%PROCESSOR_ARCHITECTURE%%"=="ARM64" ^(
echo     set "LAUNCHER=%%SCRIPT_DIR%%launcher-windows-arm64.exe"
echo     goto :launch
echo ^)
echo if "%%PROCESSOR_ARCHITECTURE%%"=="ARM" ^(
echo     set "LAUNCHER=%%SCRIPT_DIR%%launcher-windows-arm.exe"
echo     goto :launch
echo ^)
echo if "%%PROCESSOR_ARCHITECTURE%%"=="x86" ^(
echo     REM Check if running 32-bit on 64-bit Windows
echo     if defined PROCESSOR_ARCHITEW6432 ^(
echo         set "LAUNCHER=%%SCRIPT_DIR%%launcher-windows-amd64.exe"
echo         goto :launch
echo     ^) else ^(
echo         set "LAUNCHER=%%SCRIPT_DIR%%launcher-windows-386.exe"
echo         goto :launch
echo     ^)
echo ^)
echo.
echo REM Unknown architecture
echo echo Error: Unknown processor architecture: %%PROCESSOR_ARCHITECTURE%%
echo echo Supported architectures: AMD64 ^(x86-64^), x86 ^(32-bit^), ARM64, ARM
echo pause
echo exit /b 1
echo.
echo :launch
echo REM Check if launcher exists
echo if not exist "%%LAUNCHER%%" ^(
echo     echo Error: Launcher not found: %%LAUNCHER%%
echo     echo Please ensure all launcher files are in the same folder
echo     pause
echo     exit /b 1
echo ^)
echo.
echo REM Launch the appropriate executable
echo "%%LAUNCHER%%"
) > "%DIST_DIR%\launcher-windows.bat"
echo [32m✓[0m launcher-windows.bat (universal)

REM Create macOS universal launcher (shell script)
(
echo #!/bin/bash
echo # Universal macOS launcher - auto-detects architecture
echo.
echo # Get the directory where this script is located
echo SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo.
echo # Detect processor architecture
echo ARCH=$(uname -m)
echo.
echo if [ "$ARCH" = "arm64" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-macos-arm64"
echo elif [ "$ARCH" = "x86_64" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-macos-amd64"
echo else
echo     echo "Error: Unknown processor architecture: $ARCH"
echo     echo "Supported architectures: arm64 (Apple Silicon), x86_64 (Intel)"
echo     read -p "Press Enter to exit..."
echo     exit 1
echo fi
echo.
echo # Check if launcher exists
echo if [ ! -f "$LAUNCHER" ]; then
echo     echo "Error: Launcher not found: $LAUNCHER"
echo     echo "Please ensure all launcher files are in the same folder"
echo     read -p "Press Enter to exit..."
echo     exit 1
echo fi
echo.
echo # Make sure launcher is executable
echo chmod +x "$LAUNCHER"
echo.
echo # Launch the appropriate executable
echo "$LAUNCHER"
) > "%DIST_DIR%\launcher-macos.command"
echo [32m✓[0m launcher-macos.command (universal)

REM Create Linux universal launcher (shell script)
(
echo #!/bin/bash
echo # Universal Linux launcher - auto-detects architecture
echo.
echo # Get the directory where this script is located
echo SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo.
echo # Detect processor architecture
echo ARCH=$(uname -m)
echo.
echo if [ "$ARCH" = "x86_64" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-linux-amd64"
echo elif [ "$ARCH" = "i686" ] ^|^| [ "$ARCH" = "i386" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-linux-386"
echo elif [ "$ARCH" = "aarch64" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-linux-arm64"
echo elif [ "$ARCH" = "armv7l" ] ^|^| [ "$ARCH" = "armv6l" ]; then
echo     LAUNCHER="$SCRIPT_DIR/launcher-linux-arm"
echo else
echo     echo "Error: Unknown processor architecture: $ARCH"
echo     echo "Supported architectures: x86_64, i686/i386 (32-bit), aarch64 (ARM64), armv7l/armv6l (ARM32)"
echo     read -p "Press Enter to exit..."
echo     exit 1
echo fi
echo.
echo # Check if launcher exists
echo if [ ! -f "$LAUNCHER" ]; then
echo     echo "Error: Launcher not found: $LAUNCHER"
echo     echo "Please ensure all launcher files are in the same folder"
echo     read -p "Press Enter to exit..."
echo     exit 1
echo fi
echo.
echo # Make sure launcher is executable
echo chmod +x "$LAUNCHER"
echo.
echo # Launch the appropriate executable
echo "$LAUNCHER"
) > "%DIST_DIR%\launcher-linux.sh"
echo [32m✓[0m launcher-linux.sh (universal)

echo.
echo All builds complete!
echo.
echo Distribution files in dist/:
echo.
echo RECOMMENDED FOR USERS (auto-detect architecture):
echo   - launcher-windows.bat         Windows (all versions)
echo   - launcher-macos.command       macOS (all versions)
echo   - launcher-linux.sh            Linux (all versions)
echo.
echo Architecture-specific (advanced users):
echo   - launcher-macos-arm64         macOS Apple Silicon
echo   - launcher-macos-amd64         macOS Intel
echo   - launcher-windows-amd64.exe   Windows 7/8/8.1/10/11 x86-64
echo   - launcher-windows-386.exe     Windows 7/8/8.1/10 x86 32-bit
echo   - launcher-windows-arm64.exe   Windows 10/11 ARM64
echo   - launcher-windows-arm.exe     Windows 8/8.1/10 ARM 32-bit
echo   - launcher-linux-amd64         Linux x86-64
echo   - launcher-linux-386           Linux x86 32-bit
echo   - launcher-linux-arm64         Linux ARM64
echo   - launcher-linux-arm           Linux ARM 32-bit
echo.
echo To distribute: Copy ALL launcher files + test.jar to users
echo Users choose: launcher-windows.bat OR launcher-macos.command OR launcher-linux.sh
echo.
pause
