package launcher

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// JRE download URLs for different platforms (Java 8)
var jreDownloadURLs = map[string]string{
	"darwin-arm64":  "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-macosx_aarch64.tar.gz",
	"darwin-amd64":  "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-macosx_x64.tar.gz",
	"windows-amd64": "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_x64.zip",
	"windows-386":   "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_i686.zip",
	"windows-arm64": "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_aarch64.zip",
	"linux-amd64":   "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-linux_x64.tar.gz",
}

func EnsureJava(gameDir string) (string, error) {
	// Reverted to native architecture (arm64 on M1) because we are now patching the natives.
	jreDir := filepath.Join(gameDir, fmt.Sprintf("jre-%s-%s", runtime.GOOS, runtime.GOARCH))

	// Check if exists
	if execPath := findJavaExecutable(jreDir); execPath != "" {
		return execPath, nil
	}

	// Download
	fmt.Println("JRE not found, downloading...")
	if err := downloadAndInstallJRE(gameDir, jreDir, runtime.GOARCH); err != nil {
		return "", err
	}

	return findJavaExecutable(jreDir), nil
}

func findJavaExecutable(jrePath string) string {
	// On Windows, prefer javaw.exe (no console) over java.exe
	targetExecs := []string{"java"}
	if runtime.GOOS == "windows" {
		targetExecs = []string{"javaw.exe", "java.exe"}
	}

	// Common paths
	var possiblePaths []string
	for _, execName := range targetExecs {
		possiblePaths = append(possiblePaths,
			filepath.Join(jrePath, "bin", execName),
			filepath.Join(jrePath, "Contents", "Home", "bin", execName),
			filepath.Join(jrePath, "zulu-8.jdk", "Contents", "Home", "bin", execName),
		)
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Deep search
	var found string
	filepath.Walk(jrePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			name := info.Name()
			for _, target := range targetExecs {
				if name == target {
					if strings.Contains(path, "bin") {
						found = path
						// If we found the best candidate (first one), stop immediately
						if target == targetExecs[0] {
							return filepath.SkipAll
						}
					}
				}
			}
		}
		return nil
	})
	return found
}

func downloadAndInstallJRE(baseDir, jrePath, arch string) error {
	key := fmt.Sprintf("%s-%s", runtime.GOOS, arch)
	url, ok := jreDownloadURLs[key]
	if !ok {
		return fmt.Errorf("unsupported platform for auto-java: %s", key)
	}

	tmpFile := filepath.Join(baseDir, "java_install.tmp")
	defer os.Remove(tmpFile)

	if err := downloadFile(url, tmpFile); err != nil {
		return err
	}

	// Extract
	fmt.Println("Extracting Java...")
	if strings.HasSuffix(url, ".zip") {
		if err := extractZip(tmpFile, baseDir); err != nil {
			return err
		}
	} else {
		if err := extractTarGz(tmpFile, baseDir); err != nil {
			return err
		}
	}

	// Rename extracted folder to standard name
	return renameExtractedJRE(baseDir, jrePath)
}

func renameExtractedJRE(baseDir, targetPath string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && (strings.Contains(strings.ToLower(entry.Name()), "zulu") || strings.Contains(strings.ToLower(entry.Name()), "jdk")) {
			// This is likely it. Rename it.
			return os.Rename(filepath.Join(baseDir, entry.Name()), targetPath)
		}
	}
	return fmt.Errorf("could not locate extracted JRE folder")
}

// Helpers reused from legacy logic (simplified)
func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, h.Name)
		if h.Typeflag == tar.TypeDir {
			os.MkdirAll(target, 0755)
		} else if h.Typeflag == tar.TypeReg {
			os.MkdirAll(filepath.Dir(target), 0755)
			w, err := os.Create(target)
			if err != nil {
				return err
			}
			io.Copy(w, tr)
			w.Close()
			os.Chmod(target, os.FileMode(h.Mode))
		}
	}
	return nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
		} else {
			os.MkdirAll(filepath.Dir(path), 0755)
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			rc, err := f.Open()
			if err != nil {
				w.Close()
				return err
			}
			io.Copy(w, rc)
			w.Close()
			rc.Close()
		}
	}
	return nil
}
