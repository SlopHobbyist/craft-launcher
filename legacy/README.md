# Java Game Launcher

A cross-platform game launcher that automatically downloads and manages Java 8 for your application. No Java installation required from users!

**Note:** This project was completely generated using AI (Claude Code).

## Features

- **Zero Java installation required** - Automatically downloads and installs Java 8 on first run
- **Cross-platform** - Supports macOS, Windows, and Linux (x86-64, x86 32-bit, ARM64, ARM 32-bit)
- **Portable** - Everything stays in the game folder, no system-wide installation
- **Offline-friendly** - After first run, works completely offline
- **Multi-platform folders** - One game folder can contain launchers for all platforms

## Directory Structure

```
java-launcher/
├── src/              Source code and build scripts
│   ├── main.go       Main launcher logic
│   ├── launch_windows.go  Windows-specific code
│   ├── launch_unix.go     Unix/macOS-specific code
│   ├── go.mod        Go module file
│   └── build_all.sh  Build script for all platforms
├── dist/             Distribution files (ready to ship!)
│   ├── launcher-macos-arm64
│   ├── launcher-macos-amd64
│   ├── launcher-windows-amd64.exe
│   ├── launcher-windows-386.exe
│   ├── launcher-windows-arm64.exe
│   ├── launcher-windows-arm.exe
│   ├── launcher-linux-amd64
│   ├── launcher-linux-386
│   ├── launcher-linux-arm64
│   ├── launcher-linux-arm
│   ├── test.jar      Your game JAR file
│   └── README.txt    User instructions
└── README.md         This file
```

## Building

### Prerequisites

You need Go 1.21 or later installed. Download from [https://go.dev/dl/](https://go.dev/dl/)

### Build Scripts

**On macOS/Linux:**
```bash
./src/build_all.sh
```

**On Windows:**
```cmd
src\build_all.bat
```

Both scripts will:
- Check if Go is installed (and provide download link if not)
- Build launchers for all supported platforms
- Create universal launcher scripts

This will create launchers for:
- macOS ARM64 (Apple Silicon)
- macOS x86-64 (Intel)
- Windows x86-64 (64-bit)
- Windows x86 (32-bit)
- Windows ARM64 (64-bit)
- Windows ARM (32-bit)
- Linux x86-64 (64-bit)
- Linux x86 (32-bit)
- Linux ARM64 (64-bit)
- Linux ARM (32-bit)

## Distribution

To distribute your game:

1. Copy the appropriate launcher(s) from `dist/` folder
2. Include your `test.jar` (or rename to your game's JAR)
3. Optionally include `dist/README.txt` for users

Users simply run the launcher and it handles everything else!

## How It Works

1. Launcher detects the operating system and architecture
2. Checks for existing JRE in platform-specific folder (e.g., `jre-windows-amd64/`)
3. If not found, downloads Azul Zulu Java 8 LTS for the platform
4. Extracts and configures the JRE
5. Launches the game using the embedded JRE
6. On Windows, detaches the game process so it continues after launcher exits

## Supported Platforms

- macOS 10.9+ (Intel and Apple Silicon)
- Windows 7/8/8.1/10/11 (x86-64)
- Windows 7/8/8.1/10 (x86 32-bit)
- Windows 10/11 (ARM64)
- Windows 8/8.1/10 (ARM 32-bit)
- Linux (glibc-based distros: Debian, Ubuntu, Fedora, Arch, etc.)
  - x86-64 (64-bit)
  - x86 (32-bit)
  - ARM64 (64-bit)
  - ARM (32-bit - armv6l/armv7l)
  - **Note:** Requires glibc. Does NOT work on musl-based distros (Alpine, Void Linux musl variant)

## Requirements

- Go 1.21+ (for building)
- ~150MB disk space per platform (for JRE)
- Internet connection (first run only)

## Customization

To change the JAR filename, edit `jarFile` constant in `src/main.go`:

```go
const (
    jarFile = "your-game.jar"
)
```

## License

**CC0 1.0 Universal (Public Domain)**

This work has been dedicated to the public domain under CC0 1.0 Universal.

You can:
- Use this code for any purpose (commercial, personal, educational, etc.)
- Modify and distribute it freely
- Use it without any attribution required

To the extent possible under law, the author has waived all copyright and related rights to this work.

For more information: [https://creativecommons.org/publicdomain/zero/1.0/](https://creativecommons.org/publicdomain/zero/1.0/)
