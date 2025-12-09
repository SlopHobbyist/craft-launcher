package main

import (
	"context"
	"craft-launcher/launcher"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// LaunchGame starts the game
func (a *App) LaunchGame(username string) string {
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

	// Determine RAM allocation based on architecture
	// 32-bit processes can't allocate 2GB heap due to address space limitations
	ramMB := 2048
	if runtime.GOARCH == "386" {
		ramMB = 1024 // 1GB for 32-bit Windows
	}

	opts := launcher.LaunchOptions{
		Username:  username,
		GameDir:   gameDir,
		RamMB:     ramMB,
		VersionID: "1.8.9",
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
