package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	PatchedLwjglUrl     = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjglnatives/liblwjgl.dylib"
	PatchedOpenalUrl    = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjglnatives/libopenal.dylib"
	PatchedLwjglFatUrl  = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjglfat.jar"
	PatchedLwjglUtilUrl = "https://raw.githubusercontent.com/GreeniusGenius/m1-prism-launcher-hack-1.8.9/master/lwjgl_util.jar"
)

// EnsureM1Libraries verifies we are on M1/ARM64 and downloads the necessary patched libraries
func EnsureM1Libraries(gameDir string) error {
	// Only run on macOS ARM64
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		return nil
	}

	fmt.Println("Applying Apple Silicon (M1/M2) patches (Libraries & Natives)...")

	nativesDir := filepath.Join(gameDir, "natives")

	// 1. Download Natives
	// liblwjgl.dylib
	lwjglPath := filepath.Join(nativesDir, "liblwjgl.dylib")
	if err := downloadFile(PatchedLwjglUrl, lwjglPath); err != nil {
		return fmt.Errorf("failed to patch liblwjgl.dylib: %w", err)
	}

	// libopenal.dylib
	openalPath := filepath.Join(nativesDir, "libopenal.dylib")
	if err := downloadFile(PatchedOpenalUrl, openalPath); err != nil {
		return fmt.Errorf("failed to patch libopenal.dylib: %w", err)
	}

	// 2. Download Fat Jars (to replace vanilla cp entries)
	// We'll store them in a "m1_libs" folder to keep them separate
	m1LibsDir := filepath.Join(gameDir, "m1_libs")
	if err := os.MkdirAll(m1LibsDir, 0755); err != nil {
		return fmt.Errorf("failed to create m1_libs dir: %w", err)
	}

	// lwjglfat.jar
	fatJarPath := filepath.Join(m1LibsDir, "lwjglfat.jar")
	if err := downloadFile(PatchedLwjglFatUrl, fatJarPath); err != nil {
		return fmt.Errorf("failed to download lwjglfat.jar: %w", err)
	}

	// lwjgl_util.jar
	utilJarPath := filepath.Join(m1LibsDir, "lwjgl_util.jar")
	if err := downloadFile(PatchedLwjglUtilUrl, utilJarPath); err != nil {
		return fmt.Errorf("failed to download lwjgl_util.jar: %w", err)
	}

	// 3. Remove incompatible x86 natives
	// These cause crashes on M1 when the JVM tries to load them
	badLibs := []string{
		"libjinput-osx.dylib",
		"libjinput-osx.jnilib",
		"libtwitchsdk.dylib",
		"openal.dylib",
	}

	for _, lib := range badLibs {
		path := filepath.Join(nativesDir, lib)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Removing incompatible library: %s\n", lib)
			if err := os.Remove(path); err != nil {
				fmt.Printf("Warning: failed to remove %s: %v\n", lib, err)
			}
		}
	}

	return nil
}
