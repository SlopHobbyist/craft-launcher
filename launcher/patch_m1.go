package launcher

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	PatchedLwjglUrl  = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjglnatives/liblwjgl.dylib"
	PatchedOpenalUrl = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjglnatives/libopenal.dylib"
)

func PatchNatives(nativesDir string) error {
	// Only run on macOS ARM64
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		return nil
	}

	fmt.Println("Applying Apple Silicon (M1/M2) patches to natives...")

	// 1. Download patched liblwjgl.dylib
	lwjglPath := filepath.Join(nativesDir, "liblwjgl.dylib")
	if err := downloadFile(PatchedLwjglUrl, lwjglPath); err != nil {
		return fmt.Errorf("failed to patch liblwjgl.dylib: %w", err)
	}

	// 2. Download patched libopenal.dylib
	openalPath := filepath.Join(nativesDir, "libopenal.dylib")
	if err := downloadFile(PatchedOpenalUrl, openalPath); err != nil {
		return fmt.Errorf("failed to patch libopenal.dylib: %w", err)
	}

	return nil
}
