@echo off
REM Windows batch script for building all platform releases

setlocal enabledelayedexpansion

set APP_NAME=craft-launcher
set BUILD_DIR=build\bin

set LDFLAGS=

echo ===========================================
echo Building %APP_NAME% for all platforms
echo ===========================================
echo.

REM Process icons if source files exist
if exist "icons\source\launcher-icon.png" (
    echo ===========================================
    echo Processing icons...
    echo ===========================================
    node icons\process-icons.js
    if %errorlevel% neq 0 (
        echo Icon processing skipped - continuing build...
    )
    echo.
)

REM macOS ARM64 (M1/M2/etc)
echo ===========================================
echo Building for macOS ARM64 (Apple Silicon)
echo ===========================================
wails build -platform darwin/arm64if exist "%BUILD_DIR%\%APP_NAME%.app" (
    move /Y "%BUILD_DIR%\%APP_NAME%.app" "%BUILD_DIR%\%APP_NAME%-macos-arm64.app" >nul
)
echo macOS ARM64 build complete
echo.

REM macOS x86-64 (Intel)
echo ===========================================
echo Building for macOS x86-64 (Intel)
echo ===========================================
wails build -platform darwin/amd64if exist "%BUILD_DIR%\%APP_NAME%.app" (
    move /Y "%BUILD_DIR%\%APP_NAME%.app" "%BUILD_DIR%\%APP_NAME%-macos-amd64.app" >nul
)
echo macOS x86-64 build complete
echo.

REM Windows x86-64
echo ===========================================
echo Building for Windows x86-64
echo ===========================================
wails build -platform windows/amd64 %LDFLAGS% -o craft-launcher-windows-amd64.exe
if %errorlevel% equ 0 (
    echo Windows x86-64 build complete
) else (
    echo Windows x86-64 build FAILED
)
echo.

REM Windows x86 (32-bit)
echo ===========================================
echo Building for Windows x86 (32-bit)
echo ===========================================
wails build -platform windows/386 %LDFLAGS% -o craft-launcher-windows-386.exe
if %errorlevel% equ 0 (
    echo Windows x86 32-bit build complete
) else (
    echo Windows x86 32-bit build FAILED
)
echo.

REM Windows ARM
echo ===========================================
echo Building for Windows ARM
echo ===========================================
wails build -platform windows/arm64 %LDFLAGS% -o craft-launcher-windows-arm64.exe
if %errorlevel% equ 0 (
    echo Windows ARM build complete
) else (
    echo Windows ARM build FAILED
)
echo.

REM Linux x86-64
echo ===========================================
echo Building for Linux x86-64
echo ===========================================
wails build -platform linux/amd64 %LDFLAGS% -o craft-launcher-linux-amd64 >nul 2>&1
if %errorlevel% equ 0 (
    echo Linux x86-64 build complete
) else (
    echo Linux x86-64 build skipped (cross-compilation not supported on Windows)
)
echo.

REM Linux ARM
echo ===========================================
echo Building for Linux ARM
echo ===========================================
wails build -platform linux/arm64 %LDFLAGS% -o craft-launcher-linux-arm64 >nul 2>&1
if %errorlevel% equ 0 (
    echo Linux ARM build complete
) else (
    echo Linux ARM build skipped (cross-compilation not supported on Windows)
)
echo.

echo ===========================================
echo All builds complete!
echo ===========================================
echo Build outputs are in: %BUILD_DIR%\
echo.
echo Note: Linux builds require building on a Linux machine
echo       or using Docker with a Linux environment.
echo ===========================================

endlocal
