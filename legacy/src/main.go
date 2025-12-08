package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	jarFile = "test.jar"
)

// JRE download URLs for different platforms
var jreDownloadURLs = map[string]string{
	"darwin-arm64":   "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-macosx_aarch64.tar.gz",
	"darwin-amd64":   "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-macosx_x64.tar.gz",
	"windows-amd64":  "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_x64.zip",
	"windows-386":    "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_i686.zip",
	"windows-arm64":  "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-win_aarch64.zip",
	"windows-arm":    "https://cdn.azul.com/zulu-embedded/bin/zulu8.60.0.21-ca-jdk8.0.322-win_aarch32sf.zip",
	"linux-amd64":    "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-linux_x64.tar.gz",
	"linux-386":      "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-linux_i686.tar.gz",
	"linux-arm64":    "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jdk8.0.412-linux_aarch64.tar.gz",
	"linux-arm":      "https://cdn.azul.com/zulu-embedded/bin/zulu8.60.0.21-ca-jdk8.0.322-linux_aarch32hf.tar.gz",
}

// getJREDir returns the OS and architecture-specific JRE directory name
func getJREDir() string {
	return fmt.Sprintf("jre-%s-%s", runtime.GOOS, runtime.GOARCH)
}

// getDownloadURL returns the download URL for the current platform
func getDownloadURL() (string, error) {
	key := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	url, ok := jreDownloadURLs[key]
	if !ok || url == "" {
		return "", fmt.Errorf("unsupported platform: %s", key)
	}
	return url, nil
}

// printSecurityNotice displays platform-specific security warnings
func printSecurityNotice() {
	fmt.Println("")
	fmt.Println("FIRST RUN NOTICE:")

	if runtime.GOOS == "windows" {
		fmt.Println("- Windows may show 'Windows protected your PC' (SmartScreen)")
		fmt.Println("  Click 'More info' then 'Run anyway' to continue")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("- macOS may show '...cannot be opened because it is from an unidentified developer'")
		fmt.Println("  Right-click the launcher and select 'Open', then click 'Open' in the dialog")
		fmt.Println("  Or go to System Preferences > Security & Privacy and click 'Open Anyway'")
	} else if runtime.GOOS == "linux" {
		fmt.Println("- You may need to make the launcher executable: chmod +x launcher-linux-*")
		fmt.Println("  Your firewall may ask for permission to connect to the internet")
	}

	fmt.Println("- This is needed ONCE to download Java (about 100MB)")
	fmt.Println("")
}

