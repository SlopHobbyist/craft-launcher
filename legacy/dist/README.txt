GAME LAUNCHER - DISTRIBUTION FILES
===================================

This folder contains pre-built launchers for all supported platforms.

USAGE (SIMPLE):
---------------
1. Choose the launcher for your operating system:
   - Windows users: Double-click "launcher-windows.bat"
   - macOS users: Double-click "launcher-macos.command"

2. The launcher will automatically:
   - Detect your processor type (Intel, ARM, Apple Silicon)
   - Download Java 8 if needed (first run only, ~100MB)
   - Launch your game!

WHAT YOU NEED:
--------------
All users should have these files in the same folder:
  - launcher-windows.bat (Windows users use this)
  - launcher-macos.command (macOS users use this)
  - launcher-windows-amd64.exe
  - launcher-windows-arm64.exe
  - launcher-windows-arm.exe
  - launcher-macos-arm64
  - launcher-macos-amd64
  - test.jar

The .bat and .command files automatically choose the right executable for your system.

ADVANCED - Architecture-Specific Launchers:
--------------------------------------------
If you prefer to run the architecture-specific executables directly:

macOS:
  - launcher-macos-arm64        Apple Silicon Macs (M1, M2, M3, M4, etc.)
  - launcher-macos-amd64        Intel Macs

Windows:
  - launcher-windows-amd64.exe  Windows 7/8/8.1/10/11 (64-bit Intel/AMD)
  - launcher-windows-arm64.exe  Windows 10/11 on ARM64 (Snapdragon laptops)
  - launcher-windows-arm.exe    Windows 8/8.1 on ARM32 (Surface RT)

FIRST RUN:
----------
On first run, the launcher will:
1. Detect your platform
2. Download the appropriate Java 8 JRE (~100MB)
3. Extract it to a platform-specific folder (jre-{os}-{arch})
4. Launch your game

IMPORTANT - Security Prompts on First Run:
------------------------------------------
Windows Users:
  - Windows SmartScreen may show "Windows protected your PC"
    Click "More info" then "Run anyway" to continue
  - This is normal for unsigned applications

macOS Users:
  - macOS Gatekeeper may show "cannot be opened because it is from an
    unidentified developer"
    Right-click the launcher and select "Open", then click "Open" in the dialog
    Or: System Preferences > Security & Privacy > click "Open Anyway"
  - This is normal for unsigned applications

All Users:
  - Your firewall (Windows Firewall, Little Snitch, etc.) may ask for permission
    to connect to the internet. This is needed ONCE to download Java.
  - The launcher will wait for you to approve/deny the connection

OFFLINE USE:
------------
After the first run, the JRE is stored locally and the launcher works completely
offline. You can move the entire folder between computers with the same OS/arch.

MULTI-PLATFORM SUPPORT:
-----------------------
You can put multiple launchers in the same folder along with test.jar. Each
launcher will download and use its own platform-specific JRE folder, so the
game folder can be shared across different operating systems.

REQUIREMENTS:
-------------
- Internet connection (first run only)
- ~150MB free disk space per platform
- No admin/root privileges required

For source code and build instructions, see the src/ folder.
