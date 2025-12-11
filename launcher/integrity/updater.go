package integrity

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	LocalManifest   = ".client_manifest.json"
	WhitelistedFile = "options.txt" // Basic whitelist logic
)

// CheckAndUpdate handles the entire update flow
func CheckAndUpdate(gameDir string, serverURL string, statusCallback func(string)) error {
	// Ensure game dir exists
	if err := os.MkdirAll(gameDir, 0755); err != nil {
		return err
	}

	// 1. Fetch Server Manifest
	statusCallback("Checking for updates...")
	serverManifest, err := fetchManifest(serverURL)
	if err != nil {
		// STRICT REQUIREMENT: Refuse to start if unable to connect to server.
		// OBFUSCATION: Do not show IP or detailed error.
		return fmt.Errorf("can't connect to server. You either need to:\n1. Connect to the internet\n2. Wait 30 seconds and try again")
	}

	// 2. Read Local Manifest
	localManifestPath := filepath.Join(gameDir, LocalManifest)
	localManifest, _ := loadLocalManifest(localManifestPath)

	currentVersion := 0
	if localManifest != nil {
		currentVersion = localManifest.Version
	}

	statusCallback(fmt.Sprintf("Local Version: %d, Server Version: %d", currentVersion, serverManifest.Version))

	// 3. Update or Verify
	// We now ALWAYS sync to ensure that even if the version number is the same,
	// any file changes (added/removed/modified) are reflected.
	statusCallback(fmt.Sprintf("Checking for updates (v%d)...", serverManifest.Version))

	// Clean up old files that are not in the new manifest
	if localManifest != nil {
		cleanupOldFiles(gameDir, localManifest, serverManifest, statusCallback)
	}

	if err := syncingUpdate(gameDir, serverURL, serverManifest, statusCallback); err != nil {
		return err
	}
	// Save new manifest as local state
	if err := saveLocalManifest(localManifestPath, serverManifest); err != nil {
		return fmt.Errorf("failed to save local manifest: %w", err)
	}
	statusCallback("Integrity verified & up to date.")

	return nil
}

func fetchManifest(serverURL string) (*Manifest, error) {
	var resp *http.Response
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(serverURL + "/manifest.json")
		if err == nil {
			break
		}

		// Wait and retry if it's a network error (likely macOS permission prompt blocking)
		if i < maxRetries-1 {
			// fmt.Printf("Failed to connect (attempt %d/%d): %v. Retrying in 2s...\n", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned status %s", resp.Status)
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func syncingUpdate(gameDir string, serverURL string, manifest *Manifest, cb func(string)) error {
	for i, file := range manifest.Files {
		// If file exists and override is false, skip it (preserve user data)
		// We only download if it's missing entirely
		localPath := filepath.Join(gameDir, file.Path)
		if !file.Override {
			if _, err := os.Stat(localPath); err == nil {
				// File exists, skipping
				// cb(fmt.Sprintf("Skipping user file: %s", file.Path))
				continue
			}
		}

		// Check if file exists and matches checksum to avoid unnecessary re-download
		if _, err := os.Stat(localPath); err == nil {
			valid, err := verifyFileChecksum(gameDir, file)
			if err == nil && valid {
				// File exists and is valid, skip download
				// cb(fmt.Sprintf("Skipping unchanged file: %s", file.Path))
				continue
			}
		}

		cb(fmt.Sprintf("Downloading [%d/%d]: %s", i+1, len(manifest.Files), file.Path))

		if err := downloadFile(gameDir, serverURL, file); err != nil {
			return fmt.Errorf("failed to download %s: %w", file.Path, err)
		}

		// Verify immediately
		valid, err := verifyFileChecksum(gameDir, file)
		if err != nil {
			return fmt.Errorf("failed to verify %s: %w", file.Path, err)
		}
		if !valid {
			return fmt.Errorf("checksum mismatch after download for %s", file.Path)
		}
	}
	return nil
}

func cleanupOldFiles(gameDir string, oldManifest *Manifest, newManifest *Manifest, cb func(string)) {
	// Create map of new files for quick lookup
	newFiles := make(map[string]bool)
	for _, f := range newManifest.Files {
		newFiles[f.Path] = true
	}

	for _, oldFile := range oldManifest.Files {
		// If old file is NOT in the new manifest, delete it
		if !newFiles[oldFile.Path] {
			cb(fmt.Sprintf("Removing obsolete file: %s", oldFile.Path))
			fullPath := filepath.Join(gameDir, oldFile.Path)
			if err := os.Remove(fullPath); err != nil {
				// Warn but don't fail update
				fmt.Printf("Warning: Failed to remove %s: %v\n", fullPath, err)
			}
		}
	}
}

func downloadFile(gameDir string, serverURL string, file FileInfo) error {
	localPath := filepath.Join(gameDir, file.Path)

	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}

	// Use /files/ prefix as per user example logic (implied or standard)
	// User code: url := fmt.Sprintf("%s/files/%s", ServerURL, file.Path)
	url := fmt.Sprintf("%s/files/%s", serverURL, file.Path)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server download failed: %s", resp.Status)
	}

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func verifyFileChecksum(gameDir string, file FileInfo) (bool, error) {
	localPath := filepath.Join(gameDir, file.Path)
	f, err := os.Open(localPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	checksum := hex.EncodeToString(h.Sum(nil))
	return checksum == file.Checksum, nil
}

func loadLocalManifest(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var m Manifest
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func saveLocalManifest(path string, m *Manifest) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}
