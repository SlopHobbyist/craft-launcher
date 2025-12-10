package main

import (
	"context"
	"craft-launcher/launcher"
	"craft-launcher/launcher/integrity"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pbnjay/memory"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	cmd     *exec.Cmd
	cmdLock sync.Mutex
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
func (a *App) LaunchGame(username string, ramMB int, useFabric bool, serverURL string) string {
	a.cmdLock.Lock()
	if a.cmd != nil {
		a.cmdLock.Unlock()
		return "Game is already running!"
	}
	a.cmdLock.Unlock()

	// Portable: Use the directory of the executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("Error getting exe path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	gameDir := filepath.Join(exeDir, "data")
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

	// Integrity Check & Remote Update
	statusCallback := func(msg string) {
		wailsruntime.EventsEmit(a.ctx, "update-status", msg)
	}

	if err := integrity.CheckAndUpdate(gameDir, serverURL, statusCallback); err != nil {
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Update Error: %v", err))
		return fmt.Sprintf("Update Error: %v", err)
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
		// Log platform information
		wailsruntime.EventsEmit(a.ctx, "update-status", "=== PLATFORM INFO ===")
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("OS: %s", runtime.GOOS))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Architecture: %s", runtime.GOARCH))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Username: %s", username))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("RAM Allocation: %d GiB (%d MiB)", ramMB/1024, ramMB))
		wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("System RAM: %d GiB (%d MiB)", sysInfo.TotalRAM/1024, sysInfo.TotalRAM))
		wailsruntime.EventsEmit(a.ctx, "update-status", "Version: 1.8.9")
		wailsruntime.EventsEmit(a.ctx, "update-status", "=====================")

		fmt.Printf("Starting launch for %s...\n", username)

		// Launch and store command
		cmd, err := launcher.Launch(opts)
		if err != nil {
			fmt.Printf("Error launching: %v\n", err)
			wailsruntime.EventsEmit(a.ctx, "update-status", fmt.Sprintf("Error: %v", err))
			return
		}

		a.cmdLock.Lock()
		a.cmd = cmd
		a.cmdLock.Unlock()

		wailsruntime.EventsEmit(a.ctx, "update-status", "Running")

		// Wait for game to exit
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Game process exited with error: %v\n", err)
			wailsruntime.EventsEmit(a.ctx, "update-status", "Crashed")
		} else {
			fmt.Printf("Game process exited normally\n")
			wailsruntime.EventsEmit(a.ctx, "update-status", "Ready to Launch")
		}

		// Cleanup
		a.cmdLock.Lock()
		a.cmd = nil
		a.cmdLock.Unlock()
	}()

	return "Launching..."
}

// ForceStopGame kills the running game process
func (a *App) ForceStopGame() string {
	a.cmdLock.Lock()
	defer a.cmdLock.Unlock()

	if a.cmd != nil && a.cmd.Process != nil {
		if err := a.cmd.Process.Kill(); err != nil {
			return fmt.Sprintf("Error killing process: %v", err)
		}
		return "Force stopped."
	}
	return "No game running."
}
