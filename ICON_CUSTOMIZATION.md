# Icon Customization Guide

This guide explains how to customize the launcher icon and the Minecraft game icon on all platforms.

## Launcher Icon Customization

The launcher icon appears in the taskbar/dock and window title bar when the launcher is running.

### Step 1: Prepare Your Icon Files

You need to create icon files in different formats for each platform:

#### For Windows:
- **File**: `build/windows/icon.ico`
- **Format**: Windows ICO file
- **Recommended sizes**: 16x16, 32x32, 48x48, 256x256 pixels
- **Tool**: Use online converters like [ConvertICO](https://convertio.co/png-ico/) or tools like GIMP

#### For macOS:
- **File**: `build/darwin/icon.icns`
- **Format**: macOS ICNS file
- **Recommended sizes**: Multiple sizes from 16x16 to 1024x1024
- **Tool**: Use `iconutil` command-line tool (built into macOS) or [Image2Icon](https://img2icnsapp.com/)

##### Creating ICNS file using iconutil:
```bash
# 1. Create a folder named icon.iconset
mkdir icon.iconset

# 2. Add your PNG files with these exact names:
icon_16x16.png
icon_16x16@2x.png (32x32)
icon_32x32.png
icon_32x32@2x.png (64x64)
icon_128x128.png
icon_128x128@2x.png (256x256)
icon_256x256.png
icon_256x256@2x.png (512x512)
icon_512x512.png
icon_512x512@2x.png (1024x1024)

# 3. Convert to ICNS
iconutil -c icns icon.iconset -o build/darwin/icon.icns
```

#### For Linux:
- **File**: `build/linux/icon.png`
- **Format**: PNG file
- **Recommended size**: 512x512 pixels or larger

### Step 2: Update wails.json (Optional)

You can also specify icons in the `wails.json` configuration file:

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "craft-launcher",
  "outputfilename": "craft-launcher",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "SlopHobbyist",
    "email": "eandeawesome@gmail.com"
  },
  "info": {
    "companyName": "SlopHobbyist",
    "productName": "Craft Launcher",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2024 SlopHobbyist",
    "comments": "A custom Minecraft 1.8.9 launcher"
  }
}
```

### Step 3: Rebuild the Application

After replacing the icon files, rebuild the application:

```bash
# For all platforms (run build_all.sh or build_all.bat)
./build_all.sh

