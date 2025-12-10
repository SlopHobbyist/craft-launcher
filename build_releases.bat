@echo off
REM Windows batch script for building releases

setlocal

set APP_NAME=craft-launcher
set BUILD_DIR=build\bin

REM Read Server URL
set SERVER_URL=
if exist ".server_url" (
    set /p SERVER_URL=<.server_url
)

if "%SERVER_URL%"=="" (
    echo Warning: .server_url not found or empty. Using default.
    set "LDFLAGS="
) else (
    echo Using Server URL: %SERVER_URL%
    set "LDFLAGS=-ldflags "-X 'craft-launcher/launcher/integrity.ServerURL=%SERVER_URL%'""
)


REM Process icons if source files exist
if exist "icons\source\launcher-icon.png" (
    echo ===========================================
    echo Processing icons...
    echo ===========================================
    node icons\process-icons.js
    if %errorlevel% neq 0 (
        echo ⚠ Icon processing skipped - continuing build...
    )
    echo.
)

echo ===========================================
echo Building %APP_NAME% for macOS (ARM64)
echo ===========================================
wails build -platform darwin/arm64 %LDFLAGS%
if %errorlevel% neq 0 (
    echo Error building for macOS ARM64
    exit /b %errorlevel%
)

echo.
echo ===========================================
echo Building %APP_NAME% for Windows (x64)
echo ===========================================
REM Note: This requires a C cross-compiler (usually mingw-w64) if specific CGO features are used
wails build -platform windows/amd64 %LDFLAGS%
if %errorlevel% neq 0 (
    echo Error building for Windows x64
    exit /b %errorlevel%
)

echo.
echo ===========================================
echo Build Complete!
echo macOS:   %BUILD_DIR%\%APP_NAME%.app
echo Windows: %BUILD_DIR%\%APP_NAME%.exe
echo ===========================================

endlocal
