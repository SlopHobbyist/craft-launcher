package main

import (
	"context"
	"craft-launcher/launcher"
	"craft-launcher/launcher/integrity"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pbnjay/memory"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// SystemInfo holds system information for the frontend
type SystemInfo struct {
	TotalRAM   uint64 `json:"totalRAM"`   // Total RAM in MiB
	Is32Bit    bool   `json:"is32Bit"`    // Whether running 32-bit
	DefaultRAM int    `json:"defaultRAM"` // Default RAM allocation in MiB
	MinRAM     int    `json:"minRAM"`     // Minimum RAM in MiB
	MaxRAM     int    `json:"maxRAM"`     // Maximum RAM in MiB
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetSystemInfo returns system information for RAM configuration
func (a *App) GetSystemInfo() SystemInfo {
	totalRAM := memory.TotalMemory() / (1024 * 1024) // Convert to MiB
	is32Bit := runtime.GOARCH == "386"

	var defaultRAM, minRAM, maxRAM int

	if is32Bit {
		// 32-bit Windows has strict limitations
		defaultRAM = 1024 // 1 GiB
		minRAM = 256      // 0.25 GiB
		maxRAM = 1024     // 1 GiB
	} else {
		// 64-bit systems
		defaultRAM = 2048      // 2 GiB
		minRAM = 2048          // 2 GiB
		maxRAM = int(totalRAM) // System max
	}

	return SystemInfo{
		TotalRAM:   totalRAM,
		Is32Bit:    is32Bit,
		DefaultRAM: defaultRAM,
		MinRAM:     minRAM,
		MaxRAM:     maxRAM,
	}
}

// LaunchGame starts the game
func (a *App) LaunchGame(username string, ramMB int, useFabric bool) string {
	// Portable: Use the directory of the executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("Error getting exe path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// If on macOS inside .app bundle, verify we are not writing inside the bundle if signed?
	// For "portable" request, usually means "alongside the launcher".
	// In macOS .app, the exe is in Contents/MacOS. We probably want to go up to the .app level or just use a folder next to it.
	// But commonly "portable" means "data folder next to binary".

	gameDir := filepath.Join(exeDir, "data")

	// Fix for macOS specific "translocation" or read-only bundle issues?
	// If user runs directly `build/bin/craft-launcher.app/Contents/MacOS/craft-launcher`, exeDir is .../MacOS.
	// We'll put data in .../MacOS/data for now to be strictly portable relative to binary.

	if err := os.MkdirAll(gameDir, 0755); err != nil {
		return fmt.Sprintf("Error creating game dir: %v", err)
	}

	// Validate RAM allocation
	sysInfo := a.GetSystemInfo()
	if ramMB < sysInfo.MinRAM {
		ramMB = sysInfo.MinRAM
	}
	if ramMB > sysInfo.MaxRAM {
		ramMB = sysInfo.MaxRAM
	}

	// Integrity Check
	wailsruntime.EventsEmit(a.ctx, "update-status", "Verifying File Integrity...")
	restored, err := integrity.VerifyAndRestore(gameDir)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Integrity Error: %v", err))
		return fmt.Sprintf("Integrity Error: %v", err)
	}
	if len(restored) > 0 {
		msg := fmt.Sprintf("Restored %d tampered files", len(restored))
		wailsruntime.EventsEmit(a.ctx, "update-status", msg)
		for _, f := range restored {
			wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf(" - Restored: %s", f))
		}
	} else {
		wailsruntime.EventsEmit(a.ctx, "update-status", "Integrity Check: OK")
	}

	opts := launcher.LaunchOptions{
		Username:  username,
		GameDir:   gameDir,
		RamMB:     ramMB,
		VersionID: "1.8.9",
		UseFabric: useFabric,
		StatusCallback: func(status string) {
			wailsruntime.EventsEmit(a.ctx, "update-status", status)
		},
		LogCallback: func(data string) {
			wailsruntime.EventsEmit(a.ctx, "log-data", data)
		},
	}

	go func() {
		// Launch is blocking in terms of download, but run.go's Run() starts command non-blocking.
		// However, we want to run the whole logic in background so GUI doesn't freeze during download.

		// Log platform information for debugging as status messages
		wailsruntime.EventsEmit(a.ctx, "update-status", "=== PLATFORM INFO ===")
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("OS: %s", runtime.GOOS))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Architecture: %s", runtime.GOARCH))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Username: %s", username))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("RAM Allocation: %d GiB (%d MiB)", ramMB/1024, ramMB))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("System RAM: %d GiB (%d MiB)", sysInfo.TotalRAM/1024, sysInfo.TotalRAM))
		wailsruntime.EventsEmit(a.ctx, "update-status", "Version: 1.8.9")
		wailsruntime.EventsEmit(a.ctx, "update-status", "=====================")

		fmt.Printf("Starting launch for %s...\n", username)
		cmd, err := launcher.Launch(opts)
		if err != nil {
			fmt.Printf("Error launching: %v\n", err)
			wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Error: %v", err))
			return
		}

		wailsruntime.EventsEmit(a.ctx, "update-status", "Running")

		// Wait for game to exit
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Game process exited with error: %v\n", err)
			// Decide if we want to show "Crashed" or just "Ready"
			// Usually non-zero exit means crash or force quit
			wailsruntime.EventsEmit(a.ctx, "update-status", "Crashed")
		} else {
			fmt.Printf("Game process exited normally\n")
			wailsruntime.EventsEmit(a.ctx, "update-status", "Ready to Launch")
		}
	}()

	return "Launching..."
}