# Or for a specific platform
wails build -platform darwin/arm64     # macOS Apple Silicon
wails build -platform darwin/amd64     # macOS Intel
wails build -platform windows/amd64    # Windows 64-bit
wails build -platform windows/386      # Windows 32-bit
```

---

## Minecraft Game Icon Customization

The Minecraft game process icon appears when Minecraft is running. This is more complex because Java doesn't easily allow setting window icons via command-line arguments.

### Approach 1: Using Java System Properties (Cross-platform)

You can specify an icon for Java applications, but it requires more setup:

1. **Create a custom icon file** and place it in your game directory
2. **Modify the launcher code** to pass icon properties to Java

Add to [launcher/run.go:104](launcher/run.go#L104):

```go
args := []string{
    fmt.Sprintf("-Xmx%dM", opts.RamMB),
    fmt.Sprintf("-Djava.library.path=%s", nativesDir),
    // Add icon specification (path to your custom icon)
    fmt.Sprintf("-Dicon.path=%s", filepath.Join(opts.GameDir, "icon.png")),
    "-cp", realCp,
    pkg.MainClass,
}
```

### Approach 2: Platform-Specific Solutions

#### Windows:
The easiest way is to create a custom wrapper executable with an icon:

1. **Option A**: Use Resource Hacker or similar tools to edit the Java executable's icon
2. **Option B**: Create a launcher wrapper with your icon that starts Java

#### macOS:
For macOS, you can create an `.app` bundle for Minecraft:

1. Create `Minecraft.app/Contents/Info.plist`
2. Add your icon as `Minecraft.app/Contents/Resources/icon.icns`
3. Modify the launcher to start the `.app` instead of direct Java execution

Add this function to [launcher/run.go](launcher/run.go):

```go
// CreateMacOSAppBundle creates a macOS .app bundle for Minecraft with custom icon
func CreateMacOSAppBundle(gameDir, iconPath string) (string, error) {
    appPath := filepath.Join(gameDir, "Minecraft.app")
    contentsDir := filepath.Join(appPath, "Contents")
    macOSDir := filepath.Join(contentsDir, "MacOS")
    resourcesDir := filepath.Join(contentsDir, "Resources")

    // Create directory structure
    os.MkdirAll(macOSDir, 0755)
    os.MkdirAll(resourcesDir, 0755)

    // Copy icon
    if iconPath != "" {
        iconDest := filepath.Join(resourcesDir, "icon.icns")
        // Copy icon file (implement file copy)
        exec.Command("cp", iconPath, iconDest).Run()
    }

    // Create Info.plist
    plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>minecraft</string>
    <key>CFBundleIconFile</key>
    <string>icon</string>
    <key>CFBundleName</key>
    <string>Minecraft</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
</dict>
</plist>`

    plistPath := filepath.Join(contentsDir, "Info.plist")
    os.WriteFile(plistPath, []byte(plist), 0644)

    return appPath, nil
}
```

Then modify the Launch function to use this when on macOS.

#### Linux:
For Linux with X11/Wayland:

1. Set the `_NET_WM_ICON` property
2. Use `xseticon` or similar tools
3. Or create a `.desktop` file with an icon specification

### Approach 3: Recommended Simple Solution

**For most users**, the simplest approach is to:

1. **Accept that the Minecraft window will show the default Java icon** (this is standard behavior)
2. **Only customize the launcher icon** (which is fully supported and documented above)

Minecraft's own official launcher also uses the Java icon for the game window - this is normal Java application behavior.

---

## Quick Reference: Icon File Locations

```
craft-launcher/
├── build/
│   ├── appicon.png          # Generic app icon (PNG, used as fallback)
│   ├── windows/
│   │   └── icon.ico        # Windows launcher icon (REPLACE THIS)
│   ├── darwin/
│   │   └── icon.icns       # macOS launcher icon (CREATE THIS)
│   └── linux/
│       └── icon.png        # Linux launcher icon (CREATE THIS)
```

## Icon Design Tips

1. **Keep it simple**: Icons should be recognizable at small sizes (16x16)
2. **Use high contrast**: Make sure the icon stands out against light and dark backgrounds
3. **Square aspect ratio**: Icons should be square (1:1 ratio)
4. **Transparent background**: Use PNG with transparency for better integration
5. **Test at multiple sizes**: View your icon at 16x16, 32x32, 48x48, and 256x256

## Tools for Creating Icons

### Free Online Tools:
- [ConvertICO](https://convertio.co/png-ico/) - Convert PNG to ICO
- [CloudConvert](https://cloudconvert.com/) - Convert between formats
- [Favicon.io](https://favicon.io/) - Generate favicons and icons

### Desktop Applications:
- **GIMP** (Free, Windows/macOS/Linux) - Full-featured image editor
- **Paint.NET** (Free, Windows) - Simple image editor with ICO export
- **Image2Icon** (Free, macOS) - Easy ICNS creation
- **Adobe Photoshop** (Paid) - Professional icon design

### Command-Line Tools:
- **ImageMagick** - `convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico`
- **iconutil** (macOS) - Create ICNS from iconset folder
- **png2icns** (Linux) - Create ICNS files on Linux

---

## Troubleshooting

**Issue**: Icon doesn't change after rebuild
- **Solution**: Clear the build cache: `rm -rf build/bin/*` then rebuild

**Issue**: macOS icon appears as generic folder
- **Solution**: Run `touch build/bin/*.app` to update modification time, then relaunch

**Issue**: Windows icon doesn't update
- **Solution**: Clear Windows icon cache:
  ```
  ie4uinit.exe -show
  ```

**Issue**: Icon looks blurry
- **Solution**: Ensure you're providing all required icon sizes, especially high-DPI versions (2x)

---

For additional help, see the [Wails documentation on application icons](https://wails.io/docs/guides/application-development#application-icon).
