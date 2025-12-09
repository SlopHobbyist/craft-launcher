# Quality of Life Improvements - Implementation Summary

This document describes all the quality of life improvements that have been implemented in the Craft Launcher.

## 1. Auto-Hide Log on Normal Game Quit ✅

**What it does:**
- When the game exits normally (not crashed), the log window automatically closes
- This allows users to see the launcher options again without manually closing the log
- If the game crashes, the log stays open for debugging

**Implementation:**
- [frontend/src/App.tsx:19-23](frontend/src/App.tsx#L19-L23) - Added logic to close console when status changes to "Ready to Launch"
- The crash detection still auto-opens the log for debugging

**User Experience:**
- Game quits normally → Log closes automatically → User sees launcher UI
- Game crashes → Log stays open → User can read crash logs

---

## 2. Show Log Controls ✅

**What it does:**
- Added "Show Log While Running" checkbox - when enabled, log opens automatically when launching
- Added "SHOW LOG" button that appears when logs exist but the console is closed
- Users can now reopen logs from previous sessions without restarting

**Implementation:**
- [frontend/src/App.tsx:70-79](frontend/src/App.tsx#L70-L79) - Checkbox and button UI
- [frontend/src/App.tsx:10](frontend/src/App.tsx#L10) - Renamed state from `showLog` to `showLogWhileRunning` for clarity
- [frontend/src/App.css:226-243](frontend/src/App.css#L226-L243) - Styling for the "SHOW LOG" button

**User Experience:**
- Users can choose to auto-open logs when game starts
- After game closes, users can view logs again by clicking "SHOW LOG" button
- Log history is preserved until next launch

---

## 3. Launcher Status Messages in Log Window ✅

**What it does:**
- All launcher status messages (like "Downloading Assets...", "Launching...") now appear at the top of the log window
- Status messages are color-coded in green and bold for easy identification
- A separator line divides launcher messages from game output
- When log is open, users can see what the launcher is doing

**Implementation:**
- [frontend/src/App.tsx:12](frontend/src/App.tsx#L12) - Added `statusHistory` state to track all status updates
- [frontend/src/App.tsx:18](frontend/src/App.tsx#L18) - Status updates are logged with `[LAUNCHER]` prefix
- [frontend/src/components/Console.tsx:33-41](frontend/src/components/Console.tsx#L33-L41) - Displays status history above game logs
- [frontend/src/App.css:205-218](frontend/src/App.css#L205-L218) - Styling for launcher status messages

**User Experience:**
- Users with log open can see launcher progress (downloading, installing, etc.)
- No more wondering if the launcher froze - they can see "Downloading Assets..." etc.
- Status messages persist in log, useful for debugging slow launches

---

## 4. RAM Allocation Input ✅

**What it does:**
- Added RAM allocation input box with smart constraints
- **64-bit systems**: Minimum 2GB, maximum = your computer's total RAM
- **32-bit Windows**: Fixed at 1GB (grayed out and unchangeable) due to 32-bit memory limitations
- Shows system RAM info next to the label
- Validates input automatically - you can't set RAM higher than your system has

**Implementation:**

### Backend (Go):
- [app.go:25-32](app.go#L25-L32) - `SystemInfo` struct with RAM constraints
- [app.go:40-66](app.go#L40-L66) - `GetSystemInfo()` method to detect system RAM and architecture
- [app.go:69](app.go#L69) - Updated `LaunchGame()` to accept `ramMB` parameter
- [app.go:92-99](app.go#L92-L99) - Validates RAM allocation against system limits
- [go.mod:11](go.mod#L11) - Added `github.com/pbnjay/memory` dependency for RAM detection

### Frontend (TypeScript/React):
- [frontend/src/App.tsx:15-16](frontend/src/App.tsx#L15-L16) - Added RAM state and system info
- [frontend/src/App.tsx:20-23](frontend/src/App.tsx#L20-L23) - Fetches system info on startup
- [frontend/src/App.tsx:74-102](frontend/src/App.tsx#L74-L102) - RAM input UI with validation
- [frontend/src/App.tsx:53](frontend/src/App.tsx#L53) - Passes RAM value to backend
- [frontend/src/App.css:58-85](frontend/src/App.css#L58-L85) - Styling for RAM input and disabled state

### Wails Bindings:
- [frontend/wailsjs/go/main/App.d.ts](frontend/wailsjs/go/main/App.d.ts) - Auto-generated TypeScript definitions
- [frontend/wailsjs/go/models.ts](frontend/wailsjs/go/models.ts) - Auto-generated SystemInfo model

**User Experience:**
- 64-bit users can customize RAM from 2GB up to their system's total RAM
- 32-bit Windows users see grayed-out input fixed at 1GB with explanation
- Input shows helpful context: "(System: 16GB)" or "(32-bit limited to 1GB)"
- Invalid values are automatically clamped to valid range

**Technical Details:**
- 32-bit limitation: 32-bit processes have a 2GB address space limit (Windows) or 3GB (Linux), but Minecraft + JVM overhead means 1GB heap is the safe maximum
- macOS note: The minimum was kept at 2GB for 64-bit systems as Minecraft 1.8.9 with mods can struggle with less

---

## 5. Launcher Icon Customization ✅

**What it does:**
- Provides complete documentation for changing the launcher's icon on all platforms
- Created platform-specific icon directories
- Comprehensive guide with tools, tips, and troubleshooting

**Implementation:**
- [ICON_CUSTOMIZATION.md](ICON_CUSTOMIZATION.md) - Complete guide for customizing launcher icons
- [build/darwin/](build/darwin/) - Directory for macOS icon (`.icns` file)
- [build/linux/](build/linux/) - Directory for Linux icon (`.png` file)
- [build/windows/icon.ico](build/windows/icon.ico) - Windows launcher icon (existing)

**User Experience:**
- Clear instructions for creating icons for each platform
- Tool recommendations (free and paid)
- Quick reference table for icon locations
- Design tips for creating effective icons
- Troubleshooting section for common issues

**How to Use:**
1. Replace icon files in `build/` directories with your custom icons
2. Rebuild the launcher using `./build_all.sh` or `build_all.bat`
3. Icon will appear in taskbar, dock, and window title

---

## 6. Minecraft Game Icon Customization ✅

**What it does:**
- Provides infrastructure for setting a custom icon for the Minecraft game window
- Platform-specific implementations (best support on macOS)
- Documentation for all approaches including workarounds

**Implementation:**
- [launcher/icon.go](launcher/icon.go) - Icon customization functions:
  - `SetupMinecraftIcon()` - Sets DOCK_ICON env variable on macOS
  - `CreateMacOSAppBundle()` - Creates macOS .app bundle with custom icon
  - Helper functions for file operations
- [ICON_CUSTOMIZATION.md#minecraft-game-icon-customization](ICON_CUSTOMIZATION.md#minecraft-game-icon-customization) - Complete documentation

**User Experience:**
- **macOS**: Place `minecraft-icon.icns` in game data directory for custom icon
- **Windows/Linux**: Documented that Java shows default icon (normal behavior)
- **Advanced users**: Can create .app bundle for macOS or use platform-specific tools

**How to Use:**
1. Create an `.icns` icon file (macOS) or `.png` (other platforms)
2. Place it as `minecraft-icon.icns` in the game directory (`data/minecraft-icon.icns`)
3. Launcher will automatically use it on supported platforms

**Note:** The official Minecraft launcher also shows the Java icon for the game window. Customizing this is a nice-to-have but not critical.

---

## Files Modified

### Go Backend:
- [app.go](app.go) - Added RAM configuration, system info API
- [go.mod](go.mod) - Added memory detection dependency
- [launcher/icon.go](launcher/icon.go) - New file for icon customization

### TypeScript/React Frontend:
- [frontend/src/App.tsx](frontend/src/App.tsx) - All UI improvements
- [frontend/src/components/Console.tsx](frontend/src/components/Console.tsx) - Status history display
- [frontend/src/App.css](frontend/src/App.css) - All styling updates

### Wails Auto-Generated:
- [frontend/wailsjs/go/main/App.d.ts](frontend/wailsjs/go/main/App.d.ts) - TypeScript bindings
- [frontend/wailsjs/go/models.ts](frontend/wailsjs/go/models.ts) - SystemInfo model

### Documentation:
- [ICON_CUSTOMIZATION.md](ICON_CUSTOMIZATION.md) - New comprehensive icon guide
- [QUALITY_OF_LIFE_IMPROVEMENTS.md](QUALITY_OF_LIFE_IMPROVEMENTS.md) - This file

### Build System:
- [build/darwin/](build/darwin/) - Created for macOS icons
- [build/linux/](build/linux/) - Created for Linux icons

---

## Testing Checklist

Before releasing, test these scenarios:

### Log Management:
- [ ] Launch game → Game quits normally → Log auto-closes
- [ ] Launch game → Game crashes → Log stays open
- [ ] Enable "Show Log While Running" → Launch → Log opens automatically
- [ ] Disable "Show Log While Running" → Launch → Log doesn't open
- [ ] After game closes → Click "SHOW LOG" → Previous logs appear
- [ ] Status messages appear at top of log in green
- [ ] "--- GAME OUTPUT ---" separator appears between launcher and game logs

### RAM Configuration:
- [ ] On 64-bit system → RAM input shows system RAM in label
- [ ] On 64-bit system → Can change RAM from 2GB to system max
- [ ] On 64-bit system → Entering value > max clamps to max
- [ ] On 64-bit system → Entering value < 2GB clamps to 2GB
- [ ] On 32-bit Windows → RAM input is disabled (grayed out)
- [ ] On 32-bit Windows → RAM shows "1024" and can't be changed
- [ ] On 32-bit Windows → Label shows "(32-bit limited to 1GB)"
- [ ] Launch with custom RAM → Game uses correct RAM (check with JVM monitoring)

### Icon Customization:
- [ ] Replace launcher icon → Rebuild → New icon appears in taskbar/dock
- [ ] macOS: .icns icon works correctly at all sizes
- [ ] Windows: .ico icon works correctly at all sizes
- [ ] Linux: .png icon works correctly

### Build Process:
- [ ] `wails generate module` regenerates bindings correctly
- [ ] `./build_all.sh` builds all platforms successfully
- [ ] All platform binaries launch without errors

---

## Future Enhancement Ideas

Potential improvements for future versions:

1. **Progress Bars**: Show download progress with percentage and speed
2. **RAM Presets**: Quick buttons for 2GB, 4GB, 8GB allocations
3. **Log Filtering**: Search/filter logs for specific messages
4. **Log Export**: Save logs to file automatically
5. **Settings Persistence**: Remember username, RAM, and preferences
6. **Theme Customization**: Light/dark mode toggle
7. **Java Version Selection**: Allow choosing between Java 8, 11, 17
8. **Mod Profile Support**: Different RAM and settings for different mod packs
9. **Launch Options**: Custom JVM arguments, game arguments
10. **Update Checker**: Notify when launcher updates are available

---

## Build Instructions

To build the launcher with all improvements:

```bash
# Install dependencies
npm install
go mod download

# Generate Wails bindings (after modifying Go code)
export PATH=$PATH:$HOME/go/bin
wails generate module

# Build for all platforms
./build_all.sh        # macOS/Linux
# or
build_all.bat         # Windows

# Build for specific platform
wails build -platform darwin/arm64
wails build -platform windows/amd64
```

---

## Version History

**v1.1.0** - Quality of Life Update
- ✅ Auto-hide log on normal quit
- ✅ Show log controls (button + checkbox)
- ✅ Launcher status messages in log
- ✅ RAM allocation input with smart constraints
- ✅ Icon customization support and documentation
- ✅ Improved user feedback and visibility

**v1.0.0** - Initial Release
- Basic Minecraft 1.8.9 launcher
- Portable design
- Java auto-download
- macOS M1/M2/M3 support
- Cross-platform (Windows, macOS, Linux)

---

## Credits

**Launcher Developer**: SlopHobbyist (eandeawesome@gmail.com)
**Framework**: Wails v2 (Go + React)
**Quality of Life Improvements**: Implemented December 2024

---

For questions or issues, please refer to:
- [README.md](README.md) - General launcher documentation
- [ICON_CUSTOMIZATION.md](ICON_CUSTOMIZATION.md) - Icon customization guide
- [GitHub Issues](https://github.com/SlopHobbyist/craft-launcher/issues) - Report bugs or request features