func main() {
	fmt.Println("Starting game launcher...")
	fmt.Printf("Detected platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Get the directory where the launcher executable is located
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(exePath)

	// Paths relative to launcher location
	jreDir := getJREDir()
	jrePath := filepath.Join(baseDir, jreDir)
	jarPath := filepath.Join(baseDir, jarFile)

	fmt.Printf("Looking for JRE in: %s\n", jreDir)

	// Check if JRE is installed
	if !isJREInstalled(jrePath) {
		// Show security notice only on first run when download is needed
		printSecurityNotice()

		fmt.Println("JRE not found. Downloading Java 8...")
		if err := downloadAndInstallJRE(baseDir, jrePath); err != nil {
			fmt.Printf("Error installing JRE: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("JRE installed successfully!")
	} else {
		fmt.Println("JRE found.")
	}

	// Find the java executable
	javaExec := findJavaExecutable(jrePath)
	if javaExec == "" {
		fmt.Println("Error: Could not find java executable in JRE")
		os.Exit(1)
	}

	// Launch the game
	fmt.Printf("Launching %s...\n", jarFile)
	if err := launchGame(javaExec, jarPath); err != nil {
		fmt.Printf("Error launching game: %v\n", err)
		os.Exit(1)
	}
}

// isJREInstalled checks if the JRE directory exists and contains a valid Java installation
func isJREInstalled(jrePath string) bool {
	if _, err := os.Stat(jrePath); os.IsNotExist(err) {
		return false
	}

	// Check if java executable exists
	javaExec := findJavaExecutable(jrePath)
	return javaExec != ""
}

// findJavaExecutable searches for the java binary in the JRE directory
func findJavaExecutable(jrePath string) string {
	// Determine the java executable name based on OS
	javaExecName := "java"
	if runtime.GOOS == "windows" {
		javaExecName = "java.exe"
	}

	// Common paths within extracted JRE
	possiblePaths := []string{
		filepath.Join(jrePath, "bin", javaExecName),
		filepath.Join(jrePath, "Contents", "Home", "bin", javaExecName),
	}

	// Search for any java executable in subdirectories
	var foundPath string
	filepath.Walk(jrePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, filepath.Join("bin", javaExecName)) {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if foundPath != "" {
		return foundPath
	}

	// Check common paths
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// downloadAndInstallJRE downloads and extracts the JRE
func downloadAndInstallJRE(baseDir, jrePath string) error {
	// Get download URL for current platform
	downloadURL, err := getDownloadURL()
	if err != nil {
		return err
	}

	// Determine file extension based on platform
	var tmpFile string
	if runtime.GOOS == "windows" {
		tmpFile = filepath.Join(baseDir, "java_download.zip")
	} else {
		tmpFile = filepath.Join(baseDir, "java_download.tar.gz")
	}
	defer os.Remove(tmpFile)

	// Download JRE
	fmt.Printf("Downloading JRE from: %s\n", downloadURL)
	if err := downloadFile(tmpFile, downloadURL); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Extract JRE based on file type
	fmt.Println("Extracting JRE...")
	if runtime.GOOS == "windows" {
		if err := extractZip(tmpFile, baseDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	} else {
		if err := extractTarGz(tmpFile, baseDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	}

	// The extracted folder might have a different name, so we need to find it and rename it
	if err := renameExtractedJRE(baseDir, jrePath); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	return nil
}

// downloadFile downloads a file from a URL to a local path with timeout and progress
func downloadFile(filepath string, url string) error {
	// Create HTTP client with generous timeout for large downloads
	// This prevents hanging if firewall blocks silently, but allows time for firewall prompts
	client := &http.Client{
		Timeout: 10 * time.Minute, // 10 minutes total timeout
	}

	// Create context with timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	fmt.Println("Connecting to download server...")
	fmt.Println("Note: If you see a firewall popup, please allow the connection.")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed (check firewall settings): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a progress counter
	counter := &writeCounter{total: resp.ContentLength}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	fmt.Println() // New line after progress
	return err
}

// writeCounter counts bytes written and displays progress
type writeCounter struct {
	total   int64
	written int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.written += int64(n)
	wc.printProgress()
	return n, nil
}

func (wc *writeCounter) printProgress() {
	fmt.Printf("\r")
	if wc.total > 0 {
		percent := float64(wc.written) / float64(wc.total) * 100
		fmt.Printf("Downloading... %.0f%% (%d MB / %d MB)", percent, wc.written/1024/1024, wc.total/1024/1024)
	} else {
		fmt.Printf("Downloading... %d MB", wc.written/1024/1024)
	}
}

// extractTarGz extracts a tar.gz file to a destination directory
func extractTarGz(tarGzPath, destDir string) error {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// extractZip extracts a zip file to a destination directory
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create parent directory if needed
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Create file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// renameExtractedJRE finds the extracted JRE directory and renames it to "jre"
func renameExtractedJRE(baseDir, targetPath string) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	// Look for a directory that looks like a JRE (contains "zulu" or "jdk")
	for _, entry := range entries {
		if entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.Contains(name, "zulu") || strings.Contains(name, "jdk") {
				oldPath := filepath.Join(baseDir, entry.Name())
				// Check if this directory contains a java executable
				if findJavaExecutable(oldPath) != "" {
					return os.Rename(oldPath, targetPath)
				}
			}
		}
	}

	return fmt.Errorf("could not find extracted JRE directory")
}

// launchGame launches the JAR file using the embedded JRE
func launchGame(javaExec, jarPath string) error {
	cmd := exec.Command(javaExec, "-jar", jarPath)

	// Configure process attributes for background execution
	setupProcessAttributes(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Println("Game launched successfully!")
	fmt.Println("Launcher will now exit. The game is running in the background.")

	return nil
}

func init() {
	// Verify we're running on supported platform
	key := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	if _, ok := jreDownloadURLs[key]; !ok {
		fmt.Printf("Warning: Platform %s may not be fully supported yet.\n", key)
	}
}
