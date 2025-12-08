package main

import (
	"context"
	"craft-launcher/launcher"
	"fmt"
	"os"
	"path/filepath"
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
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf("Error getting home dir: %v", err)
	}

	gameDir := filepath.Join(home, ".craft-launcher")
	if err := os.MkdirAll(gameDir, 0755); err != nil {
		return fmt.Sprintf("Error creating game dir: %v", err)
	}

	opts := launcher.LaunchOptions{
		Username:  username,
		GameDir:   gameDir,
		RamMB:     2048,
		VersionID: "1.8.9",
	}

	go func() {
		// Launch is blocking in terms of download, but run.go's Run() starts command non-blocking.
		// However, we want to run the whole logic in background so GUI doesn't freeze during download.
		fmt.Printf("Starting launch for %s...\n", username)
		if err := launcher.Launch(opts); err != nil {
			fmt.Printf("Error launching: %v\n", err)
		}
	}()

	return "Launching..."
}
