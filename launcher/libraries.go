package launcher

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// DownloadLibraries downloads all required libraries and extracts natives
func DownloadLibraries(libs []Library, gameDir string) (string, error) {
	libsDir := filepath.Join(gameDir, "libraries")
	nativesDir := filepath.Join(gameDir, "natives")

	if err := os.MkdirAll(nativesDir, 0755); err != nil {
		return "", err
	}

	var cp []string

	for _, lib := range libs {
		if !shouldDownloadLibrary(lib) {
			continue
		}

		// Handle Main Artifact
		if lib.Downloads.Artifact != nil {
			path := filepath.Join(libsDir, lib.Downloads.Artifact.Path)
			if err := ensureLibrary(lib.Downloads.Artifact, path); err != nil {
				fmt.Printf("Failed to download library %s: %v\n", lib.Name, err)
			} else {
				// Convert to absolute path for classpath
				absPath, _ := filepath.Abs(path)
				cp = append(cp, absPath)
			}
		}

		// Handle Natives
		if lib.Natives != nil {
			nativeKey := ""
			switch runtime.GOOS {
			case "windows":
				nativeKey = "windows"
			case "darwin": // macOS
				nativeKey = "osx"
			case "linux":
				nativeKey = "linux"
			}

			if classifier, ok := lib.Natives[nativeKey]; ok {
				if artifact, exists := lib.Downloads.Classifiers[classifier]; exists {
					path := filepath.Join(libsDir, artifact.Path)
					if err := ensureLibrary(artifact, path); err == nil {
						// Extract native
						if err := extractNative(path, nativesDir); err != nil {
							fmt.Printf("Failed to extract native %s: %v\n", lib.Name, err)
						}
					}
				}
			}
		}
	}

	return strings.Join(cp, string(os.PathListSeparator)), nil
}

func shouldDownloadLibrary(lib Library) bool {
	if len(lib.Rules) == 0 {
		return true
	}
	allow := false
	for _, rule := range lib.Rules {
		if rule.Action == "allow" {
			if rule.OS.Name == "" || isOSMatch(rule.OS.Name) {
				allow = true
			}
		} else if rule.Action == "disallow" {
			if rule.OS.Name == "" || isOSMatch(rule.OS.Name) {
				allow = false
			}
		}
	}
	return allow
}

func isOSMatch(osName string) bool {
	switch osName {
	case "osx":
		return runtime.GOOS == "darwin"
	default:
		return runtime.GOOS == osName
	}
}

func ensureLibrary(artifact *Artifact, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return downloadFile(artifact.URL, dest)
	}
	return nil
}

func extractNative(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() || strings.Contains(f.Name, "META-INF") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, f.Name)
		out, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
