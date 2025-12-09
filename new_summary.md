# Craft Launcher - Project Summary

## What It Is

Craft Launcher is a **portable, cross-platform, "one-click" Minecraft launcher** designed to make it incredibly easy for anyone to launch Minecraft 1.8.9 - even complete beginners with zero technical knowledge.

The launcher handles everything automatically:
- Downloads and installs Java 8 (no user installation required)
- Downloads and installs Minecraft 1.8.9
- Launches the game with a single click
- Works completely offline after the first run

**Target Audience**: Players who want to join your modpack and server without dealing with technical setup.

**Current Version**: Minecraft 1.8.9 with Java 8

## Technology Stack

- **Backend**: Go 1.23 with Wails v2.11.0 framework
- **Frontend**: React 18.2.0 + TypeScript + Vite
- **UI**: Native desktop app using Wails (Chromium-based WebView)
- **Build System**: Wails cross-compilation for Windows and macOS

## Current Implementation

### Core Features (Implemented ✅)

#### 1. **Automatic Java Management**
- Detects if Java 8 is already installed locally
- If not found, automatically downloads Azul Zulu JDK 8 for the user's platform
- Extracts and sets up Java in a portable location (no system installation)
- Supports all platforms: Windows (x64, x86, ARM64), macOS (Intel, Apple Silicon), Linux

**Location**: [launcher/java.go](launcher/java.go)

#### 2. **Minecraft Installation & Asset Management**
- Fetches version manifests from Mojang's official API
- Downloads Minecraft 1.8.9 client JAR
- Downloads all required libraries (40+ JAR files)
- Downloads all game assets (textures, sounds, etc.)
- Extracts platform-specific native libraries (.dylib, .dll, .so)

**Locations**:
- [launcher/manifest.go](launcher/manifest.go) - Version lookup
- [launcher/assets.go](launcher/assets.go) - Asset downloading
- [launcher/libraries.go](launcher/libraries.go) - Library management

#### 3. **Game Launching**
- Constructs proper JVM arguments (classpath, memory, native libraries)
- Sets up offline mode authentication (no Microsoft account needed)
- Launches Minecraft with default window size (854x480)
- Captures game logs and displays them in real-time
- Default RAM allocation: 2048MB

**Location**: [launcher/run.go](launcher/run.go)

#### 4. **Apple Silicon (M1/M2/M3) Support**
- **Critical**: Automatically patches LWJGL native libraries for ARM64 compatibility
- Replaces `liblwjgl.dylib` and `libopenal.dylib` with ARM64 versions
- Based on community patch from GreeniusGenius
- **Important**: Does NOT use `-XstartOnFirstThread` flag (causes crashes with LWJGL 2)

**Location**: [launcher/patch_m1.go](launcher/patch_m1.go)
**Documentation**: [macos_support.md](macos_support.md)

#### 5. **User Interface**
- Clean, dark-themed React interface
- Username input field
- Launch button with status updates
- Real-time logging console with copy-to-clipboard
- Auto-opens log on game crash
- Status display showing current operation

**Locations**:
- [frontend/src/App.tsx](frontend/src/App.tsx) - Main UI component
- [frontend/src/components/Console.tsx](frontend/src/components/Console.tsx) - Log display
- [frontend/src/App.css](frontend/src/App.css) - Styling

#### 6. **Cross-Platform Build System**
- Build script for Windows x64 and macOS ARM64
- Outputs standalone executables/apps
- Uses Wails for native compilation

**Location**: [build_releases.sh](build_releases.sh)

### How It Works

```
User clicks "Launch Game"
         ↓
App.tsx calls Go backend (LaunchGame)
         ↓
launcher.Launch() runs in background goroutine
         ↓
Step 1: Check for Java → Download if missing
Step 2: Fetch Minecraft version manifest from Mojang
Step 3: Download all assets (textures, sounds, etc.)
Step 4: Download all libraries (JARs + natives)
Step 5: Download Minecraft client JAR
Step 6: Apply M1 patches if on Apple Silicon
Step 7: Construct JVM arguments
Step 8: Execute: java -cp {classpath} {main_class} {args}
         ↓
Game process runs, logs stream to UI
         ↓
User plays Minecraft!
```

### Data Storage Architecture

Everything is **portable** - stored relative to the launcher:
```
craft-launcher/
├── data/                      # Game directory
│   ├── jre-{os}-{arch}/      # Portable Java installation
│   ├── versions/1.8.9/       # Minecraft JAR
│   ├── libraries/            # All library JARs
│   ├── assets/               # Game textures, sounds, etc.
│   ├── natives/              # Platform-specific native libs
│   └── saves/                # Minecraft worlds (when user plays)
```

No system-wide installation needed - entire folder can be copied to USB drive.

## Future Plans

### Planned Features 🚧

#### 1. **Legacy Fabric Mod Support**
- Add mod loading capability via Legacy Fabric
- Enable running mods on Minecraft 1.8.9
- Maintain compatibility with vanilla server

#### 2. **Automatic Mod & Config Distribution**
- Download mods from internet source (your server/CDN)
- Auto-sync mod configurations
- One-click modpack installation for users
- Update checking for mod changes

#### 3. **Client-Side Anticheat**
- Robust anticheat system running on client
- Prevent common hacks/cheats
- Maintain fair gameplay on your server
- Balance between security and user privacy

## Why This Architecture?

**Go Backend**: Fast, cross-compiles easily, small binaries, excellent for system operations

**React Frontend**: Modern UI development, component reusability, easy to style

**Wails Framework**: Bridges Go + React, creates native desktop apps, no Electron bloat

**Portable Java**: Users don't need admin rights or existing Java installation

**Offline Support**: After first launch, no internet required - perfect for LAN parties

## Building the Launcher

### Prerequisites
- Go 1.23+
- Node.js and npm
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Development
```bash
wails dev
```

### Production Build
```bash
# Build for current platform
wails build

# Build releases (Windows + macOS)
./build_releases.sh
```

Output in `build/bin/`:
- macOS: `craft-launcher.app`
- Windows: `craft-launcher.exe`

## Key Technical Decisions

1. **Java 8 specifically**: Minecraft 1.8.9 requires Java 8 (newer versions have compatibility issues)

2. **Offline UUID authentication**: Uses hardcoded UUID and "null" access token for offline mode - no Microsoft account needed

3. **Azul Zulu JDK**: Free, redistributable, well-maintained Java 8 builds for all platforms

4. **M1 native patching**: Community-sourced solution for Apple Silicon compatibility

5. **Hardcoded 1.8.9**: Launcher is version-specific by design - simplifies code and user experience

## Project Status

**Current State**: Fully functional launcher for Minecraft 1.8.9
- ✅ Java auto-installation
- ✅ Game downloading
- ✅ Cross-platform support (Windows, macOS, Linux)
- ✅ Apple Silicon support
- ✅ Offline mode
- ✅ Real-time logging

**Next Phase**: Mod support + distribution system

**Long-term**: Client-side anticheat integration

---

*Built with the goal of making Minecraft accessible to everyone, regardless of technical skill.*
