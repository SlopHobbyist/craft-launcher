# Project Summary: Craft Launcher

## Overview
A custom, lightweight, and portable Minecraft 1.8.9 launcher built with Go and Wails. Designed for simplicity, portability, and "offline" tournament/modpack play without requiring Microsoft authentication.

## Technology Stack
*   **Backend**: Go (Golang)
*   **Frontend**: Wails (React + TypeScript + Vite)
*   **Target Platforms**: macOS (Intel/Apple Silicon), Windows (x64)

## Key Features
*   **Offline Authentication**: Simple username input, no Mojang/Microsoft login required (uses offline UUIDs).
*   **Automatic Dependency Management**:
    *   Fetches 1.8.9 manifests and packages directly from Mojang's servers.
    *   Downloads Assets (indexes/objects) and Libraries.
    *   **Java 8 Management**: Automatically checks for and downloads the correct Zulu OpenJDK 8 JRE (native architecture) for the host OS if missing.
*   **Apple Silicon (M1/M2/M3) Support**: 
    *   Minecraft 1.8.9's legacy LWJGL 2 libraries are incompatible with Apple Silicon.
    *   **Solution**: We implemented a custom patcher (`launcher/patch_m1.go`) that hot-swaps the vanilla natives with ARM64-compatible community builds (`liblwjgl.dylib`, `libopenal.dylib`) at runtime.
    *   This allows the game to run with **Native ARM64 Java 8**, offering significantly better performance than Rosetta translation.
*   **Portability**: 
    *   All game data (equivalent to `.minecraft`) is stored in a `data/` folder directly adjacent to the launcher executable.
    *   On macOS, this resides inside the bundle at `Craft Launcher.app/Contents/MacOS/data`.
*   **Clean UI**: Minimalist Dark Mode design with a simple "Play" button.

## Architecture Guidelines
*   **`launcher/`**: Core logic package.
    *   `manifest.go`: Structs and fetchers for Mojang version manifests.
    *   `run.go`: Main orchestrator. Downloads dependencies and constructs the complex Minecraft Java argument string.
    *   `java.go`: JRE detection and installation logic. Configured to download native-arch JREs (amd64/arm64). Handles Windows `javaw.exe` preference.
    *   `patch_m1.go`: Specific logic to download patched native libraries for macOS ARM64.
*   **`app.go`**: Wails Application struct. 
    *   Exposes `LaunchGame(username)` to the frontend.
    *   Handles "Portable" path resolution using `os.Executable()`.
*   **Frontend**: 
    *   Standard React app in `frontend/src`.
    *   Invokes backend via `window.go.main.App.LaunchGame`.

## Recent Change Log (Development History)
1.  **Initialization**: Ported legacy Go launcher logic into a new Wails React-TS template.
2.  **Portability Fix**: Switched storage from `UserHomeDir/.craft-launcher` to `executable_dir/data` to meet portability requirements.
3.  **macOS ARM64 Fix**: 
    *   Encountered `SIGSEGV` with vanilla libraries and Metal renderer.
    *   Implemented `PatchNatives` to download patched libraries from `GreeniusGenius/m1-prism-launcher-hack-1.8.9`.
    *   Ensured `EnsureJava` downloads `darwin-arm64` JRE instead of forcing x86.
4.  **Windows Console Fix**: Updated `launcher/java.go` to prefer finding `javaw.exe` over `java.exe` on Windows to prevent a persistent console window.
5.  **Build Automation**: Created `build_releases.sh` to compile for both macOS and Windows in one step.

## Future Plans
*   **Mod Support**: Logic needed to install/launch Forge (currently vanilla only).
*   **Anti-Cheat**: Foundation laid for client-side integrity checks.
