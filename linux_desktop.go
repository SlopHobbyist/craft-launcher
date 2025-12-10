package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var desktopFilePath string

// setupLinuxDesktopFile creates a temporary .desktop file for proper taskbar icon display
// This is called on startup for Linux systems only
func setupLinuxDesktopFile() error {
	if runtime.GOOS != "linux" {
		return nil // Not Linux, skip
	}

	// Get executable path and directory
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Icon path - the build script copies launcher-icon.png to the build directory
	// alongside the binary
	iconPath := filepath.Join(exeDir, "launcher-icon.png")

	// Verify icon exists
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		// No icon found - .desktop file will still work but might show generic icon
		fmt.Println("Warning: launcher-icon.png not found, using default icon")
		iconPath = "" // Empty icon path will be handled by desktop environment
	}

	// Get user's .local/share/applications directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	if err := os.MkdirAll(desktopDir, 0755); err != nil {
		return fmt.Errorf("failed to create desktop directory: %w", err)
	}

	// Create .desktop file
	desktopFilePath = filepath.Join(desktopDir, "craft-launcher-temp.desktop")

	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Craft Launcher
Comment=Minecraft 1.8.9 Launcher
Exec=%s
Icon=%s
Terminal=false
Categories=Game;
StartupWMClass=craft-launcher
`, exePath, iconPath)

	if err := os.WriteFile(desktopFilePath, []byte(desktopContent), 0644); err != nil {
		return fmt.Errorf("failed to write .desktop file: %w", err)
	}

	// Set up signal handlers for cleanup
	setupCleanupHandlers()

	return nil
}

// cleanupLinuxDesktopFile removes the temporary .desktop file
func cleanupLinuxDesktopFile() {
	if desktopFilePath != "" {
		os.Remove(desktopFilePath)
		fmt.Println("Cleaned up temporary .desktop file")
	}
}

// setupCleanupHandlers ensures cleanup happens on various exit scenarios
func setupCleanupHandlers() {
	// Handle graceful shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		cleanupLinuxDesktopFile()
		os.Exit(0)
	}()
}
