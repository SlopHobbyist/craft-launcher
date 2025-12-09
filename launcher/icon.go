package launcher

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// SetupMinecraftIcon configures a custom icon for the Minecraft process
// This is platform-specific and may not work on all systems
func SetupMinecraftIcon(gameDir string, cmd *exec.Cmd) error {
	// Only attempt on macOS for now
	if runtime.GOOS != "darwin" {
		return nil // Not supported on this platform, but not an error
	}

	// Check if icon file exists in game directory
	iconPath := filepath.Join(gameDir, "minecraft-icon.icns")
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		// No custom icon provided, use default Java icon
		return nil
	}

	// On macOS, we can set the DOCK_ICON environment variable
	// Note: This is a best-effort approach and may not work with all Java versions
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, "DOCK_ICON="+iconPath)

	return nil
}

// CreateMacOSAppBundle creates a macOS .app bundle for Minecraft with custom icon
// This provides the most reliable way to set an icon on macOS
func CreateMacOSAppBundle(gameDir, javaPath string, args []string, iconPath string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", nil
	}

	appPath := filepath.Join(gameDir, "Minecraft.app")
	contentsDir := filepath.Join(appPath, "Contents")
	macOSDir := filepath.Join(contentsDir, "MacOS")
	resourcesDir := filepath.Join(contentsDir, "Resources")

	// Create directory structure
	if err := os.MkdirAll(macOSDir, 0755); err != nil {
		return "", err
	}
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return "", err
	}

	// Copy icon if provided
	if iconPath != "" && fileExists(iconPath) {
		iconDest := filepath.Join(resourcesDir, "minecraft.icns")
		if err := copyFile(iconPath, iconDest); err != nil {
			return "", err
		}
	}

	// Create Info.plist
	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>minecraft-launcher</string>
	<key>CFBundleIconFile</key>
	<string>minecraft</string>
	<key>CFBundleName</key>
	<string>Minecraft 1.8.9</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleIdentifier</key>
	<string>com.craft.minecraft</string>
	<key>CFBundleVersion</key>
	<string>1.8.9</string>
</dict>
</plist>`

	plistPath := filepath.Join(contentsDir, "Info.plist")
	if err := os.WriteFile(plistPath, []byte(plist), 0644); err != nil {
		return "", err
	}

	// Create launcher script
	scriptContent := "#!/bin/bash\n"
	scriptContent += "cd \"" + gameDir + "\"\n"
	scriptContent += "exec \"" + javaPath + "\""
	for _, arg := range args {
		scriptContent += " \"" + arg + "\""
	}
	scriptContent += "\n"

	scriptPath := filepath.Join(macOSDir, "minecraft-launcher")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return "", err
	}

	return appPath, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
